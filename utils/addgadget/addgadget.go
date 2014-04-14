package main

import (
	"fmt"
	"flag"
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

func doDel() {
	gadgets := models.GetGadgets()
	if len(gadgets) == 0 {
		fmt.Println("no gadgets to delete")
		return
	}
	for i, g := range gadgets {
		fmt.Printf("%d   %s\n", i + 1, g.Name)
	}
	var j int
	var confirm string
	fmt.Print("select a gadget: ")
	fmt.Scanf("%d", &j)
	fmt.Printf("really delete %d (y/N)? ", j)
	fmt.Scanf("%s", &confirm)
	if confirm == "y" || confirm == "Y" {
		g := gadgets[j - 1]
		g.Delete()
	}
}

func add() {
	g := models.Gadget{}
	fmt.Print("name: ")
	fmt.Scanf("%s", &g.Name)
	fmt.Print("host: ")
	fmt.Scanf("%s", &g.Host)
	fmt.Print(fmt.Sprintf("really save gadget (name: %s, host: %s)? (Y/n) ", g.Name, g.Host))
	var save string
	fmt.Scanf("%s", &save)
	if save == "y" || save == "Y" || save == "" {
		fmt.Println(g.Save())
	} else {
		fmt.Println("not saving")
	}
}
