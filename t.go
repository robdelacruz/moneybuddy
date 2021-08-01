package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var _log *log.Logger
var _termW, _termH int

func main() {
	rand.Seed(time.Now().UnixNano())

	err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
func run(args []string) error {
	flog, err := os.Create("./log.txt")
	if err != nil {
		return err
	}
	defer flog.Close()
	_log = log.New(flog, "", 0)

	sw, parms := parseArgs(args)

	// [-i new_file]  Create and initialize db file
	if sw["i"] != "" {
		dbfile := sw["i"]
		if fileExists(dbfile) {
			return fmt.Errorf("File '%s' already exists. Can't initialize it.\n", dbfile)
		}
		createTables(dbfile)
		return nil
	}

	// Need to specify a db file as first parameter.
	if len(parms) == 0 {
		s := `Usage:

   Specify database file:
	t <db file>

   To initialize new database file:
	t -i <new db file>

`
		fmt.Printf(s)
		return nil
	}

	// Exit if db file doesn't exist.
	dbfile := parms[0]
	if !fileExists(dbfile) {
		return fmt.Errorf(`Database file '%s' doesn't exist. Create one using:
	t -i <filename>
   `, dbfile)
	}

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return fmt.Errorf("Error opening '%s' (%s)\n", dbfile, err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	//http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./"))))
	http.HandleFunc("/", indexHandler(db))
	http.HandleFunc("/api/rootdata", rootdataHandler(db))
	http.HandleFunc("/api/currencies", currenciesHandler(db))
	http.HandleFunc("/api/accounts", accountsHandler(db))
	http.HandleFunc("/api/accountstxns", accountstxnsHandler(db))
	http.HandleFunc("/api/account", accountHandler(db))
	http.HandleFunc("/api/txn", txnHandler(db))
	http.HandleFunc("/api/subscriberoot", subscriberootHandler(db))
	http.HandleFunc("/api/whos", whosHandler(db))

	port := "8000"
	if len(parms) > 1 {
		port = parms[1]
	}
	fmt.Printf("Listening on %s...\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	return err

	return nil
}

func listContains(ss []string, v string) bool {
	for _, s := range ss {
		if v == s {
			return true
		}
	}
	return false
}
func fileExists(file string) bool {
	_, err := os.Stat(file)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func createTables(newfile string) {
	if fileExists(newfile) {
		s := fmt.Sprintf("File '%s' already exists. Can't initialize it.\n", newfile)
		fmt.Printf(s)
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", newfile)
	if err != nil {
		fmt.Printf("Error opening '%s' (%s)\n", newfile, err)
		os.Exit(1)
	}

	ss := []string{
		"CREATE TABLE currency (currency_id INTEGER PRIMARY KEY NOT NULL, currency TEXT, usdrate REAL);",
		"CREATE TABLE account (account_id INTEGER PRIMARY KEY NOT NULL, code TEXT, name TEXT, accounttype INTEGER, currency_id INTEGER);",
		"CREATE TABLE txn (txn_id INTEGER PRIMARY KEY NOT NULL, account_id INTEGER, date TEXT, ref TEXT, desc TEXT, amt REAL);",
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("DB error (%s)\n", err)
		os.Exit(1)
	}
	for _, s := range ss {
		_, err := txexec(tx, s)
		if err != nil {
			tx.Rollback()
			log.Printf("DB error (%s)\n", err)
			os.Exit(1)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("DB error (%s)\n", err)
		os.Exit(1)
	}

	fmt.Printf("Creating test data... ")
	initTestData(db)
	fmt.Printf("Done\n")
}

func initTestData(db *sql.DB) {
	c1 := Currency{
		Currency: "USD",
		Usdrate:  1.0,
	}
	c2 := Currency{
		Currency: "PHP",
		Usdrate:  48.0,
	}
	_, err := createCurrency(db, &c1)
	if err != nil {
		panic(err)
	}
	_, err = createCurrency(db, &c2)
	if err != nil {
		panic(err)
	}

	naccounts := 5 + rand.Intn(25)
	for i := 0; i < naccounts; i++ {
		_, err := createRandomAccount(db)
		if err != nil {
			panic(err)
		}
	}
}

func parseArgs(args []string) (map[string]string, []string) {
	switches := map[string]string{}
	parms := []string{}

	standaloneSwitches := []string{}
	definitionSwitches := []string{"i"}
	fNoMoreSwitches := false
	curKey := ""

	for _, arg := range args {
		if fNoMoreSwitches {
			// any arg after "--" is a standalone parameter
			parms = append(parms, arg)
		} else if arg == "--" {
			// "--" means no more switches to come
			fNoMoreSwitches = true
		} else if strings.HasPrefix(arg, "--") {
			switches[arg[2:]] = "y"
			curKey = ""
		} else if strings.HasPrefix(arg, "-") {
			if listContains(definitionSwitches, arg[1:]) {
				// -a "val"
				curKey = arg[1:]
				continue
			}
			for _, ch := range arg[1:] {
				// -a, -b, -ab
				sch := string(ch)
				if listContains(standaloneSwitches, sch) {
					switches[sch] = "y"
				}
			}
		} else if curKey != "" {
			switches[curKey] = arg
			curKey = ""
		} else {
			// standalone parameter
			parms = append(parms, arg)
		}
	}

	return switches, parms
}

// Helper function to make fmt.Fprintf(w, ...) calls shorter.
// Ex.
// Replace:
//   fmt.Fprintf(w, "<p>Some text %s.</p>", str)
//   fmt.Fprintf(w, "<p>Some other text %s.</p>", str)
// with the shorter version:
//   P := makeFprintf(w)
//   P("<p>Some text %s.</p>", str)
//   P("<p>Some other text %s.</p>", str)
func makeFprintf(w io.Writer) func(format string, a ...interface{}) (n int, err error) {
	return func(format string, a ...interface{}) (n int, err error) {
		return fmt.Fprintf(w, format, a...)
	}
}

func jsonstr(v interface{}) string {
	bs, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return ""
	}
	return string(bs)
}

func indexHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		P := makeFprintf(w)
		printHtmlOpen(P, "Website Title", []string{"/static/bundle.js"})
		printContainerOpen(P)

		printContainerClose(P)
		printHtmlClose(P)
	}
}

func rootdataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rootdata, err := findRootdata(db)
		if err != nil {
			handleErr(w, err, "rootdataHandler")
		}

		w.Header().Set("Content-Type", "application/json")
		P := makeFprintf(w)
		P("%s", jsonstr(rootdata))
	}
}

func currenciesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cc, err := findAllCurrencies(db)
		if err != nil {
			handleErr(w, err, "currenciesHandler")
		}

		w.Header().Set("Content-Type", "application/json")
		P := makeFprintf(w)
		P("%s", jsonstr(cc))
	}
}

func accountsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aa, err := findAllAccounts(db)
		if err != nil {
			handleErr(w, err, "accountsHandler")
		}

		w.Header().Set("Content-Type", "application/json")
		P := makeFprintf(w)
		P("%s", jsonstr(aa))
	}
}

func accountstxnsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aa, err := findAllAccountsTxns(db)
		if err != nil {
			handleErr(w, err, "accountstxnsHandler")
		}

		w.Header().Set("Content-Type", "application/json")
		P := makeFprintf(w)
		P("%s", jsonstr(aa))
	}
}

func accountHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			qid := idtoi(r.FormValue("id"))
			if qid == 0 {
				http.Error(w, "Not found.", 404)
				return
			}
			a, err := findAccount(db, qid)
			if err != nil {
				handleErr(w, err, "accountHandler")
				return
			}
			if a == nil {
				http.Error(w, "Not found.", 404)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			P := makeFprintf(w)
			P("%s", jsonstr(a))
			return
		} else if r.Method == "POST" {
			bs, err := ioutil.ReadAll(r.Body)
			if err != nil {
				handleErr(w, err, "POST accountHandler")
				return
			}
			var a Account
			err = json.Unmarshal(bs, &a)
			if err != nil {
				handleErr(w, err, "POST accountHandler")
				return
			}
			newid, err := createAccount(db, &a)
			if err != nil {
				handleErr(w, err, "POST accountHandler")
				return
			}
			a.Accountid = newid

			w.Header().Set("Content-Type", "application/json")
			P := makeFprintf(w)
			P("%s", jsonstr(a))
			return
		} else if r.Method == "PUT" {
			bs, err := ioutil.ReadAll(r.Body)
			if err != nil {
				handleErr(w, err, "PUT accountHandler")
				return
			}
			var a Account
			err = json.Unmarshal(bs, &a)
			if err != nil {
				handleErr(w, err, "PUT accountHandler")
				return
			}
			err = editAccount(db, &a)
			if err != nil {
				handleErr(w, err, "PUT accountHandler")
				return
			}

			// Inform all data subscribers that a data change occured.
			signalSubs(_subs, &_mu1)

			// Close all subscribers.
			closeSubs(_subs, &_mu1)
			_subs = nil

			_whos = nil

			savedAccount, err := findAccount(db, a.Accountid)
			if err != nil {
				handleErr(w, err, "PUT accountHandler")
				return
			}

			w.Header().Set("Content-Type", "application/json")
			P := makeFprintf(w)
			P("%s", jsonstr(savedAccount))
			return
		} else if r.Method == "DELETE" {
			qid := idtoi(r.FormValue("id"))
			if qid == 0 {
				http.Error(w, "Not found.", 404)
				return
			}
			err := delAccount(db, qid)
			if err != nil {
				handleErr(w, err, "DEL accountHandler")
				return
			}
			return
		}

		http.Error(w, "Use GET/POST/PUT/DELETE", 401)
	}
}
func txnHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			qid := idtoi(r.FormValue("id"))
			if qid == 0 {
				http.Error(w, "Not found.", 404)
				return
			}
			t, err := findTxn(db, qid)
			if err != nil {
				handleErr(w, err, "GET txnHandler")
				return
			}
			if t == nil {
				http.Error(w, "Not found.", 404)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			P := makeFprintf(w)
			P("%s", jsonstr(t))
			return
		} else if r.Method == "POST" {
			bs, err := ioutil.ReadAll(r.Body)
			if err != nil {
				handleErr(w, err, "POST txnHandler")
				return
			}
			var t Txn
			err = json.Unmarshal(bs, &t)
			if err != nil {
				handleErr(w, err, "POST txnHandler")
				return
			}
			newid, err := createTxn(db, &t)
			if err != nil {
				handleErr(w, err, "POST txnHandler")
				return
			}
			t.Txnid = newid

			w.Header().Set("Content-Type", "application/json")
			P := makeFprintf(w)
			P("%s", jsonstr(t))
			return
		} else if r.Method == "PUT" {
			bs, err := ioutil.ReadAll(r.Body)
			if err != nil {
				handleErr(w, err, "PUT txnHandler")
				return
			}
			var t Txn
			err = json.Unmarshal(bs, &t)
			if err != nil {
				handleErr(w, err, "PUT txnHandler")
				return
			}
			err = editTxn(db, &t)
			if err != nil {
				handleErr(w, err, "PUT txnHandler")
				return
			}

			// Inform all data subscribers that a data change occured.
			signalSubs(_subs, &_mu1)

			// Close all subscribers.
			closeSubs(_subs, &_mu1)
			_subs = nil

			w.Header().Set("Content-Type", "application/json")
			P := makeFprintf(w)
			P("%s", jsonstr(t))
			return
		} else if r.Method == "DELETE" {
			qid := idtoi(r.FormValue("id"))
			if qid == 0 {
				http.Error(w, "Not found.", 404)
				return
			}
			err := delTxn(db, qid)
			if err != nil {
				handleErr(w, err, "DEL txnHandler")
				return
			}
			return
		}

		http.Error(w, "Use GET/POST/PUT/DELETE", 401)
	}
}

