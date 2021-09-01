package main

import (
	"database/sql"
	"errors"
	"fmt"
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
	return &u, nil
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
