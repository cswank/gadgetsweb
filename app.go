package main

import (
	"fmt"
	"bitbucket.com/cswank/gadgetsweb/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", GetHome).Methods("GET")
	r.HandleFunc("/history/locations/summary", GetSummary).Methods("GET")
	r.HandleFunc("/history/locations/{location}/directions/{direction}/devices/{device}", GetTimeseries).Methods("GET")
	r.HandleFunc("/socket", GetSocket)

	http.Handle("/", r)
	fmt.Println("listening on 0.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	
}

func GetTimeseries(w http.ResponseWriter, r *http.Request) {
	controllers.GetTimeseries(w, r)
}

func GetSummary(w http.ResponseWriter, r *http.Request) {
	controllers.GetSummary(w, r)
}

func GetSocket(w http.ResponseWriter, r *http.Request) {
	controllers.HandleSocket(w, r)
}
