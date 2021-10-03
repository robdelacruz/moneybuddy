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

type SignalChan chan struct{}
type DataSync struct {
	subs map[int64][]SignalChan
	mu   sync.RWMutex
}

var _datasync DataSync

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
	http.HandleFunc("/api/currency", currencyHandler(db))
	http.HandleFunc("/api/book", bookHandler(db))
	http.HandleFunc("/api/account", accountHandler(db))
	http.HandleFunc("/api/txn", txnHandler(db))
	http.HandleFunc("/api/subscriberoot", subscriberootHandler(db))
	http.HandleFunc("/api/rptdata", rptdataHandler(db))

	http.HandleFunc("/api/login", loginHandler(db))
	http.HandleFunc("/api/signup", signupHandler(db))
	http.HandleFunc("/api/password", passwordHandler(db))
	http.HandleFunc("/api/user", userHandler(db))

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

func initTestData(db *sql.DB, username string) {
	u := User{
		Username: username,
	}
	userid, err := createUser(db, &u)
	if err != nil {
		panic(err)
	}
	u.Userid = userid

	err = initUserData(db, &u)
	if err != nil {
		panic(err)
	}

	bb, err := findUserBooks(db, userid)
	if err != nil {
		panic(err)
	}

	for _, b := range bb {
		if b.BookType == SystemBook {
			continue
		}

		naccounts := 5 + rand.Intn(10)
		fmt.Printf("Creating %d random bank accounts for book %d: %s...\n", naccounts, b.Bookid, b.Name)
		for i := 0; i < naccounts; i++ {
			_, err := createRandomBankAccount(db, b.Bookid)
			if err != nil {
				panic(err)
			}
		}

		stocks := []string{"IBM", "CAT", "GE", "AMZN", "AAPL", "MSFT"}
		unitprices := []float64{143.59, 217.71, 103.35, 3298.99, 151.12, 294.60}
		naccounts = rand.Intn(len(stocks)+1) + 2
		if naccounts > len(stocks) {
			naccounts = len(stocks)
		}
		fmt.Printf("Creating %d random stock accounts for book %d: %s...\n", naccounts, b.Bookid, b.Name)
		for i := 0; i < naccounts; i++ {
			_, err := createRandomStockAccount(db, b.Bookid, stocks[i], unitprices[i])
			if err != nil {
				panic(err)
			}
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
		if r.Method != "GET" {
			http.Error(w, "Use GET", 401)
			return
		}
		requser := validateApiUser(db, r)
		if requser == nil {
			http.Error(w, "Invalid user", 401)
			return
		}

		quserid := idtoi(r.FormValue("userid"))
		if quserid == 0 {
			http.Error(w, "Not found.", 404)
			return
		}
		if quserid != requser.Userid {
			http.Error(w, "Invalid user", 401)
			return
		}

		rootdata, err := findRootdata(db, quserid)
		if err != nil {
			handleErr(w, err, "rootdataHandler")
		}

		w.Header().Set("Content-Type", "application/json")
		P := makeFprintf(w)
		P("%s", jsonstr(rootdata))
	}
}

func getbookid(r *http.Request) int64 {
	qbookid := idtoi(r.FormValue("bookid"))
	if qbookid == 0 {
		qbookid = 1 // Default to first book if not specified.
	}
	return qbookid
}

func currencyHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requser := validateApiUser(db, r)
		if requser == nil {
			http.Error(w, "Invalid user", 401)
			return
		}

		if r.Method == "GET" {
			qid := idtoi(r.FormValue("id"))
			if qid == 0 {
				http.Error(w, "Not found.", 404)
				return
			}
			c, err := findCurrency(db, qid)
			if err != nil {
				handleErr(w, err, "GET currencyHandler")
				return
			}
			if c == nil {
				http.Error(w, "Not found.", 404)
				return
			}
			if c.Userid != requser.Userid {
				http.Error(w, "Invalid user", 401)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			P := makeFprintf(w)
			P("%s", jsonstr(c))
			return
		} else if r.Method == "POST" {
			bs, err := ioutil.ReadAll(r.Body)
			if err != nil {
				handleErr(w, err, "POST currencyHandler")
				return
			}
			var c Currency
			err = json.Unmarshal(bs, &c)
			if err != nil {
				handleErr(w, err, "POST currencyHandler")
				return
			}
			if c.Userid != requser.Userid {
				http.Error(w, "Invalid user", 401)
				return
			}

			newid, err := createCurrency(db, &c)
			if err != nil {
				handleErr(w, err, "POST currencyHandler")
				return
			}
			c.Currencyid = newid

			// Inform all data subscribers that a data change occured.
			signalAndCloseSubs(&_datasync, requser.Userid)

			w.Header().Set("Content-Type", "application/json")
			P := makeFprintf(w)
			P("%s", jsonstr(c))
			return
		} else if r.Method == "PUT" {
			bs, err := ioutil.ReadAll(r.Body)
			if err != nil {
				handleErr(w, err, "PUT currencyHandler")
				return
			}
			var c Currency
			err = json.Unmarshal(bs, &c)
			if err != nil {
				handleErr(w, err, "PUT currencyHandler")
				return
			}
			if c.Userid != requser.Userid {
				http.Error(w, "Invalid user", 401)
				return
			}

			err = editCurrency(db, &c)
			if err != nil {
				handleErr(w, err, "PUT currencyHandler")
				return
			}

			// Inform all data subscribers that a data change occured.
			signalAndCloseSubs(&_datasync, requser.Userid)

			w.Header().Set("Content-Type", "application/json")
			P := makeFprintf(w)
			P("%s", jsonstr(c))
			return
		} else if r.Method == "DELETE" {
			qid := idtoi(r.FormValue("id"))
			if qid == 0 {
				http.Error(w, "Not found.", 404)
				return
			}
			c, err := findCurrency(db, qid)
			if err != nil {
				handleErr(w, err, "DEL currencyHandler")
				return
			}
			if c == nil {
				http.Error(w, "Not found.", 404)
				return
			}
			if c.Userid != requser.Userid {
				http.Error(w, "Invalid user", 401)
				return
			}
			err = delCurrency(db, qid)
			if err != nil {
				handleErr(w, err, "DEL currencyHandler")
				return
			}

			// Inform all data subscribers that a data change occured.
			signalAndCloseSubs(&_datasync, requser.Userid)

			return
		}

		http.Error(w, "Use GET/POST/PUT/DELETE", 401)
	}
}

func bookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requser := validateApiUser(db, r)
		if requser == nil {
			http.Error(w, "Invalid user", 401)
			return
		}

		if r.Method == "GET" {
			qid := idtoi(r.FormValue("id"))
			if qid == 0 {
				http.Error(w, "Not found.", 404)
				return
			}
			b, err := findBook(db, qid)
			if err != nil {
				handleErr(w, err, "GET bookHandler")
				return
			}
			if b == nil {
				http.Error(w, "Not found.", 404)
				return
			}
			if b.Userid != requser.Userid {
				http.Error(w, "Invalid user", 401)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			P := makeFprintf(w)
			P("%s", jsonstr(b))
			return
		} else if r.Method == "POST" {
			bs, err := ioutil.ReadAll(r.Body)
			if err != nil {
				handleErr(w, err, "POST bookHandler")
				return
			}
			var b Book
			err = json.Unmarshal(bs, &b)
			if err != nil {
				handleErr(w, err, "POST bookHandler")
				return
			}
			if b.Userid != requser.Userid {
				http.Error(w, "Invalid user", 401)
				return
			}

			newid, err := createBook(db, &b)
			if err != nil {
				handleErr(w, err, "POST bookHandler")
				return
			}
			b.Bookid = newid

			// Inform all data subscribers that a data change occured.
			signalAndCloseSubs(&_datasync, requser.Userid)

			w.Header().Set("Content-Type", "application/json")
			P := makeFprintf(w)
			P("%s", jsonstr(b))
			return
		} else if r.Method == "PUT" {
			bs, err := ioutil.ReadAll(r.Body)
			if err != nil {
				handleErr(w, err, "PUT bookHandler")
				return
			}
			var b Book
			err = json.Unmarshal(bs, &b)
			if err != nil {
				handleErr(w, err, "PUT bookHandler")
				return
			}
			if b.Userid != requser.Userid {
				http.Error(w, "Invalid user", 401)
				return
			}

			err = editBook(db, &b)
			if err != nil {
				handleErr(w, err, "PUT bookHandler")
				return
			}

			// Inform all data subscribers that a data change occured.
			signalAndCloseSubs(&_datasync, requser.Userid)

			w.Header().Set("Content-Type", "application/json")
			P := makeFprintf(w)
			P("%s", jsonstr(b))
			return
		} else if r.Method == "DELETE" {
			qid := idtoi(r.FormValue("id"))
			if qid == 0 {
				http.Error(w, "Not found.", 404)
				return
			}
			b, err := findBook(db, qid)
			if err != nil {
				handleErr(w, err, "DEL bookHandler")
				return
			}
			if b == nil {
				http.Error(w, "Not found.", 404)
				return
			}
			if b.Userid != requser.Userid {
				http.Error(w, "Invalid user", 401)
				return
			}
			err = delBook(db, qid)
			if err != nil {
				handleErr(w, err, "DEL bookHandler")
				return
			}

			// Inform all data subscribers that a data change occured.
			signalAndCloseSubs(&_datasync, requser.Userid)

			return
		}

		http.Error(w, "Use GET/POST/PUT/DELETE", 401)
	}
}

func accountHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requser := validateApiUser(db, r)
		if requser == nil {
			http.Error(w, "Invalid user", 401)
			return
		}

		if r.Method == "GET" {
			qid := idtoi(r.FormValue("id"))
			if qid == 0 {
				http.Error(w, "Not found.", 404)
				return
			}
			a, err := findAccount(db, qid)
			if err != nil {
				handleErr(w, err, "GET accountHandler")
				return
			}
			if a == nil {
				http.Error(w, "Not found.", 404)
				return
			}
			// Check if requester user has access.
			accountUserid, err := findAccountUserid(db, qid)
			if err != nil {
				handleErr(w, err, "GET accountHandler")
				return
			}
			if accountUserid != requser.Userid {
				http.Error(w, "Invalid user", 401)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			P := makeFprintf(w)
			P("%s", jsonstr(a))
			return
		} else if r.Method == "POST" {
			qbookid := getbookid(r)
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
			// Check if requester user has access.
			accountUserid, err := findAccountUserid(db, a.Accountid)
			if err != nil {
				handleErr(w, err, "POST accountHandler")
				return
			}
			if accountUserid != requser.Userid {
				http.Error(w, "Invalid user", 401)
				return
			}

			newid, err := createAccount(db, &a, qbookid)
			if err != nil {
				handleErr(w, err, "POST accountHandler")
				return
			}
			a.Accountid = newid

			// Inform all data subscribers that a data change occured.
			signalAndCloseSubs(&_datasync, requser.Userid)

			w.Header().Set("Content-Type", "application/json")
			P := makeFprintf(w)
			P("%s", jsonstr(a))
			return
		} else if r.Method == "PUT" {
			qbookid := getbookid(r)
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
			// Check if requester user has access.
			accountUserid, err := findAccountUserid(db, a.Accountid)
			if err != nil {
				handleErr(w, err, "PUT accountHandler")
				return
			}
			if accountUserid != requser.Userid {
				http.Error(w, "Invalid user", 401)
				return
			}

			err = editAccount(db, &a)
			if err != nil {
				handleErr(w, err, "PUT accountHandler")
				return
			}

			if qbookid > 0 {
				err = assignAccountToBook(db, a.Accountid, qbookid)
				if err != nil {
					handleErr(w, err, "PUT accountHandler")
					return
				}
			}

			// Inform all data subscribers that a data change occured.
			signalAndCloseSubs(&_datasync, requser.Userid)

			w.Header().Set("Content-Type", "application/json")
			P := makeFprintf(w)
			P("%s", jsonstr(a))
			return
		} else if r.Method == "DELETE" {
			qid := idtoi(r.FormValue("id"))
			if qid == 0 {
				http.Error(w, "Not found.", 404)
				return
			}
			// Check if requester user has access.
			accountUserid, err := findAccountUserid(db, qid)
			if err != nil {
				handleErr(w, err, "DEL accountHandler")
				return
			}
			if accountUserid != requser.Userid {
				http.Error(w, "Invalid user", 401)
				return
			}

			err = delAccount(db, qid)
			if err != nil {
				handleErr(w, err, "DEL accountHandler")
				return
			}

			// Inform all data subscribers that a data change occured.
			signalAndCloseSubs(&_datasync, requser.Userid)

			return
		}

		http.Error(w, "Use GET/POST/PUT/DELETE", 401)
	}
}
func txnHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requser := validateApiUser(db, r)
		if requser == nil {
			http.Error(w, "Invalid user", 401)
			return
		}

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
			// Check if requester user has access.
			txnUserid, err := findTxnUserid(db, qid)
			if err != nil {
				handleErr(w, err, "GET txnHandler")
				return
			}
			if txnUserid != requser.Userid {
				http.Error(w, "Invalid user", 401)
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
			// Check if requester user has access.
			txnUserid, err := findTxnUserid(db, t.Txnid)
			if err != nil {
				handleErr(w, err, "POST txnHandler")
				return
			}
			if txnUserid != requser.Userid {
				http.Error(w, "Invalid user", 401)
				return
			}

			newid, err := createTxn(db, &t)
			if err != nil {
				handleErr(w, err, "POST txnHandler")
				return
			}
			t.Txnid = newid

			// Inform all data subscribers that a data change occured.
			signalAndCloseSubs(&_datasync, requser.Userid)

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
			// Check if requester user has access.
			txnUserid, err := findTxnUserid(db, t.Txnid)
			if err != nil {
				handleErr(w, err, "PUT txnHandler")
				return
			}
			if txnUserid != requser.Userid {
				http.Error(w, "Invalid user", 401)
				return
			}

			err = editTxn(db, &t)
			if err != nil {
				handleErr(w, err, "PUT txnHandler")
				return
			}

			// Inform all data subscribers that a data change occured.
			signalAndCloseSubs(&_datasync, requser.Userid)

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
			// Check if requester user has access.
			txnUserid, err := findTxnUserid(db, qid)
			if err != nil {
				handleErr(w, err, "DEL txnHandler")
				return
			}
			if txnUserid != requser.Userid {
				http.Error(w, "Invalid user", 401)
				return
			}

			err = delTxn(db, qid)
			if err != nil {
				handleErr(w, err, "DEL txnHandler")
				return
			}

			// Inform all data subscribers that a data change occured.
			signalAndCloseSubs(&_datasync, requser.Userid)

			return
		}

		http.Error(w, "Use GET/POST/PUT/DELETE", 401)
	}
}

func rptdataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Use GET", 401)
			return
		}
		requser := validateApiUser(db, r)
		if requser == nil {
			http.Error(w, "Invalid user", 401)
			return
		}

		quserid := idtoi(r.FormValue("userid"))
		if quserid == 0 {
			http.Error(w, "Not found.", 404)
			return
		}
		if quserid != requser.Userid {
			http.Error(w, "Invalid user", 401)
			return
		}

		qcurrencyid := idtoi(r.FormValue("currencyid"))
		if qcurrencyid == 0 {
			qcurrencyid = 1
		}

		rptdata, err := findRptdata(db, quserid, qcurrencyid)
		if err != nil {
			handleErr(w, err, "rptdataHandler")
		}

		w.Header().Set("Content-Type", "application/json")
		P := makeFprintf(w)
		P("%s", jsonstr(rptdata))
	}
}

func subscriberootHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Use GET", 401)
			return
		}
		requser := validateApiUser(db, r)
		if requser == nil {
			http.Error(w, "Invalid user", 401)
			return
		}

		quserid := idtoi(r.FormValue("userid"))
		if quserid == 0 {
			http.Error(w, "Not found.", 404)
			return
		}
		if quserid != requser.Userid {
			http.Error(w, "Invalid user", 401)
			return
		}

		sub := make(SignalChan, 1)
		addSub(&_datasync, quserid, sub)
		<-sub

		rootdata, err := findRootdata(db, quserid)
		if err != nil {
			handleErr(w, err, "subscriberootHandler")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		P := makeFprintf(w)
		P("%s", jsonstr(rootdata))
	}
}

