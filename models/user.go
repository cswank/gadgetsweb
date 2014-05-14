package models

import (
	"fmt"
	"errors"
	"code.google.com/p/go.crypto/bcrypt"
)

func GetUsers() []User {
	db := getDB()
	users := make([]User, len(db.Users))
	var i int
	for _, u := range db.Users {
		users[i] = u
		i++
	}
	return users
}

//Is authorized if the username is in the db and has the
//requested permissions
func (u *User)IsAuthorized(permission string) bool {
	db := getDB()
	user := db.Users[u.Username]
	a := len(user.Username) != 0 && user.Username == u.Username
	if permission == "write" {
		a = a && user.Permission == "write"
	}
	return a
}

func (u *User)Save() error {
	if len(u.Password) < 8 {
		return errors.New("password is too short")
	}
	u.hashPassword()
	db := getDB()
	db.Users[u.Username] = *u
	return db.Save()
}

func (u *User)Delete() error {
	db := getDB()
	delete (db.Users, u.Username)
	return db.Save()
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
	db := getDB()
	user, ok := db.Users[u.Username]
	if !ok {
		return errors.New("user not found in database")
	}
	u.HashedPassword = user.HashedPassword
	return nil
}


