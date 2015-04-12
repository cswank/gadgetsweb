package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cswank/gadgetsweb/auth"
	"github.com/cswank/gadgetsweb/controllers"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var (
	hashKey      []byte
	blockKey     []byte
	cert         string
	key          string
	SecureCookie *securecookie.SecureCookie
)

func init() {
	hashKey = []byte(os.Getenv("GADGETS_HASH_KEY"))
	blockKey = []byte(os.Getenv("GADGETS_BLOCK_KEY"))
	cert = os.Getenv("GADGETS_CERT")
	key = os.Getenv("GADGETS_KEY")
	if len(cert) == 0 || len(key) == 0 {
		log.Fatal("you must set CERT and KEY env vars")
	}
	SecureCookie = securecookie.New(hashKey, blockKey)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/login", auth.Login).Methods("POST")
	r.HandleFunc("/api/logout", auth.Logout).Methods("POST")
	r.HandleFunc("/api/socket", GetSocket).Methods("GET")
	r.HandleFunc("/api/gadgets", GetGadgets).Methods("GET")
	r.HandleFunc("/api/gadgets", AddGadgets).Methods("POST")
	r.HandleFunc("/api/gadgets/types", GetGadgetTypes).Methods("GET")
	r.HandleFunc("/api/gadgets/{name}/methods", GetMethods).Methods("GET")
	r.HandleFunc("/api/gadgets/{name}/status", GetStatus).Methods("GET")
	r.HandleFunc("/api/gadgets/{name}/commands", SendCommand).Methods("POST")
	r.HandleFunc("/api/gadgets/{name}/methods", AddMethod).Methods("POST")
	r.HandleFunc("/api/gadgets/{name}/notes", GetNotes).Methods("GET")
	r.HandleFunc("/api/gadgets/{name}/notes", AddNotes).Methods("POST")
	r.HandleFunc("/api/gadgets/{name}/methods/{methodId}", UpdateMethod).Methods("PUT")
	r.HandleFunc("/api/gadgets/{name}/methods/{methodId}", DeleteMethod).Methods("DELETE")
	r.HandleFunc("/api/history/gadgets/{gadget}/devices", GetDevices).Methods("GET")
	r.HandleFunc("/api/history/gadgets/{gadget}/locations/{location}/devices/{device}", GetTimeseries).Methods("GET")

	http.Handle("/", r)
	fmt.Println("listening on 0.0.0.0:443")
	err := http.ListenAndServeTLS(":443", cert, key, nil)
	log.Println(err)
}

func GetGadgets(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.GetGadgets, "read")
}

func GetGadgetTypes(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.GetGadgetTypes, "read")
}

func AddGadgets(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.AddGadgets, "write")
}

func GetMethods(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.GetMethods, "write")
}

func AddMethod(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.SaveMethod, "write")
}

func GetNotes(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.GetNotes, "write")
}

func AddNotes(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.SaveNote, "write")
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

func GetStatus(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.GetStatus, "read")
}

func SendCommand(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.SendCommand, "read")
}

func GetSocket(w http.ResponseWriter, r *http.Request) {
	auth.CheckAuth(w, r, controllers.HandleSocket, "read")
}
