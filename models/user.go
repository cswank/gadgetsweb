package models

import (
	"errors"
	"code.google.com/p/go.crypto/bcrypt"
)

var (
	getPasswordQuery = "SELECT password FROM users WHERE username = ?"
	getUserQuery = "SELECT username, permission FROM users WHERE username = ?"
	saveUserQuery = "INSERT INTO users (username, password, permission) VALUES (?, ?, ?)"
)

type User struct {
	Id uint64 `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
	HashedPassword []byte `json:"-"`
	Permission string `json:"permission"`
}

//Is authorized if the username is in the db
func (u *User)IsAuthorized(permission string) bool {
	db, err := getDB()
	defer db.Close()
	if err != nil {
		return false
	}
	var uname string
	err = db.QueryRow(getUserQuery, u.Username).Scan(&uname, &u.Permission)
	a := len(uname) != 0 && uname == u.Username
	if permission == "write" {
		a = a && u.Permission == "write"
	}
	return err == nil && a
}

func (u *User)Save() error {
	if len(u.Password) < 8 {
		return errors.New("password is too short")
	}
	u.hashPassword()
	db, err := getDB()
	defer db.Close()
	if err != nil {
		return err
	}
	_, err = db.Query(saveUserQuery, u.Username, u.HashedPassword, u.Permission)
	return nil
}

func (u *User)CheckPassword() (bool, error) {
	err := u.getHashedPassword()
	if err != nil {
		return false, err
	}
	isGood:= bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(u.Password))
	return isGood == nil, err
}

func (u *User)hashPassword() {
	u.HashedPassword, _ = bcrypt.GenerateFromPassword(
		[]byte(u.Password),
		bcrypt.DefaultCost,
	)
}

func (u *User)getHashedPassword() error {
	db, err := getDB()
	defer db.Close()
	if err != nil {
		return err
	}
	var hashedPassword string
	err = db.QueryRow(getPasswordQuery, u.Username).Scan(&hashedPassword)
	if err == nil {
		u.HashedPassword = []byte(hashedPassword)
	}
	return err
}


