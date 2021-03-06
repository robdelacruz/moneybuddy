package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"strings"
)

type Ledger struct {
	Ledgerid int64  `json:"ledgerid"`
	Name     string `json:"name"`
}

type AccountType int

const (
	BankAccount AccountType = iota
	StockAccount
)

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