func addSub(ds *DataSync, userid int64, sub SignalChan) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if ds.subs == nil {
		ds.subs = make(map[int64][]SignalChan)
	}
	ds.subs[userid] = append(ds.subs[userid], sub)
}
func signalSubs(ds *DataSync, userid int64) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	for _, ch := range ds.subs[userid] {
		ch <- struct{}{}
	}
}
func closeSubs(ds *DataSync, userid int64) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	for _, sub := range ds.subs[userid] {
		close(sub)
	}
	ds.subs[userid] = nil
}
func signalAndCloseSubs(ds *DataSync, userid int64) {
	signalSubs(ds, userid)
	closeSubs(ds, userid)
}

type LoginResult struct {
	Userid   int64  `json:"userid"`
	Username string `json:"username"`
	Sig      string `json:"sig"`
	Error    string `json:"error"`
}

func loginHandler(db *sql.DB) http.HandlerFunc {
	type LoginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Use POST method", 401)
			return
		}
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleErr(w, err, "loginHandler")
			return
		}
		var loginreq LoginReq
		err = json.Unmarshal(bs, &loginreq)
		if err != nil {
			handleErr(w, err, "loginHandler")
			return
		}

		u, err := findUserByUsername(db, loginreq.Username)
		if err != nil {
			handleErr(w, err, "loginHandler")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		P := makeFprintf(w)

		var result LoginResult
		if u == nil {
			result.Error = fmt.Sprintf("Username '%s' not found", loginreq.Username)
			bs, _ = json.MarshalIndent(result, "", "\t")
			P("%s\n", string(bs))
			return
		}

		// Log in user.
		sig, err := login(u, loginreq.Password)
		if err != nil {
			result.Error = fmt.Sprintf("%s", err)
		}
		result.Userid = u.Userid
		result.Username = u.Username
		result.Sig = sig

		bs, _ = json.MarshalIndent(result, "", "\t")
		P("%s\n", string(bs))
	}
}

func signupHandler(db *sql.DB) http.HandlerFunc {
	type SignupReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Use POST method", 401)
			return
		}

		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleErr(w, err, "signupHandler")
			return
		}
		var signupreq SignupReq
		err = json.Unmarshal(bs, &signupreq)
		if err != nil {
			handleErr(w, err, "signupHandler")
			return
		}
		if signupreq.Username == "" {
			http.Error(w, "username required", 401)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		P := makeFprintf(w)

		// Attempt to sign up new user.
		var result LoginResult
		u, err := signup(db, signupreq.Username, signupreq.Password)
		if err != nil {
			result.Error = fmt.Sprintf("%s", err)
			bs, _ := json.MarshalIndent(result, "", "\t")
			P("%s\n", string(bs))
			return
		}

		// Log in the newly signed up user.
		sig, err := login(u, signupreq.Password)
		if err != nil {
			result.Error = fmt.Sprintf("%s", err)
		}
		result.Userid = u.Userid
		result.Username = u.Username
		result.Sig = sig

		bs, _ = json.MarshalIndent(result, "", "\t")
		P("%s\n", string(bs))
	}
}

