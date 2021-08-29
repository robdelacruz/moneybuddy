package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Rootdata struct {
	Currencies []*Currency `json:"currencies"`
	Books      []*Book     `json:"books"`
}

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
