package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	Bookid        int64      `json:"bookid"`
	Name          string     `json:"name"`
	BankAccounts  []*Account `json:"bankaccounts"`
	StockAccounts []*Account `json:"stockaccounts"`
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
