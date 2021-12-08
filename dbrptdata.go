package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Rptdata struct {
	Userid   int64      `json:"userid"`
	Currency *Currency  `json:"currency"`
	BookRpts []*BookRpt `json:"bookrpts"`
}

type BookRpt struct {
	Bookid     int64       `json:"bookid"`
	Bookname   string      `json:"bookname"`
	SummaryRpt *SummaryRpt `json:"summaryrpt"`
}

type SummaryRpt struct {
	Heading string            `json:"heading"`
	Items   []*SummaryRptItem `json:"rptitems"`
}

type SummaryRptItem struct {
	Code    string       `json:"code"`
	Caption string       `json:"caption"`
	Cols    []RptItemCol `json:"cols"`
}

type RptItemCol interface{}

type CurrencyAmt struct {
	CurrencyName string  `json:"currencyname"`
	Amt          float64 `json:"amt"`
}

func findRptdata(db *sql.DB, userid, currencyid int64) (*Rptdata, error) {
	c, err := findCurrency(db, currencyid)
	if err != nil {
		return nil, err
	}
	if c == nil {
		c = &Currency{1, "USD", 1.0, 0}
	}

	bb, err := findUserBooks(db, userid)
	if err != nil {
		return nil, err
	}

	rr := []*BookRpt{}
	for _, b := range bb {
		if b.BookType != UserBook {
			continue
		}

		var r BookRpt
		r.Bookid = b.Bookid
		r.Bookname = b.Name
		summaryrpt, err := findSummaryRpt(db, b, c)
		if err != nil {
			return nil, err
		}
		r.SummaryRpt = summaryrpt
		rr = append(rr, &r)
	}

	var rptdata Rptdata
	rptdata.Userid = userid
	rptdata.Currency = c
	rptdata.BookRpts = rr
	return &rptdata, nil
}

func makeSummaryRptItem(code, caption string, cols ...RptItemCol) *SummaryRptItem {
	if cols == nil {
		cols = []RptItemCol{}
	}
	return &SummaryRptItem{
		Code:    code,
		Caption: caption,
		Cols:    cols,
	}
}

var EmptyCurrencyAmt = CurrencyAmt{"", 0.0}
var EmptySummaryRptItem = makeSummaryRptItem("", "")

func findSummaryRpt(db *sql.DB, b *Book, c *Currency) (*SummaryRpt, error) {
	var bankBal, stockBal, totalBal float64
	for _, a := range b.BankAccounts {
		bankBal += convToCurrency(a.Balance, a.Currency, c)
	}
	for _, a := range b.StockAccounts {
		stockBal += convToCurrency(a.Balance, a.Currency, c)
	}
	totalBal = bankBal + stockBal

	var items []*SummaryRptItem
	items = append(items, makeSummaryRptItem("", "All Accounts", CurrencyAmt{c.Name, totalBal}))
	items = append(items, makeSummaryRptItem("", "Bank Accounts", CurrencyAmt{c.Name, bankBal}))
	items = append(items, makeSummaryRptItem("", "Stocks", CurrencyAmt{c.Name, stockBal}))
	items = append(items, EmptySummaryRptItem)

	items = append(items, makeSummaryRptItem("", "# Bank Accounts", "Converted", "Original"))
	for _, a := range b.BankAccounts {
		bankbal := convToCurrency(a.Balance, a.Currency, c)
		items = append(items, makeSummaryRptItem("", a.Name, CurrencyAmt{c.Name, bankbal}, CurrencyAmt{a.Currency.Name, a.Balance}))
	}
	items = append(items, EmptySummaryRptItem)

	p := message.NewPrinter(language.English)

	items = append(items, makeSummaryRptItem("", "# Stocks", "Converted", "Original"))
	for _, a := range b.StockAccounts {
		nshares, err := accountSumAmt(db, a.Accountid)
		if err != nil {
			return nil, err
		}
		stockdesc := p.Sprintf("%s (%.2f shares)", a.Name, nshares)
		stockbal := convToCurrency(a.Balance, a.Currency, c)
		items = append(items, makeSummaryRptItem("", stockdesc, CurrencyAmt{c.Name, stockbal}, CurrencyAmt{a.Currency.Name, a.Balance}))
	}

	var summaryrpt SummaryRpt
	summaryrpt.Heading = fmt.Sprintf("Summary Report for '%s'", b.Name)
	summaryrpt.Items = items
	return &summaryrpt, nil
}
