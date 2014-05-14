package main

import (
	"fmt"
	"os"
	"net/http"
	"bitbucket.org/cswank/gadgetsweb/controllers"
	"bitbucket.org/cswank/gadgetsweb/auth"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var (
	hashKey        = []byte(os.Getenv("HASH_KEY"))
	blockKey       = []byte(os.Getenv("BLOCK_KEY"))
	SecureCookie   = securecookie.New(hashKey, blockKey)
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/login", auth.Login).Methods("POST")
	r.HandleFunc("/api/logout", auth.Logout).Methods("POST")
	r.HandleFunc("/api/socket/in", GetInSocket).Methods("GET")
	r.HandleFunc("/api/socket/out", GetOutSocket).Methods("GET")
	r.HandleFunc("/api/gadgets", GetGadgets).Methods("GET")
	r.HandleFunc("/api/gadgets", AddGadgets).Methods("POST")
	r.HandleFunc("/api/gadgets/types", GetGadgetTypes).Methods("GET")
	r.HandleFunc("/api/gadgets/{name}/methods", GetMethods).Methods("GET")
	r.HandleFunc("/api/gadgets/{name}/methods", AddMethod).Methods("POST")
	r.HandleFunc("/api/gadgets/{name}/methods/{methodId}", UpdateMethod).Methods("PUT")
	r.HandleFunc("/api/gadgets/{name}/methods/{methodId}", DeleteMethod).Methods("DELETE")
	r.HandleFunc("/api/history/gadgets/{gadget}/devices", GetDevices).Methods("GET")
	r.HandleFunc("/api/history/gadgets/{gadget}/locations/{location}/devices/{device}", GetTimeseries).Methods("GET")
	
	http.Handle("/", r)
	fmt.Println("listening on 0.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}

func GetGadgets(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r , controllers.GetGadgets, "read")
}

func GetGadgetTypes(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r , controllers.GetGadgetTypes, "read")
}

func AddGadgets(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r , controllers.AddGadgets, "write")
}

func GetMethods(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r , controllers.GetMethods, "write")
}

func AddMethod(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.SaveMethod, "write")
}

func UpdateMethod(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.SaveMethod, "write")
}

func DeleteMethod(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.DeleteMethod, "write")
}

func GetTimeseries(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.GetTimeseries, "read")
}

func GetDevices(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.GetDevices, "read")
}

func GetInSocket(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.HandleInSocket, "read")
}

func GetOutSocket(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.HandleOutSocket, "write")
}
