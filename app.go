package main

import (
	"log"
	"fmt"
	"os"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bitbucket.org/cswank/gadgetsweb/controllers"
	"bitbucket.org/cswank/gadgetsweb/models"
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
	r.HandleFunc("/api/socket", GetSocket)
	r.HandleFunc("/api/gadgets", GetGadgets).Methods("GET")
	r.HandleFunc("/api/gadgets", AddGadgets).Methods("POST")
	r.HandleFunc("/api/gadgets/types", GetGadgetTypes).Methods("GET")
	r.HandleFunc("/api/gadgets/{name}/methods", GetMethods).Methods("GET")
	r.HandleFunc("/api/gadgets/{name}/methods", AddMethod).Methods("POST")
	r.HandleFunc("/api/gadgets/{name}/methods/{methodId}", UpdateMethod).Methods("PUT")
	r.HandleFunc("/api/gadgets/{name}/methods/{methodId}", DeleteMethod).Methods("DELETE")
	r.HandleFunc("/api/history/gadgets/{gadget}/devices", GetDevices).Methods("GET")
	r.HandleFunc("/api/history/gadgets/{gadget}/locations/{location}/devices/{device}", GetTimeseries).Methods("GET")
	r.HandleFunc("/recipes/{name}", GetRecipe).Methods("GET")
	
	http.Handle("/", r)
	fmt.Println("listening on 0.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}

func GetGadgets(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r , controllers.GetGadgets)
}

func GetGadgetTypes(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r , controllers.GetGadgetTypes)
}

func AddGadgets(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r , controllers.AddGadgets)
}

func GetMethods(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r , controllers.GetMethods)
}

func AddMethod(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.SaveMethod)
}

func UpdateMethod(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.SaveMethod)
}

func DeleteMethod(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.DeleteMethod)
}

func GetRecipe(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.GetRecipe)
}

func GetTimeseries(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.GetTimeseries)
}

func GetDevices(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.GetDevices)
}

func GetSocket(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.HandleSocket)
}
