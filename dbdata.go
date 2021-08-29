package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type AccountType int

const (
	BankAccount AccountType = iota
	StockAccount
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
}
type Rootdata struct {
	Currencies []*Currency `json:"currencies"`
	Books      []*Book     `json:"books"`
}
type Book struct {
	Bookid        int64      `json:"bookid"`
	Name          string     `json:"name"`
	BankAccounts  []*Account `json:"bankaccounts"`
	StockAccounts []*Account `json:"stockaccounts"`
	Userid        int64      `json:"userid"`
}
type Account struct {
	Accountid   int64       `json:"accountid"`
	Code        string      `json:"code"`
	Name        string      `json:"name"`
	AccountType AccountType `json:"accounttype"`
	Unitprice   float64     `json:"unitprice"`
	Currencyid  int64       `json:"currencyid"`
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
}

//** Currency functions **
func createCurrency(db *sql.DB, c *Currency) (int64, error) {
	s := "INSERT INTO currency (currency, usdrate) VALUES (?, ?)"
	result, err := sqlexec(db, s, c.Currency, c.Usdrate)
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
	s := "UPDATE currency SET currency = ?, usdrate = ? WHERE currency_id = ?"
	_, err := sqlexec(db, s, c.Currency, c.Usdrate, c.Currencyid)
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
	s := "SELECT currency_id, currency, usdrate FROM currency WHERE currency_id = ?"
	row := db.QueryRow(s, currencyid)
	var c Currency
	err := row.Scan(&c.Currencyid, &c.Currency, &c.Usdrate)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
func findCurrencies(db *sql.DB, swhere string) ([]*Currency, error) {
	s := fmt.Sprintf("SELECT currency_id, currency, usdrate FROM currency WHERE %s", swhere)
	rows, err := db.Query(s)
	if err != nil {
		return nil, err
	}
	cc := []*Currency{}
	for rows.Next() {
		var c Currency
		rows.Scan(&c.Currencyid, &c.Currency, &c.Usdrate)
		cc = append(cc, &c)
	}
	return cc, nil
}
func findAllCurrencies(db *sql.DB) ([]*Currency, error) {
	return findCurrencies(db, "1=1 ORDER BY currency_id")
}

func convToCurrency(n float64, ncur, tocur *Currency) float64 {
	if ncur.Usdrate == 0.0 {
		return n * tocur.Usdrate
	}
	return n / ncur.Usdrate * tocur.Usdrate
}

//** Rootdata, Book functions
func findRootdata(db *sql.DB) (*Rootdata, error) {
	cc, err := findAllCurrencies(db)
	if err != nil {
		return nil, err
	}
	bb, err := findAllBooks(db)
	if err != nil {
		return nil, err
	}

	var rd Rootdata
	rd.Currencies = cc
	rd.Books = bb
	return &rd, nil
}

func createBook(db *sql.DB, b *Book) (int64, error) {
	s := "INSERT INTO book (name) VALUES (?)"
	result, err := sqlexec(db, s, b.Name)
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
	s := "UPDATE book SET name = ? WHERE book_id = ?"
	_, err := sqlexec(db, s, b.Name, b.Bookid)
	if err != nil {
		return err
	}
	return nil
}
func delBook(db *sql.DB, bookid int64) error {
	s := "DELETE FROM book WHERE book_id = ?"
	_, err := sqlexec(db, s, bookid)
	if err != nil {
		return err
	}
	return nil
}

func assignBookAccounts(b *Book, aa []*Account) {
	var bb []*Account
	var ss []*Account
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
	s := "SELECT book_id, name FROM book WHERE book_id = ?"
	row := db.QueryRow(s, bookid)
	var b Book
	err := row.Scan(&b.Bookid, &b.Name)
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

func findBooks(db *sql.DB, swhere string) ([]*Book, error) {
	s := fmt.Sprintf("SELECT book_id, name FROM book WHERE %s", swhere)
	rows, err := db.Query(s)
	if err != nil {
		return nil, err
	}
	bb := []*Book{}
	for rows.Next() {
		var b Book
		rows.Scan(&b.Bookid, &b.Name)

		aa, err := findAllAccountsByType(db, b.Bookid)
		if err != nil {
			return nil, err
		}
		assignBookAccounts(&b, aa)
		bb = append(bb, &b)
	}
	return bb, nil
}
func findAllBooks(db *sql.DB) ([]*Book, error) {
	return findBooks(db, "1=1 ORDER BY book_id")
}

//** Account functions **
func createAccount(db *sql.DB, a *Account, bookid int64) (int64, error) {
	s := "INSERT INTO account (code, name, accounttype, currency_id, unitprice) VALUES (?, ?, ?, ?, ?)"
	result, err := sqlexec(db, s, a.Code, a.Name, a.AccountType, a.Currencyid, a.Unitprice)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	s = "INSERT INTO bookaccount (book_id, account_id) VALUES (?, ?)"
	_, err = sqlexec(db, s, bookid, id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
func editAccount(db *sql.DB, a *Account) error {
	s := "UPDATE account SET code = ?, name = ?, accounttype = ?, currency_id = ?, unitprice = ? WHERE account_id = ?"
	_, err := sqlexec(db, s, a.Code, a.Name, a.AccountType, a.Currencyid, a.Unitprice, a.Accountid)
	if err != nil {
		return err
	}
	return nil
}
func delAccount(db *sql.DB, accountid int64) error {
	s := "DELETE FROM account WHERE account_id = ?"
	_, err := sqlexec(db, s, accountid)
	if err != nil {
		return err
	}

	s = "DELETE FROM bookaccount where account_id = ?"
	_, err = sqlexec(db, s, accountid)
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

	s := `SELECT account_id, code, name, accounttype, a.unitprice, a.currency_id, IFNULL(cur.currency, ''), IFNULL(cur.Usdrate, 1.0), 
(SELECT IIF(a.accounttype = 0, IFNULL(SUM(txn.amt), 0.0), IFNULL(SUM(txn.amt)*a.unitprice, 0.0))
  FROM txn WHERE txn.account_id = a.account_id) AS bal
FROM account a 
LEFT OUTER JOIN currency cur ON cur.currency_id = a.currency_id 
WHERE account_id = ?`
	row := db.QueryRow(s, accountid)
	var a Account
	var c Currency
	err := row.Scan(&a.Accountid, &a.Code, &a.Name, &a.AccountType, &a.Unitprice, &c.Currencyid, &c.Currency, &c.Usdrate, &a.Balance)
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
SELECT a.account_id, a.code, a.name, a.accounttype, a.unitprice, a.currency_id, IFNULL(cur.currency, ''), IFNULL(cur.Usdrate, 1.0),
(SELECT IIF(a.accounttype = 0, IFNULL(SUM(txn.amt), 0.0), IFNULL(SUM(txn.amt)*a.unitprice, 0.0))
  FROM txn WHERE txn.account_id = a.account_id) AS bal
FROM account a 
LEFT OUTER JOIN currency cur ON cur.currency_id = a.currency_id 
INNER JOIN bookaccount ba ON ba.account_id = a.account_id
WHERE ba.book_id = ? AND %s`, swhere)
	rows, err := db.Query(s, bookid)
	if err != nil {
		return nil, err
	}
	aa := []*Account{}
	for rows.Next() {
		var a Account
		var c Currency
		rows.Scan(&a.Accountid, &a.Code, &a.Name, &a.AccountType, &a.Unitprice, &c.Currencyid, &c.Currency, &c.Usdrate, &a.Balance)
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
	return findAccounts(db, bookid, "1=1 ORDER BY accounttype, name")
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

func createRandomBankAccount(db *sql.DB, bookid int64) (int64, error) {
	banks := []string{"BPI", "Security", "Wells Fargo", "Bank of America", "FirstTech", "MetroBank"}
	descs := []string{"Savings", "Checking", "Time Deposit", "Money Market", "Individual", "COD", "Maxi"}
	opts := []string{"", "Cash", "Ext", "Other"}

	ibank := rand.Intn(len(banks))
	idesc := rand.Intn(len(descs))
	iopt := rand.Intn(len(opts))
	name := strings.TrimSpace(fmt.Sprintf("%s %s %s", banks[ibank], descs[idesc], opts[iopt]))

	a := Account{
		Code:        "",
		Name:        name,
		AccountType: BankAccount,
		Unitprice:   1.0,
		Currencyid:  1,
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

func createRandomStockAccount(db *sql.DB, bookid int64, ticker string, unitprice float64) (int64, error) {
	a := Account{
		Code:        "",
		Name:        ticker,
		AccountType: StockAccount,
		Unitprice:   unitprice,
		Currencyid:  1,
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
	s := "INSERT INTO txn (account_id, date, ref, desc, amt) VALUES (?, ?, ?, ?, ?)"
	result, err := sqlexec(db, s, t.Accountid, t.Date, t.Ref, t.Desc, t.Amt)
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
	s := "UPDATE txn SET date = ?, ref = ?, desc = ?, amt = ? WHERE txn_id = ?"
	_, err := sqlexec(db, s, t.Date, t.Ref, t.Desc, t.Amt, t.Txnid)
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
	s := "SELECT txn_id, account_id, date, ref, desc, amt FROM txn WHERE txn_id = ?"
	row := db.QueryRow(s, txnid)
	var t Txn
	err := row.Scan(&t.Txnid, &t.Accountid, &t.Date, &t.Ref, &t.Desc, &t.Amt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}
func findTxns(db *sql.DB, swhere string) ([]*Txn, error) {
	s := fmt.Sprintf("SELECT txn_id, account_id, date, ref, desc, amt FROM txn WHERE %s", swhere)
	rows, err := db.Query(s)
	if err != nil {
		return nil, err
	}
	tt := []*Txn{}
	for rows.Next() {
		var t Txn
		rows.Scan(&t.Txnid, &t.Accountid, &t.Date, &t.Ref, &t.Desc, &t.Amt)
		tt = append(tt, &t)
	}
	return tt, nil
}
func findTxnsOfAccount(db *sql.DB, accountid int64, swhere string) ([]*Txn, error) {
	s := fmt.Sprintf("SELECT txn_id, account_id, date, ref, desc, amt FROM txn WHERE account_id = ? AND %s", swhere)
	rows, err := db.Query(s, accountid)
	if err != nil {
		return nil, err
	}
	tt := []*Txn{}
	for rows.Next() {
		var t Txn
		rows.Scan(&t.Txnid, &t.Accountid, &t.Date, &t.Ref, &t.Desc, &t.Amt)
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
	s := "INSERT INTO txn (account_id, date, ref, desc, amt) VALUES (?, ?, ?, ?, ?)"

	for i := 0; i < ntxns; i++ {
		amt := float64(rand.Intn(5000000))/100.0 - 25000
		t := Txn{
			Accountid: accountid,
			Date:      isodate(randdate(2000, 2021)),
			Ref:       "",
			Desc:      createRandomWords(words),
			Amt:       amt,
		}
		_, err := txexec(tx, s, t.Accountid, t.Date, t.Ref, t.Desc, t.Amt)
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
	s := "INSERT INTO txn (account_id, date, ref, desc, amt) VALUES (?, ?, ?, ?, ?)"

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
		}
		_, err := txexec(tx, s, t.Accountid, t.Date, t.Ref, t.Desc, t.Amt)
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
	s := "INSERT INTO user (user_id, username, password) VALUES (?, ?, ?)"
	result, err := sqlexec(db, s, u.Userid, u.Username, u.Password)
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