type SignalChan chan struct{}

var _subs []SignalChan
var _mu1 sync.RWMutex
var _whos []string

func whosHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		P := makeFprintf(w)
		P("num _subs: %d\n", len(_subs))
		P("num _whos: %d\n", len(_whos))

		for i := 0; i < len(_whos); i++ {
			P("_whos[%d]: %s\n", i, _whos[i])
		}
	}
}

func subscriberootHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		qwho := r.FormValue("who")
		if qwho == "" {
			qwho = "(noname)"
		}

		_whos = append(_whos, qwho)
		sub := addSub(&_mu1)
		<-sub

		model, err := findRootdata(db)
		if err != nil {
			handleErr(w, err, "subscribemodelHandler")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		P := makeFprintf(w)
		P("%s", jsonstr(model))
	}
}

func addSub(mu *sync.RWMutex) SignalChan {
	mu.Lock()
	defer mu.Unlock()

	sub := make(SignalChan, 1)
	_subs = append(_subs, sub)
	return sub
}

func signalSubs(subs []SignalChan, mu *sync.RWMutex) {
	fmt.Printf("signalSubs()\n")
	mu.RLock()
	defer mu.RUnlock()

	for _, sub := range subs {
		sub <- struct{}{}
	}
}

func closeSubs(subs []SignalChan, mu *sync.RWMutex) {
	mu.Lock()
	defer mu.Unlock()

	for _, sub := range subs {
		close(sub)
	}
}
