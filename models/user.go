package models

import (
	"errors"

	"code.google.com/p/go.crypto/bcrypt"
)

var (
	getUsersQuery    = "SELECT username FROM users"
	getPasswordQuery = "SELECT password FROM users WHERE username = ?"
	deleteUserQuery  = "DELETE FROM users WHERE username = ?"
	getUserQuery     = "SELECT username, permission FROM users WHERE username = ?"
	saveUserQuery    = "INSERT INTO users (username, password, permission) VALUES (?, ?, ?)"
)

type User struct {
	Id             uint64 `json:"-"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	HashedPassword []byte `json:"-"`
	Permission     string `json:"permission"`
}

func GetUsers() ([]User, error) {
	users := []User{}
	rows, err := DB.Query(getUsersQuery)
	if err != nil {
		return users, err
	}
	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.Username)
		if err != nil {
			return []User{}, err
		}
		users = append(users, u)
	}
	return users, err
}

//Is authorized if the username is in the db
func (u *User) IsAuthorized(permission string) bool {
	var uname string
	err := DB.QueryRow(getUserQuery, u.Username).Scan(&uname, &u.Permission)
	a := len(uname) != 0 && uname == u.Username
	if permission == "write" {
		a = a && u.Permission == "write"
	}
	return err == nil && a
}

func (u *User) Save() error {
	if len(u.Password) < 8 {
		return errors.New("password is too short")
	}
	u.hashPassword()
	_, err := DB.Exec(saveUserQuery, u.Username, u.HashedPassword, u.Permission)
	return err
}

func (u *User) Delete() error {
	_, err := DB.Exec(deleteUserQuery, u.Username)
	return err
}

func (u *User) CheckPassword() (bool, error) {
	err := u.getHashedPassword()
	if err != nil {
		return false, err
	}
	return bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(u.Password)) == nil, err
}

func (u *User) hashPassword() {
	u.HashedPassword, _ = bcrypt.GenerateFromPassword(
		[]byte(u.Password),
		bcrypt.DefaultCost,
	)
}

func (u *User) getHashedPassword() error {
	var hashedPassword string
	err := DB.QueryRow(getPasswordQuery, u.Username).Scan(&hashedPassword)
	if err == nil {
		u.HashedPassword = []byte(hashedPassword)
	}
	return err
}
