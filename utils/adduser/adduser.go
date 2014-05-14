package main

import (
	"flag"
	"fmt"
	"code.google.com/p/gopass"
	"bitbucket.org/cswank/gadgetsweb/models"
)

var (
	del = flag.Bool("d", false, "delete a gadget")
)

func main() {
	flag.Parse()
	if *del {
		doDel()
	} else {
		add()
	}
}

func add() {
	u := models.User{}
	fmt.Print("username: ")
	fmt.Scanf("%s", &u.Username)
	fmt.Print("can write? (y/N): ")
	var perm string
	fmt.Scanf("%s", &perm)
	if perm == "y" || perm == "Y" {
		u.Permission = "write"
	}
	u.Password, _ = gopass.GetPass("password: ")
	fmt.Println(u)
	u.Save()
}

func doDel() {
	users := models.GetUsers()
	if len(users) == 0 {
		fmt.Println("no users to delete")
		return
	}
	for i, u := range users {
		fmt.Printf("%d   %s\n", i + 1, u.Username)
	}
	var j int
	var confirm string
	fmt.Print("select a user: ")
	fmt.Scanf("%d", &j)
	fmt.Printf("really delete %d (y/N)? ", j)
	fmt.Scanf("%s", &confirm)
	if confirm == "y" || confirm == "Y" {
		u := users[j - 1]
		u.Delete()
	}
}

