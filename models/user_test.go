package models

import (
	"testing"
	
)


func TestSaveUser(t *testing.T) {
	db, _ := getDB()
	defer db.Close()
	db.Query("CREATE TABLE users(username text PRIMARY KEY, password text)")
	db.Query("DELETE FROM users")
	u := User{
		Username: "craig",
		Password: "xyatooks",
	}
	err := u.Save()
	if err != nil {
		t.Error(err)
	}
	var pw string
	err = db.QueryRow(getPasswordQuery, "craig").Scan(&pw)
	if err != nil {
		t.Error(err)
	}
	if pw != string(u.HashedPassword) {
		t.Error(pw)
	}
	db.Query("DELETE FROM users")
}

func TestIsAuthorized(t *testing.T) {
	db, _ := getDB()
	defer db.Close()
	db.Query("CREATE TABLE users(username text PRIMARY KEY, password text)")
	db.Query("DELETE FROM users")
	u := User{
		Username: "craig",
		Password: "xyatooks",
	}

	if u.IsAuthorized() {
		t.Error("shouldn't be authorized")
	}
	
	u.Save()

	if !u.IsAuthorized() {
		t.Error("should be authorized")
	}
}

func TestCheckPassword(t *testing.T) {
	db, _ := getDB()
	defer db.Close()
	db.Query("CREATE TABLE users(username text PRIMARY KEY, password text)")
	db.Query("DELETE FROM users")
	u := User{
		Username: "craig",
		Password: "xyatooks",
	}
	u.Save()

	u2 := User{
		Username: "craig",
		Password: "xyatooks",
	}
	isGood, _ := u2.CheckPassword()
	if !isGood {
		t.Error("password didn't match")
	}

	u3 := User{
		Username: "craig",
		Password: "xyatooks!",
	}
	isGood, _ = u3.CheckPassword()
	if isGood {
		t.Error("password shouldn't have matched")
	}
	//db.Query("DELETE FROM users")
}


