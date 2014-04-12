package models

import (
	"testing"
	"io/ioutil"
	"path"
)


func TestSaveUser(t *testing.T) {
	u := User{
		Username: "craig",
		Password: "xyatooks",
	}
	err := u.Save()
	if err != nil {
		t.Error(err)
	}
	db := getDB()
	u = db.Users["craig"]
	if u.Username != "craig" {
		t.Error(u)
	}
}

func TestIsAuthorized(t *testing.T) {
	tmp, _ := ioutil.TempDir("/tmp", "")
	DBPath = path.Join(tmp, "gadgets.db")
	u := User{
		Username: "me",
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
	tmp, _ := ioutil.TempDir("/tmp", "")
	DBPath = path.Join(tmp, "gadgets.db")
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
}