func passwordHandler(db *sql.DB) http.HandlerFunc {
	type PasswordReq struct {
		Userid      int64  `json:"userid"`
		Password    string `json:"password"`
		NewPassword string `json:"newpassword"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Use POST method", 401)
			return
		}
		requser := validateApiUser(db, r)
		if requser == nil {
			http.Error(w, "Invalid user", 401)
			return
		}

		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleErr(w, err, "passwordHandler")
			return
		}
		var passwordreq PasswordReq
		err = json.Unmarshal(bs, &passwordreq)
		if err != nil {
			handleErr(w, err, "passwordHandler")
			return
		}
		if passwordreq.Userid != requser.Userid {
			http.Error(w, "Invalid user", 401)
			return
		}

		u, err := findUser(db, passwordreq.Userid)
		if err != nil {
			handleErr(w, err, "passwordHandler")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		P := makeFprintf(w)

		var result LoginResult
		if u == nil {
			result.Error = fmt.Sprintf("Userid %d not found", passwordreq.Userid)
			bs, _ = json.MarshalIndent(result, "", "\t")
			P("%s\n", string(bs))
			return
		}

		if !validateHash(u.Password, passwordreq.Password) {
			result.Error = "Incorrect password"
			bs, _ = json.MarshalIndent(result, "", "\t")
			P("%s\n", string(bs))
			return
		}

		// Update the password
		u.Password = genHash(passwordreq.NewPassword)
		err = editUser(db, u)
		if err != nil {
			handleErr(w, err, "passwordHandler")
			return
		}

		// Log in again with the new password to return new user signature.
		sig, err := login(u, passwordreq.NewPassword)
		if err != nil {
			result.Error = fmt.Sprintf("%s", err)
		}
		result.Userid = u.Userid
		result.Username = u.Username
		result.Sig = sig

		bs, _ = json.MarshalIndent(result, "", "\t")
		P("%s\n", string(bs))
	}
}

func userHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requser := validateApiUser(db, r)
		if requser == nil {
			http.Error(w, "Invalid user", 401)
			return
		}

		if r.Method == "GET" {
			qid := idtoi(r.FormValue("id"))
			if qid == 0 {
				http.Error(w, "Not found.", 404)
				return
			}
			// Only admin or same user can GET
			if requser.Userid != 1 && qid != requser.Userid {
				http.Error(w, "Invalid user", 401)
				return
			}

			u, err := findUser(db, qid)
			if err != nil {
				handleErr(w, err, "GET userHandler")
				return
			}
			if u == nil {
				http.Error(w, "Not found.", 404)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			P := makeFprintf(w)
			P("%s", jsonstr(u))
			return
		} else if r.Method == "POST" {
			// Only admin can POST
			if requser.Userid != 1 {
				http.Error(w, "Invalid user", 401)
				return
			}

			bs, err := ioutil.ReadAll(r.Body)
			if err != nil {
				handleErr(w, err, "POST userHandler")
				return
			}
			var u User
			err = json.Unmarshal(bs, &u)
			if err != nil {
				handleErr(w, err, "POST userHandler")
				return
			}
			newid, err := createUser(db, &u)
			if err != nil {
				handleErr(w, err, "POST userHandler")
				return
			}
			u.Userid = newid

			w.Header().Set("Content-Type", "application/json")
			P := makeFprintf(w)
			P("%s", jsonstr(u))
			return
		} else if r.Method == "PUT" {
			bs, err := ioutil.ReadAll(r.Body)
			if err != nil {
				handleErr(w, err, "PUT userHandler")
				return
			}
			var u User
			err = json.Unmarshal(bs, &u)
			if err != nil {
				handleErr(w, err, "PUT userHandler")
				return
			}
			// Only admin or same user can PUT
			if requser.Userid != 1 && u.Userid != requser.Userid {
				http.Error(w, "Invalid user", 401)
				return
			}

			err = editUser(db, &u)
			if err != nil {
				handleErr(w, err, "PUT userHandler")
				return
			}

			w.Header().Set("Content-Type", "application/json")
			P := makeFprintf(w)
			P("%s", jsonstr(u))
			return
		} else if r.Method == "DELETE" {
			qid := idtoi(r.FormValue("id"))
			if qid == 0 {
				http.Error(w, "Not found.", 404)
				return
			}
			// Only admin or same user can DEL
			if requser.Userid != 1 && qid != requser.Userid {
				http.Error(w, "Invalid user", 401)
				return
			}

			err := delUser(db, qid)
			if err != nil {
				handleErr(w, err, "DEL userHandler")
				return
			}

			return
		}

		http.Error(w, "Use GET/POST/PUT/DELETE", 401)
	}
}
