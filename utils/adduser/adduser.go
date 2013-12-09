package main

import (
	"fmt"
	"code.google.com/p/gopass"
	"bitbucket.com/cswank/gadgetsweb/models"
	
)

func main() {
	u := models.User{}
	fmt.Print("username: ")
	fmt.Scanf("%s", &u.Username)
	u.Password, _ = gopass.GetPass("password: ")
	fmt.Println(u)
	u.Save()
}
