package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xuri/excelize"
)

func main() {
	_, parms := parseArgs(os.Args[1:])
	if len(parms) < 2 {
		fmt.Printf("importer <dbfile> <xlsfile> <bookname>\n")
		os.Exit(0)
	}

	// Exit if db file doesn't exist.
	dbfile := parms[0]
	if !fileExists(dbfile) {
		fmt.Printf("DB File '%s' doesn't exist.\n", dbfile)
		os.Exit(1)
	}
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		panic(err)
	}

	qexcelfile := parms[1]
	var qbookname string
	if len(parms) < 3 {
		qbookname = strings.Split(qexcelfile, ".")[0]
	} else {
		qbookname = parms[2]
	}

	f, err := excelize.OpenFile(qexcelfile)
	if err != nil {
		panic(err)
	}

	username := "rob"
	u, err := findUserByUsername(db, username)
	if err != nil {
		panic(err)
	}
	if u == nil {
		fmt.Printf("Username '%s' not found.\n", username)
		os.Exit(1)
	}

	// Delete any existing books with the same name.
	existingbookids := []int64{}
	s := "SELECT book_id FROM book WHERE user_id = ? AND name = ?"
	rows, err := db.Query(s, u.Userid, qbookname)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var bookid int64
		rows.Scan(&bookid)
		existingbookids = append(existingbookids, bookid)
	}
	for _, bookid := range existingbookids {
		err := delBook(db, bookid)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Deleted existing book '%s' (id: %d)\n", qbookname, bookid)
	}

	// Create new book for the new accounts.
	b := Book{
		Name:     qbookname,
		BookType: UserBook,
		Userid:   u.Userid,
		Active:   1,
	}
	newid, err := createBook(db, &b)
	if err != nil {
		panic(err)
	}
	b.Bookid = newid
	fmt.Printf("Created book: %s\n", b.Name)

	sheets := []string{
		"BPIChecking",
		"BpiUSD",
		"SecurityUSD",
		"SecurityUSD2",
		"SecuritySavingsUSD",
		"BPITimeDeposit",
		"SecuritySavings",
		"MoneyMkt",
	}

	for _, sheetname := range sheets {
		// Create new account per sheet.
		a := Account{
			Code:        strings.ToLower(sheetname),
			Name:        sheetname,
			AccountType: BankAccount,
			Unitprice:   1.0,
			Currencyid:  2,
		}
		if sheetname == "MoneyMkt" || sheetname == "BpiUSD" || sheetname == "SecurityUSD" || sheetname == "SecurityUSD2" || sheetname == "SecuritySavingsUSD" {
			a.Currencyid = 1
		}
		newid, err := createAccount(db, &a, b.Bookid)
		if err != nil {
			panic(err)
		}
		a.Accountid = newid
		fmt.Printf("Created account: %s\n", a.Name)

		importSheet(db, f, a.Accountid, sheetname)
	}
}

func parseArgs(args []string) (map[string]string, []string) {
	switches := map[string]string{}
	parms := []string{}

	standaloneSwitches := []string{}
	definitionSwitches := []string{"i"}
	fNoMoreSwitches := false
	curKey := ""

	for _, arg := range args {
		if fNoMoreSwitches {
			// any arg after "--" is a standalone parameter
			parms = append(parms, arg)
		} else if arg == "--" {
			// "--" means no more switches to come
			fNoMoreSwitches = true
		} else if strings.HasPrefix(arg, "--") {
			switches[arg[2:]] = "y"
			curKey = ""
		} else if strings.HasPrefix(arg, "-") {
			if listContains(definitionSwitches, arg[1:]) {
				// -a "val"
				curKey = arg[1:]
				continue
			}
			for _, ch := range arg[1:] {
				// -a, -b, -ab
				sch := string(ch)
				if listContains(standaloneSwitches, sch) {
					switches[sch] = "y"
				}
			}
		} else if curKey != "" {
			switches[curKey] = arg
			curKey = ""
		} else {
			// standalone parameter
			parms = append(parms, arg)
		}
	}

	return switches, parms
}

func listContains(ss []string, v string) bool {
	for _, s := range ss {
		if v == s {
			return true
		}
	}
	return false
}
func fileExists(file string) bool {
	_, err := os.Stat(file)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func importSheet(db *sql.DB, f *excelize.File, accountid int64, name string) {
	rows, err := f.GetRows(name)
	if err != nil {
		panic(err)
	}

	var headingrow int
	for i, row := range rows {
		if len(row) == 0 {
			continue
		}

		if row[0] == "Date" {
			headingrow = i
			break
		}
	}

	if headingrow == 0 {
		fmt.Printf("Sheet '%s' not imported: heading row not found.\n", name)
		return
	}

	for _, row := range rows[headingrow+1:] {
		ncols := len(row)

		// Reached the end of transaction rows
		if ncols >= 3 && strings.TrimSpace(row[2]) == "Total" {
			break
		}

		// Skip if no deposit/withdraw amount
		if ncols < 4 {
			continue
		}

		var sdate, checkno, desc, sdeposit, swithdraw string
		var amt float64

		sdate = strings.TrimSpace(row[0])
		checkno = strings.TrimSpace(row[1])
		desc = strings.TrimSpace(row[2])
		sdeposit = strings.TrimSpace(row[3])
		if ncols >= 5 {
			swithdraw = strings.TrimSpace(row[4])
		}

		if sdeposit != "" {
			amt += atof(sdeposit)
		}
		if swithdraw != "" {
			amt -= atof(swithdraw)
		}

		t := Txn{
			Accountid: accountid,
			Date:      sdate,
			Ref:       checkno,
			Desc:      desc,
			Amt:       amt,
		}
		newid, err := createTxn(db, &t)
		if err != nil {
			panic(err)
		}
		t.Txnid = newid

		fmt.Printf("Added txn: [%s], [%s], \"%s\", [%.2f]\n", t.Date, t.Ref, t.Desc, t.Amt)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}
}
