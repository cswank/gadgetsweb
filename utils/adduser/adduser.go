package main

import (
	"flag"
	"fmt"
	"log"

	"code.google.com/p/gopass"
	"github.com/cswank/gadgetsweb/models"
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
	log.Println(u.Save())
}

func doDel() {
	users, err := models.GetUsers()
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(users) == 0 {
		fmt.Println("no users to delete")
		return
	}
	for i, u := range users {
		fmt.Printf("%d   %s\n", i+1, u.Username)
	}
	var j int
	var confirm string
	fmt.Print("select a user: ")
	fmt.Scanf("%d", &j)
	fmt.Printf("really delete %d (y/N)? ", j)
	fmt.Scanf("%s", &confirm)
	if confirm == "y" || confirm == "Y" {
		u := users[j-1]
		u.Delete()
	}
}
