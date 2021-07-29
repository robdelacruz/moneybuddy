package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Model struct {
	Accounts   []*Account  `json:"accounts"`
	Currencies []*Currency `json:"currencies"`
}

func findModel(db *sql.DB) (*Model, error) {
	aa, err := findAllAccountsTxns(db)
	if err != nil {
		return nil, err
	}
	cc, err := findAllCurrencies(db)
	if err != nil {
		return nil, err
	}

	var m Model
	m.Accounts = aa
	m.Currencies = cc
	return &m, nil
}
