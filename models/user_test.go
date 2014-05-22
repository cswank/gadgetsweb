package models

import (
	"os"
	"testing"
	"io/ioutil"
	"path"
)


func TestSaveUser(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	os.Setenv("GADGETSDB", path.Join(tmp, "db"))
	db, err := GetDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	u := User{
		Username: "craig",
		Password: "xyatooks",
	}
	err = u.Save()
	if err != nil {
		t.Fatal(err)
	}
	var pw string
	err = db.QueryRow(getPasswordQuery, "craig").Scan(&pw)
	if err != nil {
		t.Fatal(err)
	}
	if pw != string(u.HashedPassword) {
		t.Error(pw)
	}
	os.RemoveAll(tmp)
}

func TestDeleteUser(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	os.Setenv("GADGETSDB", path.Join(tmp, "db"))
	db, err := GetDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	u := User{
		Username: "craig",
		Password: "xyatooks",
	}
	u.Save()
	users, err := GetUsers()
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 {
		t.Fatal(users)
	}

	u = users[0]
	err = u.Delete()
	if err != nil {
		t.Fatal(err)
	}
	users, err = GetUsers()
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 0 {
		t.Error(users)
	}
	os.RemoveAll(tmp)
}


func TestGetUsers(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	os.Setenv("GADGETSDB", path.Join(tmp, "db"))
	db, err := GetDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	u := User{
		Username: "craig",
		Password: "xyatooks",
	}
	u.Save()
	u = User{
		Username: "laura",
		Password: "xyatookss",
	}
	u.Save()

	users, err := GetUsers()
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatal(users)
	}
	if users[0].Username != "craig" {
		t.Fatal(users[0])
	}
	os.RemoveAll(tmp)
}

func TestIsAuthorized(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	os.Setenv("GADGETSDB", path.Join(tmp, "db"))
	db, _ := GetDB()
	defer db.Close()
	u := User{
		Username: "craig",
		Password: "xyatooks",
	}

	if u.IsAuthorized("read") {
		t.Error("shouldn't be authorized")
	}
	
	u.Save()

	if !u.IsAuthorized("read") {
		t.Error("should be authorized")
	}

	if u.IsAuthorized("write") {
		t.Error("should not be authorized to write")
	}
	os.RemoveAll(tmp)
}

func TestIsAuthorizedWithWrite(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	os.Setenv("GADGETSDB", path.Join(tmp, "db"))
	db, _ := GetDB()
	defer db.Close()
	u := User{
		Username: "craig",
		Password: "xyatooks",
		Permission: "write",
	}

	if u.IsAuthorized("write") {
		t.Error("shouldn't be authorized")
	}
	u.Save()
	if !u.IsAuthorized("write") {
		t.Error("should be authorized")
	}
	os.RemoveAll(tmp)
}

func TestCheckPassword(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	os.Setenv("GADGETSDB", path.Join(tmp, "db"))
	db, _ := GetDB()
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
	os.RemoveAll(tmp)
}


