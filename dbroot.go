package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Rootdata struct {
	Accounts   []*Account  `json:"accounts"`
	Currencies []*Currency `json:"currencies"`
}

func findRootdata(db *sql.DB) (*Rootdata, error) {
	aa, err := findAllAccountsTxns(db)
	if err != nil {
		return nil, err
	}
	cc, err := findAllCurrencies(db)
	if err != nil {
		return nil, err
	}

	var rd Rootdata
	rd.Accounts = aa
	rd.Currencies = cc
	return &rd, nil
}
