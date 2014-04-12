package models

import (
	"errors"
	"code.google.com/p/go.crypto/bcrypt"
)


//Is authorized if the username is in the db
func (u *User)IsAuthorized() bool {
	db := getDB()
	user := db.Users[u.Username]
	return len(user.Username) != 0 && user.Username == u.Username
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


