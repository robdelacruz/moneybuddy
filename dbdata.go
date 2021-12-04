package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type AccountType int64

const (
	BankAccount AccountType = iota
	StockAccount
)

type BookType int

const (
	UserBook BookType = iota
	SystemBook
)

type User struct {
	Userid   int64  `json:"userid"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type Currency struct {
	Currencyid int64   `json:"currencyid"`
	Currency   string  `json:"currency"`
	Usdrate    float64 `json:"usdrate"`
	Userid     int64   `json:"userid"`
}
type Rootdata struct {
	Userid     int64       `json:"userid"`
	Currencies []*Currency `json:"currencies"`
	Books      []*Book     `json:"books"`
}
type Book struct {
	Bookid        int64      `json:"bookid"`
	Name          string     `json:"name"`
	BookType      BookType   `json:"booktype"`
	BankAccounts  []*Account `json:"bankaccounts"`
	StockAccounts []*Account `json:"stockaccounts"`
	Userid        int64      `json:"userid"`
	Active        int64      `json:"active"`
}
type Account struct {
	Accountid   int64       `json:"accountid"`
	Bookid      int64       `json:"bookid"`
	Code        string      `json:"code"`
	Seq         int64       `json:"seq"`
	Name        string      `json:"name"`
	AccountType AccountType `json:"accounttype"`
	Unitprice   float64     `json:"unitprice"`
	Currencyid  int64       `json:"currencyid"`
	Ref         string      `json:"ref"`
	Memo        string      `json:"memo"`
	Currency    *Currency   `json:"currency"`
	Balance     float64     `json:"balance"`
	Txns        []*Txn      `json:"txns"`
}
type Txn struct {
	Txnid     int64   `json:"txnid"`
	Accountid int64   `json:"accountid"`
	Date      string  `json:"date"`
	Ref       string  `json:"ref"`
	Desc      string  `json:"desc"`
	Amt       float64 `json:"amt"`
	Memo      string  `json:"memo"`
}

func createTables(newfile string) *sql.DB {
	if fileExists(newfile) {
		s := fmt.Sprintf("File '%s' already exists. Can't initialize it.\n", newfile)
		fmt.Printf(s)
		return nil
	}

	db, err := sql.Open("sqlite3", newfile)
	if err != nil {
		fmt.Printf("Error opening '%s' (%s)\n", newfile, err)
		return nil
	}

	ss := []string{
		"CREATE TABLE user (user_id INTEGER PRIMARY KEY NOT NULL, username TEXT UNIQUE, password TEXT);",
		"CREATE TABLE book (book_id INTEGER PRIMARY KEY NOT NULL, name TEXT NOT NULL DEFAULT 'My Accounts', booktype INTEGER NOT NULL, user_id INTEGER NOT NULL, active INTEGER NOT NULL DEFAULT 1);",
		"CREATE TABLE currency (currency_id INTEGER PRIMARY KEY NOT NULL, currency TEXT NOT NULL, usdrate REAL NOT NULL DEFAULT 1.0, user_id INTEGER NOT NULL);",
		"CREATE TABLE account (account_id INTEGER PRIMARY KEY NOT NULL, book_id INTEGER NOT NULL, code TEXT DEFAULT '', seq INTEGER NOT NULL DEFAULT 0, name TEXT NOT NULL DEFAULT 'account', accounttype INTEGER NOT NULL, currency_id INTEGER NOT NULL, unitprice REAL NOT NULL DEFAULT 1.0, ref TEXT NOT NULL DEFAULT '', memo TEXT NOT NULL DEFAULT '');",
		"CREATE TABLE txn (txn_id INTEGER PRIMARY KEY NOT NULL, account_id INTEGER NOT NULL, date TEXT NOT NULL DEFAULT '', ref TEXT NOT NULL DEFAULT '', desc TEXT NOT NULL DEFAULT '', amt REAL NOT NULL DEFAULT 0.0, memo TEXT NOT NULL DEFAULT '');",
		"INSERT INTO user (user_id, username, password) VALUES (1, 'admin', '');",
	}

	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("DB error (%s)\n", err)
		return nil
	}
	for _, s := range ss {
		_, err := txexec(tx, s)
		if err != nil {
			tx.Rollback()
			fmt.Printf("DB error (%s)\n", err)
			return nil
		}
	}
	err = tx.Commit()
	if err != nil {
		fmt.Printf("DB error (%s)\n", err)
		return nil
	}

	return db
}

func findAccountUserid(db *sql.DB, accountid int64) (int64, error) {
	s := `SELECT IFNULL(b.user_id, 0) 
FROM account a
INNER JOIN book b ON b.book_id = a.book_id 
WHERE a.account_id = ?`
	rows, err := db.Query(s, accountid)
	if err != nil {
		return 0, err
	}
	var userid int64
	for rows.Next() {
		rows.Scan(&userid)
		break
	}
	rows.Close()
	return userid, nil
}

func findTxnUserid(db *sql.DB, txnid int64) (int64, error) {
	s := `SELECT IFNULL(b.user_id, 0) 
FROM txn t
INNER JOIN account a ON a.account_id = t.account_id
INNER JOIN book b ON b.book_id = a.book_id 
WHERE t.txn_id = ?`
	rows, err := db.Query(s, txnid)
	if err != nil {
		return 0, err
	}
	var userid int64
	for rows.Next() {
		rows.Scan(&userid)
		break
	}
	rows.Close()
	return userid, nil
}

//** Currency functions **
func createCurrency(db *sql.DB, c *Currency) (int64, error) {
	s := "INSERT INTO currency (currency, usdrate, user_id) VALUES (?, ?, ?)"
	result, err := sqlexec(db, s, c.Currency, c.Usdrate, c.Userid)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
func editCurrency(db *sql.DB, c *Currency) error {
	s := "UPDATE currency SET currency = ?, usdrate = ?, user_id = ? WHERE currency_id = ?"
	_, err := sqlexec(db, s, c.Currency, c.Usdrate, c.Userid, c.Currencyid)
	if err != nil {
		return err
	}
	return nil
}
func delCurrency(db *sql.DB, currencyid int64) error {
	s := "DELETE FROM currency WHERE currency_id = ?"
	_, err := sqlexec(db, s, currencyid)
	if err != nil {
		return err
	}
	return nil
}

func findCurrency(db *sql.DB, currencyid int64) (*Currency, error) {
	s := "SELECT currency_id, currency, usdrate, user_id FROM currency WHERE currency_id = ?"
	row := db.QueryRow(s, currencyid)
	var c Currency
	err := row.Scan(&c.Currencyid, &c.Currency, &c.Usdrate, &c.Userid)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
func findCurrencies(db *sql.DB, userid int64, swhere string) ([]*Currency, error) {
	s := fmt.Sprintf("SELECT currency_id, currency, usdrate, user_id FROM currency WHERE user_id = ? AND %s", swhere)
	rows, err := db.Query(s, userid)
	if err != nil {
		return nil, err
	}
	cc := []*Currency{}
	for rows.Next() {
		var c Currency
		rows.Scan(&c.Currencyid, &c.Currency, &c.Usdrate, &c.Userid)
		cc = append(cc, &c)
	}
	return cc, nil
}
func findUserCurrencies(db *sql.DB, userid int64) ([]*Currency, error) {
	return findCurrencies(db, userid, "1=1 ORDER BY currency_id")
}

func convToCurrency(n float64, ncur, tocur *Currency) float64 {
	if ncur.Usdrate == 0.0 {
		return n * tocur.Usdrate
	}
	return n / ncur.Usdrate * tocur.Usdrate
}

//** Rootdata, Book functions
func findRootdata(db *sql.DB, userid int64) (*Rootdata, error) {
	cc, err := findUserCurrencies(db, userid)
	if err != nil {
		return nil, err
	}
	bb, err := findUserBooks(db, userid)
	if err != nil {
		return nil, err
	}

	var rd Rootdata
	rd.Userid = userid
	rd.Currencies = cc
	rd.Books = bb
	return &rd, nil
}

func createBook(db *sql.DB, b *Book) (int64, error) {
	s := "INSERT INTO book (name, booktype, user_id, active) VALUES (?, ?, ?, ?)"
	result, err := sqlexec(db, s, b.Name, b.BookType, b.Userid, b.Active)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
func editBook(db *sql.DB, b *Book) error {
	s := "UPDATE book SET name = ?, booktype = ?, user_id = ?, active = ? WHERE book_id = ?"
	_, err := sqlexec(db, s, b.Name, b.BookType, b.Userid, b.Active, b.Bookid)
	if err != nil {
		return err
	}
	return nil
}
func delBook(db *sql.DB, bookid int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Delete book's accounts' transactions.
	s := "DELETE FROM txn WHERE account_id IN (SELECT account_id FROM account WHERE book_id = ?)"
	_, err = txexec(tx, s, bookid)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete this book's accounts.
	s = "DELETE FROM account WHERE book_id = ?"
	_, err = txexec(tx, s, bookid)
	if err != nil {
		tx.Rollback()
		return err
	}

	s = "DELETE FROM book WHERE book_id = ?"
	_, err = txexec(tx, s, bookid)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func assignBookAccounts(b *Book, aa []*Account) {
	bb := []*Account{}
	ss := []*Account{}
	for _, a := range aa {
		if a.AccountType == BankAccount {
			bb = append(bb, a)
		} else if a.AccountType == StockAccount {
			ss = append(ss, a)
		}
	}
	b.BankAccounts = bb
	b.StockAccounts = ss
}

func findBook(db *sql.DB, bookid int64) (*Book, error) {
	s := "SELECT book_id, name, booktype, user_id, active FROM book WHERE book_id = ?"
	row := db.QueryRow(s, bookid)
	var b Book
	err := row.Scan(&b.Bookid, &b.Name, &b.BookType, &b.Userid, &b.Active)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	aa, err := findAllAccountsByType(db, bookid)
	if err != nil {
		return nil, err
	}

	assignBookAccounts(&b, aa)
	return &b, nil
}

func findBooks(db *sql.DB, userid int64, swhere string) ([]*Book, error) {
	s := fmt.Sprintf("SELECT book_id, name, booktype, user_id, active FROM book WHERE user_id = ? AND %s", swhere)
	rows, err := db.Query(s, userid)
	if err != nil {
		return nil, err
	}
	bb := []*Book{}
	for rows.Next() {
		var b Book
		rows.Scan(&b.Bookid, &b.Name, &b.BookType, &b.Userid, &b.Active)

		aa, err := findAllAccountsByType(db, b.Bookid)
		if err != nil {
			return nil, err
		}
		assignBookAccounts(&b, aa)
		bb = append(bb, &b)
	}
	return bb, nil
}
func findUserBooks(db *sql.DB, userid int64) ([]*Book, error) {
	return findBooks(db, userid, "1=1 ORDER BY booktype, active, book_id")
}

//** Account functions **
func createAccount(db *sql.DB, a *Account, bookid int64) (int64, error) {
	s := fmt.Sprintf("INSERT INTO account (book_id, code, name, accounttype, currency_id, unitprice, ref, memo, seq) VALUES (?, ?, ?, ?, ?, ?, ?, ?, (SELECT IFNULL(MAX(seq), 0)+1 FROM account WHERE accounttype = %d))", a.AccountType)
	result, err := sqlexec(db, s, bookid, a.Code, a.Name, a.AccountType, a.Currencyid, a.Unitprice, a.Ref, a.Memo)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
func editAccount(db *sql.DB, a *Account) error {
	s := "UPDATE account SET book_id = ?, code = ?, seq = ?, name = ?, accounttype = ?, currency_id = ?, unitprice = ?, ref = ?, memo = ? WHERE account_id = ?"
	_, err := sqlexec(db, s, a.Bookid, a.Code, a.Seq, a.Name, a.AccountType, a.Currencyid, a.Unitprice, a.Ref, a.Memo, a.Accountid)
	if err != nil {
		return err
	}
	return nil
}
func delAccount(db *sql.DB, accountid int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Delete account's transactions.
	s := "DELETE FROM txn WHERE account_id = ?"
	_, err = txexec(tx, s, accountid)
	if err != nil {
		tx.Rollback()
		return err
	}

	s = "DELETE FROM account WHERE account_id = ?"
	_, err = txexec(tx, s, accountid)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
func assignAccountToBook(db *sql.DB, accountid, bookid int64) error {
	b, err := findBook(db, bookid)
	if err != nil {
		return err
	}
	if b == nil {
		return sql.ErrNoRows
	}

	s := "UPDATE account SET book_id = ? WHERE account_id = ?"
	_, err = sqlexec(db, s, bookid, accountid)
	if err != nil {
		return err
	}
	return nil
}

func findAccount(db *sql.DB, accountid int64) (*Account, error) {
	// account balance is calculated this way:
	// For accounttype = 0 (bankaccount):
	//   total of all txn amounts
	// For accounttype = 1 (stockaccount):
	//   total of all txn shares * acccount unit price

	// txn shares is recorded in txn.amount field

	s := `SELECT account_id, book_id, code, seq, name, accounttype, a.unitprice, a.currency_id, a.ref, a.memo, IFNULL(cur.currency, ''), IFNULL(cur.Usdrate, 1.0), 
(SELECT IIF(a.accounttype = 0, IFNULL(SUM(txn.amt), 0.0), IFNULL(SUM(txn.amt)*a.unitprice, 0.0))
  FROM txn WHERE txn.account_id = a.account_id) AS bal
FROM account a 
LEFT OUTER JOIN currency cur ON cur.currency_id = a.currency_id 
WHERE account_id = ?`
	row := db.QueryRow(s, accountid)
	var a Account
	var c Currency
	err := row.Scan(&a.Accountid, &a.Bookid, &a.Code, &a.Seq, &a.Name, &a.AccountType, &a.Unitprice, &c.Currencyid, &a.Ref, &a.Memo, &c.Currency, &c.Usdrate, &a.Balance)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	a.Currencyid = c.Currencyid
	a.Currency = &c

	tt, err := findAllTxnsOfAccount(db, accountid)
	if err != nil {
		return nil, err
	}
	a.Txns = tt
	return &a, nil
}
func findAccounts(db *sql.DB, bookid int64, swhere string) ([]*Account, error) {
	s := fmt.Sprintf(`
SELECT a.account_id, a.book_id, a.code, a.seq, a.name, a.accounttype, a.unitprice, a.currency_id, a.ref, a.memo, IFNULL(cur.currency, ''), IFNULL(cur.Usdrate, 1.0),
(SELECT IIF(a.accounttype = 0, IFNULL(SUM(txn.amt), 0.0), IFNULL(SUM(txn.amt)*a.unitprice, 0.0))
  FROM txn WHERE txn.account_id = a.account_id) AS bal
FROM account a 
LEFT OUTER JOIN currency cur ON cur.currency_id = a.currency_id 
WHERE a.book_id = ? AND %s`, swhere)
	rows, err := db.Query(s, bookid)
	if err != nil {
		return nil, err
	}
	aa := []*Account{}
	for rows.Next() {
		var a Account
		var c Currency
		rows.Scan(&a.Accountid, &a.Bookid, &a.Code, &a.Seq, &a.Name, &a.AccountType, &a.Unitprice, &c.Currencyid, &a.Ref, &a.Memo, &c.Currency, &c.Usdrate, &a.Balance)
		a.Currencyid = c.Currencyid
		a.Currency = &c

		tt, err := findAllTxnsOfAccount(db, a.Accountid)
		if err != nil {
			return nil, err
		}
		a.Txns = tt
		aa = append(aa, &a)
	}
	return aa, nil
}
func findAllAccounts(db *sql.DB, bookid int64) ([]*Account, error) {
	return findAccounts(db, bookid, "1=1 ORDER BY name")
}
func findAllAccountsByType(db *sql.DB, bookid int64) ([]*Account, error) {
	return findAccounts(db, bookid, "1=1 ORDER BY accounttype, seq, name")
}

func accountSumAmt(db *sql.DB, accountid int64) (float64, error) {
	s := "SELECT IFNULL(SUM(amt), 0.0) FROM txn WHERE account_id = ?"
	row := db.QueryRow(s, accountid)
	var bal float64
	err := row.Scan(&bal)
	if err != nil {
		return 0.0, err
	}
	return bal, nil
}

// Move accountid to seq:newseq and shift right all accounts >= newseq
func resequenceAccounts(db *sql.DB, accounttype AccountType, accountid, newseq int64) error {
	// Shift all accounts on or after newseq to the right by 1 seq.
	s := fmt.Sprintf(`
SELECT a.account_id 
FROM account a 
WHERE a.book_id = (SELECT book_id FROM account WHERE account_id = %d) 
AND a.accounttype = ? AND a.seq >= ? 
ORDER BY a.seq`, accountid)
	rows, err := db.Query(s, accounttype, newseq)
	if err != nil {
		return err
	}
	ids := []int64{}
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}
	rightseq := newseq + 1
	for _, id := range ids {
		s = "UPDATE account SET seq = ? WHERE account_id = ?"
		_, err = sqlexec(db, s, rightseq, id)
		if err != nil {
			return err
		}
		rightseq++
	}

	// Move account into position.
	s = "UPDATE account SET seq = ? WHERE account_id = ?"
	_, err = sqlexec(db, s, newseq, accountid)
	if err != nil {
		return err
	}
	return nil
}

func createRandomBankAccount(db *sql.DB, bookid int64, currencyid int64) (int64, error) {
	banks := []string{"BPI", "Security", "Wells Fargo", "Bank of America", "FirstTech", "MetroBank"}
	descs := []string{"Savings", "Checking", "Time Deposit", "Money Market", "Individual", "COD", "Maxi"}
	opts := []string{"", "Cash", "Ext", "Other"}

	ibank := rand.Intn(len(banks))
	idesc := rand.Intn(len(descs))
	iopt := rand.Intn(len(opts))
	name := strings.TrimSpace(fmt.Sprintf("%s %s %s", banks[ibank], descs[idesc], opts[iopt]))

	a := Account{
		Bookid:      bookid,
		Code:        "",
		Name:        name,
		AccountType: BankAccount,
		Unitprice:   1.0,
		Currencyid:  currencyid,
		Ref:         strconv.Itoa(rand.Intn(999999999999999)),
		Memo:        "",
	}
	accountid, err := createAccount(db, &a, bookid)
	if err != nil {
		return 0, err
	}

	// Generate transactions for this account.
	ntxns := 5 + rand.Intn(150)
	err = createRandomTxns(db, accountid, ntxns)
	if err != nil {
		return 0, err
	}

	return accountid, nil
}

func createRandomStockAccount(db *sql.DB, bookid int64, ticker string, unitprice float64, currencyid int64) (int64, error) {
	a := Account{
		Bookid:      bookid,
		Code:        "",
		Name:        ticker,
		AccountType: StockAccount,
		Unitprice:   unitprice,
		Currencyid:  currencyid,
		Ref:         strconv.Itoa(rand.Intn(999999999999999)),
		Memo:        "",
	}
	accountid, err := createAccount(db, &a, bookid)
	if err != nil {
		return 0, err
	}

	// Generate transactions for this account.
	ntxns := 5 + rand.Intn(30)
	err = createRandomStockTxns(db, accountid, ntxns)
	if err != nil {
		return 0, err
	}

	return accountid, nil
}

//** Txn functions **
func createTxn(db *sql.DB, t *Txn) (int64, error) {
	s := "INSERT INTO txn (account_id, date, ref, desc, amt, memo) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := sqlexec(db, s, t.Accountid, t.Date, t.Ref, t.Desc, t.Amt, t.Memo)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
func editTxn(db *sql.DB, t *Txn) error {
	s := "UPDATE txn SET date = ?, ref = ?, desc = ?, amt = ?, memo = ? WHERE txn_id = ?"
	_, err := sqlexec(db, s, t.Date, t.Ref, t.Desc, t.Amt, t.Memo, t.Txnid)
	if err != nil {
		return err
	}
	return nil
}
func delTxn(db *sql.DB, txnid int64) error {
	s := "DELETE FROM txn WHERE txn_id = ?"
	_, err := sqlexec(db, s, txnid)
	if err != nil {
		return err
	}
	return nil
}

func findTxn(db *sql.DB, txnid int64) (*Txn, error) {
	s := "SELECT txn_id, account_id, date, ref, desc, amt, memo FROM txn WHERE txn_id = ?"
	row := db.QueryRow(s, txnid)
	var t Txn
	err := row.Scan(&t.Txnid, &t.Accountid, &t.Date, &t.Ref, &t.Desc, &t.Amt, &t.Memo)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}
func findTxns(db *sql.DB, swhere string) ([]*Txn, error) {
	s := fmt.Sprintf("SELECT txn_id, account_id, date, ref, desc, amt, memo FROM txn WHERE %s", swhere)
	rows, err := db.Query(s)
	if err != nil {
		return nil, err
	}
	tt := []*Txn{}
	for rows.Next() {
		var t Txn
		rows.Scan(&t.Txnid, &t.Accountid, &t.Date, &t.Ref, &t.Desc, &t.Amt, &t.Memo)
		tt = append(tt, &t)
	}
	return tt, nil
}
func findTxnsOfAccount(db *sql.DB, accountid int64, swhere string) ([]*Txn, error) {
	s := fmt.Sprintf("SELECT txn_id, account_id, date, ref, desc, amt, memo FROM txn WHERE account_id = ? AND %s", swhere)
	rows, err := db.Query(s, accountid)
	if err != nil {
		return nil, err
	}
	tt := []*Txn{}
	for rows.Next() {
		var t Txn
		rows.Scan(&t.Txnid, &t.Accountid, &t.Date, &t.Ref, &t.Desc, &t.Amt, &t.Memo)
		tt = append(tt, &t)
	}
	return tt, nil
}

func findAllTxnsOfAccount(db *sql.DB, accountid int64) ([]*Txn, error) {
	return findTxnsOfAccount(db, accountid, "1=1 ORDER BY date DESC")
}

func createRandomWords(words []string) string {
	nwords := 1 + rand.Intn(11)

	var sb strings.Builder
	for i := 0; i < nwords; i++ {
		sb.WriteString(words[rand.Intn(len(words))])
		if i < nwords-1 {
			sb.WriteString(" ")
		}
	}

	return sb.String()
}

func createRandomTxns(db *sql.DB, accountid int64, ntxns int) error {
	words := []string{"interest", "dividend", "refund", "salary", "deposit from", "receive from", "withdraw", "banking fee", "tenant", "bank", "pay"}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	s := "INSERT INTO txn (account_id, date, ref, desc, amt, memo) VALUES (?, ?, ?, ?, ?, ?)"

	for i := 0; i < ntxns; i++ {
		amt := float64(rand.Intn(5000000))/100.0 - 25000
		t := Txn{
			Accountid: accountid,
			Date:      isodate(randdate(2000, 2021)),
			Ref:       "",
			Desc:      createRandomWords(words),
			Amt:       amt,
			Memo:      createRandomWords(words),
		}
		_, err := txexec(tx, s, t.Accountid, t.Date, t.Ref, t.Desc, t.Amt, t.Memo)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func createRandomStockTxns(db *sql.DB, accountid int64, ntxns int) error {
	buydescs := []string{"buy", "dividend"}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	s := "INSERT INTO txn (account_id, date, ref, desc, amt, memo) VALUES (?, ?, ?, ?, ?, ?)"

	for i := 0; i < ntxns; i++ {
		desc := "sell"
		amt := float64(rand.Intn(50000))/100.0 - 250
		if amt >= 0 {
			desc = buydescs[rand.Intn(len(buydescs))]
		}

		t := Txn{
			Accountid: accountid,
			Date:      isodate(randdate(2000, 2021)),
			Ref:       "",
			Desc:      desc,
			Amt:       amt,
			Memo:      "",
		}
		_, err := txexec(tx, s, t.Accountid, t.Date, t.Ref, t.Desc, t.Amt, t.Memo)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

//** User functions **
func createUser(db *sql.DB, u *User) (int64, error) {
	s := "INSERT INTO user (username, password) VALUES (?, ?)"
	result, err := sqlexec(db, s, u.Username, u.Password)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
func editUser(db *sql.DB, u *User) error {
	s := "UPDATE user SET username = ?, password = ? WHERE user_id = ?"
	_, err := sqlexec(db, s, u.Username, u.Password, u.Userid)
	if err != nil {
		return err
	}
	return nil
}
func delUser(db *sql.DB, userid int64) error {
	s := "DELETE FROM user WHERE user_id = ?"
	_, err := sqlexec(db, s, userid)
	if err != nil {
		return err
	}
	return nil
}

func findUser(db *sql.DB, userid int64) (*User, error) {
	s := "SELECT user_id, username, password FROM user WHERE user_id = ?"
	row := db.QueryRow(s, userid)
	var u User
	err := row.Scan(&u.Userid, &u.Username, &u.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func findUserByUsername(db *sql.DB, username string) (*User, error) {
	s := "SELECT user_id, username, password FROM user WHERE username = ?"
	row := db.QueryRow(s, username)
	var u User
	err := row.Scan(&u.Userid, &u.Username, &u.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func isUsernameExists(db *sql.DB, username string) (bool, error) {
	u, err := findUserByUsername(db, username)
	if err != nil {
		return false, err
	}
	if u != nil {
		return true, nil
	}
	return false, nil
}

func findUsers(db *sql.DB, swhere string) ([]*User, error) {
	s := fmt.Sprintf("SELECT user_id, username, password FROM user WHERE %s", swhere)
	rows, err := db.Query(s)
	if err != nil {
		return nil, err
	}
	uu := []*User{}
	for rows.Next() {
		var u User
		rows.Scan(&u.Userid, &u.Username, &u.Password)
		uu = append(uu, &u)
	}
	return uu, nil
}
func findAllUsers(db *sql.DB) ([]*User, error) {
	return findUsers(db, "1=1 ORDER BY user_id")
}
