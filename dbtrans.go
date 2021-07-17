package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"strings"
)

type Transaction struct {
	Transactionid int64   `json:"transactionid"`
	Accountid     int64   `json:"accountid"`
	Date          string  `json:"date"`
	Ref           string  `json:"ref"`
	Desc          string  `json:"desc"`
	Amt           float64 `json:"amt"`
}

func createTransaction(db *sql.DB, t *Transaction) (int64, error) {
	s := "INSERT INTO trans (account_id, date, ref, desc, amt) VALUES (?, ?, ?, ?, ?)"
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
func editTransaction(db *sql.DB, t *Transaction) error {
	s := "UPDATE trans SET date = ?, ref = ?, desc = ?, amt = ? WHERE trans_id = ?"
	_, err := sqlexec(db, s, t.Date, t.Ref, t.Desc, t.Amt, t.Transactionid)
	if err != nil {
		return err
	}
	return nil
}
func delTransaction(db *sql.DB, transactionid int64) error {
	s := "DELETE FROM trans WHERE trans_id = ?"
	_, err := sqlexec(db, s, transactionid)
	if err != nil {
		return err
	}
	return nil
}

func findTransaction(db *sql.DB, transactionid int64) (*Transaction, error) {
	s := "SELECT trans_id, account_id, date, ref, desc, amt FROM trans WHERE trans_id = ?"
	row := db.QueryRow(s, transactionid)
	var t Transaction
	err := row.Scan(&t.Transactionid, &t.Accountid, &t.Date, &t.Ref, &t.Desc, &t.Amt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}
func findTransactions(db *sql.DB, swhere string) ([]*Transaction, error) {
	s := fmt.Sprintf("SELECT trans_id, account_id, date, ref, desc, amt FROM trans WHERE %s", swhere)
	rows, err := db.Query(s)
	if err != nil {
		return nil, err
	}
	tt := []*Transaction{}
	for rows.Next() {
		var t Transaction
		rows.Scan(&t.Transactionid, &t.Accountid, &t.Date, &t.Ref, &t.Desc, &t.Amt)
		tt = append(tt, &t)
	}
	return tt, nil
}
func findAccountTransactions(db *sql.DB, accountid int64, swhere string) ([]*Transaction, error) {
	s := fmt.Sprintf("SELECT trans_id, account_id, date, ref, desc, amt FROM trans WHERE account_id = ? AND %s", swhere)
	rows, err := db.Query(s, accountid)
	if err != nil {
		return nil, err
	}
	tt := []*Transaction{}
	for rows.Next() {
		var t Transaction
		rows.Scan(&t.Transactionid, &t.Accountid, &t.Date, &t.Ref, &t.Desc, &t.Amt)
		tt = append(tt, &t)
	}
	return tt, nil
}

func findAllAccountTransactions(db *sql.DB, accountid int64) ([]*Transaction, error) {
	return findAccountTransactions(db, accountid, "1=1 ORDER BY date DESC")
}

func createRandomTransaction(db *sql.DB, accountid int64) (int64, error) {
	words := []string{"interest", "dividend", "refund", "salary", "deposit from", "receive from", "withdraw", "banking fee", "tenant", "bank", "pay"}

	nwords := 1 + rand.Intn(11)

	var sb strings.Builder
	for i := 0; i < nwords; i++ {
		sb.WriteString(words[rand.Intn(len(words))])
		if i < nwords-1 {
			sb.WriteString(" ")
		}
	}

	t := Transaction{
		Accountid: accountid,
		Date:      isodate(randdate(2000, 2021)),
		Ref:       "",
		Desc:      sb.String(),
		Amt:       rand.NormFloat64(),
	}
	return createTransaction(db, &t)
}
