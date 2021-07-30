package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"strings"
)

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
	Currencyid  int64       `json:"currencyid"`
	Currency    string      `json:"currency"`
	Balance     float64     `json:"balance"`
	Txns        []*Txn      `json:"txns"`
}

func createAccount(db *sql.DB, a *Account) (int64, error) {
	s := "INSERT INTO account (code, name, accounttype, currency_id) VALUES (?, ?, ?, ?)"
	result, err := sqlexec(db, s, a.Code, a.Name, a.AccountType, a.Currencyid)
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
	s := "UPDATE account SET code = ?, name = ?, accounttype = ?, currency_id = ? WHERE account_id = ?"
	_, err := sqlexec(db, s, a.Code, a.Name, a.AccountType, a.Currencyid, a.Accountid)
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
	return nil
}

func findAccount(db *sql.DB, accountid int64) (*Account, error) {
	s := `SELECT account_id, code, name, accounttype, a.currency_id, IFNULL(cur.currency, ''), 
(SELECT SUM(amt) FROM txn WHERE txn.account_id = a.account_id) AS bal
FROM account a 
LEFT OUTER JOIN currency cur ON cur.currency_id = a.currency_id 
WHERE account_id = ?`
	row := db.QueryRow(s, accountid)
	var a Account
	err := row.Scan(&a.Accountid, &a.Code, &a.Name, &a.AccountType, &a.Currencyid, &a.Currency, &a.Balance)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	tt, err := findAllTxnsOfAccount(db, accountid)
	if err != nil {
		return nil, err
	}
	a.Txns = tt
	return &a, nil
}
func findAccounts(db *sql.DB, swhere string) ([]*Account, error) {
	s := fmt.Sprintf(`
SELECT account_id, code, name, accounttype, a.currency_id, IFNULL(cur.currency, ''), 
(SELECT SUM(amt) FROM txn WHERE txn.account_id = a.account_id) AS bal
FROM account a 
LEFT OUTER JOIN currency cur ON cur.currency_id = a.currency_id 
WHERE %s`, swhere)
	rows, err := db.Query(s)
	if err != nil {
		return nil, err
	}
	aa := []*Account{}
	for rows.Next() {
		var a Account
		rows.Scan(&a.Accountid, &a.Code, &a.Name, &a.AccountType, &a.Currencyid, &a.Currency, &a.Balance)
		aa = append(aa, &a)
	}
	return aa, nil
}
func findAccountsTxns(db *sql.DB, swhere string) ([]*Account, error) {
	aa, err := findAccounts(db, swhere)
	if err != nil {
		return nil, err
	}
	for _, a := range aa {
		tt, err := findAllTxnsOfAccount(db, a.Accountid)
		if err != nil {
			return nil, err
		}
		a.Txns = tt
	}
	return aa, nil
}
func findAllAccounts(db *sql.DB) ([]*Account, error) {
	return findAccounts(db, "1=1 ORDER BY name")
}
func findAllAccountsTxns(db *sql.DB) ([]*Account, error) {
	return findAccountsTxns(db, "1=1 ORDER BY name")
}

func balAccount(db *sql.DB, accountid int64) float64 {
	s := "SELECT IFNULL(SUM(amt), 0.0) FROM txn WHERE account_id = ?"
	row := db.QueryRow(s, accountid)
	var bal float64
	err := row.Scan(&bal)
	if err != nil {
		return 0.0
	}
	return bal
}

func createRandomAccount(db *sql.DB) (int64, error) {
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
		Currencyid:  1,
	}
	accountid, err := createAccount(db, &a)
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
