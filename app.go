package main

import (
	"fmt"
	"net/http"
	"bitbucket.com/cswank/gadgetsweb/app"
)

func main() {
	r := app.GetRouter()
	http.Handle("/", r)
	fmt.Println("listening on 0.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}

