package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"strings"
)

type Txn struct {
	Txnid     int64   `json:"txnid"`
	Accountid int64   `json:"accountid"`
	Date      string  `json:"date"`
	Ref       string  `json:"ref"`
	Desc      string  `json:"desc"`
	Amt       float64 `json:"amt"`
}

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
		t := Txn{
			Accountid: accountid,
			Date:      isodate(randdate(2000, 2021)),
			Ref:       "",
			Desc:      createRandomWords(words),
			Amt:       rand.NormFloat64()*50000 + 30000,
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
