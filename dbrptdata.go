package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Rptdata struct {
	Currency *Currency  `json:"currency"`
	BookRpts []*BookRpt `json:"bookrpts"`
}

type BookRpt struct {
	Bookid   int64      `json:"bookid"`
	Bookname string     `json:"bookname"`
	RptItems []*RptItem `json:"rptitems"`
}

type RptItem struct {
	Caption string  `json:"caption"`
	Val     float64 `json:"val"`
}

func findRptdata(db *sql.DB, currencyid int64) (*Rptdata, error) {
	c, err := findCurrency(db, currencyid)
	if err != nil {
		return nil, err
	}
	if c == nil {
		c = &Currency{1, "USD", 1.0}
	}

	bb, err := findAllBooks(db)
	if err != nil {
		return nil, err
	}

	var rr []*BookRpt
	for _, b := range bb {
		r, err := findBookRptdata(db, b, c)
		if err != nil {
			return nil, err
		}
		rr = append(rr, r)
	}

	var rptdata Rptdata
	rptdata.Currency = c
	rptdata.BookRpts = rr
	return &rptdata, nil
}

var EmptyRptItem = RptItem{"", 0.0}

func findBookRptdata(db *sql.DB, b *Book, c *Currency) (*BookRpt, error) {
	var bankBal, stockBal, totalBal float64
	for _, a := range b.BankAccounts {
		bankBal += convToCurrency(a.Balance, a.Currency, c)
	}
	for _, a := range b.StockAccounts {
		stockBal += convToCurrency(a.Balance, a.Currency, c)
	}
	totalBal = bankBal + stockBal

	var items []*RptItem
	items = append(items, &RptItem{"All Accounts", totalBal})
	items = append(items, &RptItem{"Bank Accounts", bankBal})
	items = append(items, &RptItem{"Stocks", stockBal})
	items = append(items, &EmptyRptItem)

	items = append(items, &RptItem{"# Bank Accounts", 0})
	for _, a := range b.BankAccounts {
		bankbal := convToCurrency(a.Balance, a.Currency, c)
		items = append(items, &RptItem{a.Name, bankbal})
	}
	items = append(items, &EmptyRptItem)

	items = append(items, &RptItem{"# Stocks", 0})
	for _, a := range b.StockAccounts {
		nshares, err := accountSumAmt(db, a.Accountid)
		if err != nil {
			return nil, err
		}
		stockdesc := fmt.Sprintf("%s: %.2f shares at %.2f/share", a.Name, nshares, a.Unitprice)
		stockbal := convToCurrency(a.Balance, a.Currency, c)
		items = append(items, &RptItem{stockdesc, stockbal})
	}

	var r BookRpt
	r.Bookid = b.Bookid
	r.Bookname = b.Name
	r.RptItems = items
	return &r, nil
}
