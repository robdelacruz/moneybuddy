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
	subs []SignalChan
	mu   sync.RWMutex
	whos []string
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
	http.HandleFunc("/api/book", bookHandler(db))
	http.HandleFunc("/api/account", accountHandler(db))
	http.HandleFunc("/api/txn", txnHandler(db))
	http.HandleFunc("/api/subscriberoot", subscriberootHandler(db))
	http.HandleFunc("/api/whos", whosHandler(db))
	http.HandleFunc("/api/rptdata", rptdataHandler(db))

	http.HandleFunc("/api/login", loginHandler(db))
	http.HandleFunc("/api/signup", signupHandler(db))
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
		"CREATE TABLE user (user_id INTEGER PRIMARY KEY NOT NULL, username TEXT UNIQUE, password TEXT);",
		"CREATE TABLE book (book_id INTEGER PRIMARY KEY NOT NULL, name TEXT NOT NULL DEFAULT 'My Accounts', user_id INTEGER NOT NULL);",
		"CREATE TABLE currency (currency_id INTEGER PRIMARY KEY NOT NULL, currency TEXT NOT NULL, usdrate REAL NOT NULL DEFAULT 1.0, user_id INTEGER NOT NULL);",
		"CREATE TABLE account (account_id INTEGER PRIMARY KEY NOT NULL, code TEXT DEFAULT '', name TEXT NOT NULL DEFAULT 'account', accounttype INTEGER NOT NULL, currency_id INTEGER NOT NULL, unitprice REAL NOT NULL DEFAULT 1.0);",
		"CREATE TABLE bookaccount (book_id INTEGER NOT NULL, account_id INTEGER NOT NULL);",
		"CREATE TABLE txn (txn_id INTEGER PRIMARY KEY NOT NULL, account_id INTEGER NOT NULL, date TEXT NOT NULL DEFAULT '', ref TEXT NOT NULL DEFAULT '', desc TEXT NOT NULL DEFAULT '', amt REAL NOT NULL DEFAULT 0.0, memo TEXT NOT NULL DEFAULT '');",
		"INSERT INTO user (user_id, username, password) VALUES (1, 'admin', '');",
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

	fmt.Printf("Creating user1's test data... ")
	initTestData(db, "rob")
	fmt.Printf("Creating user2's test data... ")
	initTestData(db, "user2")
	fmt.Printf("Done\n")
}

func initTestData(db *sql.DB, username string) {
	u := User{
		Username: username,
	}
	userid, err := createUser(db, &u)
	if err != nil {
		panic(err)
	}

	c := Currency{
		Currency: "USD",
		Usdrate:  1.0,
		Userid:   userid,
	}
	_, err = createCurrency(db, &c)
	if err != nil {
		panic(err)
	}
	c = Currency{
		Currency: "PHP",
		Usdrate:  48.0,
		Userid:   userid,
	}
	_, err = createCurrency(db, &c)
	if err != nil {
		panic(err)
	}

	b := Book{
		Name:   fmt.Sprintf("%s's accounts", username),
		Userid: userid,
	}
	_, err = createBook(db, &b)
	if err != nil {
		panic(err)
	}
	b = Book{
		Name:   "Work Accounts",
		Userid: userid,
	}
	_, err = createBook(db, &b)
	if err != nil {
		panic(err)
	}

	bb, err := findUserBooks(db, userid)
	if err != nil {
		panic(err)
	}

	for _, b := range bb {
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
				handleErr(w, err, "bookHandler")
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
			signalAndCloseSubs(&_datasync)

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
			signalAndCloseSubs(&_datasync)

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
			signalAndCloseSubs(&_datasync)

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
			newid, err := createAccount(db, &a, qbookid)
			if err != nil {
				handleErr(w, err, "POST accountHandler")
				return
			}
			a.Accountid = newid

			// Inform all data subscribers that a data change occured.
			signalAndCloseSubs(&_datasync)

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
			signalAndCloseSubs(&_datasync)

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
			err := delAccount(db, qid)
			if err != nil {
				handleErr(w, err, "DEL accountHandler")
				return
			}

			// Inform all data subscribers that a data change occured.
			signalAndCloseSubs(&_datasync)

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

			// Inform all data subscribers that a data change occured.
			signalAndCloseSubs(&_datasync)

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
			signalAndCloseSubs(&_datasync)

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

			// Inform all data subscribers that a data change occured.
			signalAndCloseSubs(&_datasync)

			return
		}

		http.Error(w, "Use GET/POST/PUT/DELETE", 401)
	}
}

func whosHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		P := makeFprintf(w)
		P("num subs: %d\n", len(_datasync.subs))
		P("num whos: %d\n", len(_datasync.whos))

		for i := 0; i < len(_datasync.whos); i++ {
			P("whos[%d]: %s\n", i, _datasync.whos[i])
		}
	}
}

func rptdataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		quserid := idtoi(r.FormValue("userid"))
		if quserid == 0 {
			http.Error(w, "Not found.", 404)
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
		quserid := idtoi(r.FormValue("userid"))
		if quserid == 0 {
			http.Error(w, "Not found.", 404)
			return
		}

		qwho := r.FormValue("who")
		if qwho == "" {
			qwho = "(noname)"
		}
		_datasync.whos = append(_datasync.whos, qwho)

		sub := make(SignalChan, 1)
		addSub(&_datasync, sub)
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

func addSub(ds *DataSync, sub SignalChan) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.subs = append(ds.subs, sub)
}
func signalSubs(ds *DataSync) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	for _, ch := range ds.subs {
		ch <- struct{}{}
	}
}
func closeSubs(ds *DataSync) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	for _, sub := range ds.subs {
		close(sub)
	}
	ds.subs = nil
	ds.whos = nil
}
func signalAndCloseSubs(ds *DataSync) {
	signalSubs(ds)
	closeSubs(ds)
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

func userHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			qid := idtoi(r.FormValue("id"))
			if qid == 0 {
				http.Error(w, "Not found.", 404)
				return
			}
			a, err := findUser(db, qid)
			if err != nil {
				handleErr(w, err, "userHandler")
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
