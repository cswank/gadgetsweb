package main

import (
	"fmt"
	"bitbucket.org/cswank/gadgetsweb/models"
)

func main() {
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
