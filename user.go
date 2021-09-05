package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func genHash(sinput string) string {
	bsHash, err := bcrypt.GenerateFromPassword([]byte(sinput), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(bsHash)
}
func validateHash(shash, sinput string) bool {
	if shash == "" && sinput == "" {
		return true
	}
	err := bcrypt.CompareHashAndPassword([]byte(shash), []byte(sinput))
	if err != nil {
		return false
	}
	return true
}
func genSig(u *User) string {
	sig := genHash(fmt.Sprintf("%s_%s", u.Username, u.Password))
	return sig
}
func validateSig(sig string, u *User) bool {
	return validateHash(sig, fmt.Sprintf("%s_%s", u.Username, u.Password))
}

func signup(db *sql.DB, username, pwd string) (*User, error) {
	fexists, err := isUsernameExists(db, username)
	if err != nil {
		return nil, err
	}
	if fexists {
		return nil, fmt.Errorf("username '%s' already exists", username)
	}

	var u User
	u.Username = username
	u.Password = genHash(pwd)
	newid, err := createUser(db, &u)
	if err != nil {
		return nil, err
	}
	u.Userid = newid

	err = initUserData(db, &u)
	if err != nil {
		delUser(db, u.Userid)
		return nil, err
	}

	return &u, nil
}

// Create initial data set for new user.
func initUserData(db *sql.DB, u *User) error {
	c := Currency{
		Currency: "USD",
		Usdrate:  1.0,
		Userid:   u.Userid,
	}
	_, err := createCurrency(db, &c)
	if err != nil {
		return err
	}
	c = Currency{
		Currency: "PHP",
		Usdrate:  48.0,
		Userid:   u.Userid,
	}
	_, err = createCurrency(db, &c)
	if err != nil {
		return err
	}

	b := Book{
		Name:   fmt.Sprintf("%s's Accounts", strings.Title(u.Username)),
		Userid: u.Userid,
	}
	_, err = createBook(db, &b)
	if err != nil {
		return err
	}
	b = Book{
		Name:   "Other Accounts",
		Userid: u.Userid,
	}
	_, err = createBook(db, &b)
	if err != nil {
		return err
	}

	return nil
}

var ErrLoginIncorrect = errors.New("Incorrect username or password")

func login(u *User, pwd string) (string, error) {
	if !validateHash(u.Password, pwd) {
		return "", ErrLoginIncorrect
	}
	// Return user signature, this will be used to authenticate user per request.
	sig := genSig(u)
	return sig, nil
}
func loginUsername(db *sql.DB, username, pwd string) (string, error) {
	u, err := findUserByUsername(db, username)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", ErrLoginIncorrect
	}
	sig, err := login(u, pwd)
	if err != nil {
		return "", err
	}
	return sig, nil
}

func validateUserSig(db *sql.DB, userid int64, sig string) *User {
	if userid == 0 {
		return nil
	}
	u, err := findUser(db, userid)
	if err != nil {
		return nil
	}
	if u == nil {
		return nil
	}
	if !validateSig(sig, u) {
		return nil
	}
	return u
}

func validateApiUser(db *sql.DB, r *http.Request) *User {
	// Get user making the request. There are two ways to specify user:
	// - Through querystring siguserid and sig
	// - Through http cookies 'userid' and 'sig'
	quserid := idtoi(r.FormValue("siguserid"))
	qsig := r.FormValue("sig")
	if quserid == 0 {
		// cookie format: user=<userid>|<username>|<sig>
		vvv := strings.Split(readCookie(r, "user"), "|")
		if len(vvv) < 3 {
			return nil
		}
		quserid = idtoi(vvv[0])
		qsig = vvv[2]
	}

	fmt.Printf("validateApiUser: quserid=%d, qsig=%s\n", quserid, qsig)

	return validateUserSig(db, quserid, qsig)
}
