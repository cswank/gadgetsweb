package main

import (
	"log"
	"fmt"
	"os"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bitbucket.com/cswank/gadgetsweb/controllers"
	"bitbucket.com/cswank/gadgetsweb/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var (
	hashKey        = []byte(os.Getenv("HASH_KEY"))
	blockKey       = []byte(os.Getenv("BLOCK_KEY"))
	//hashKey        = []byte("very-secret")
	//blockKey       = []byte("a-lot-secrettttt")
	SecureCookie   = securecookie.New(hashKey, blockKey)
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", Login).Methods("POST")
	r.HandleFunc("/logout", Logout).Methods("POST")
	r.HandleFunc("/gadgets", GetGadgets).Methods("GET")
	
	r.HandleFunc("/history/locations/summary", GetSummary).Methods("GET")
	r.HandleFunc("/gadgets/{name}/methods", GetMethods).Methods("GET")
	r.HandleFunc("/gadgets/{name}/methods", AddMethod).Methods("POST")
	r.HandleFunc("/gadgets/{name}/methods/{methodId}", UpdateMethod).Methods("PUT")
	r.HandleFunc("/history/locations/{location}/directions/{direction}/devices/{device}", GetTimeseries).Methods("GET")
	r.HandleFunc("/socket", GetSocket)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/app")))
	http.Handle("/", r)
	fmt.Println("listening on 0.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:  "gadgets",
		Value: "",
		Path:  "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request 1", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, user)
	if err != nil {
		http.Error(w, "bad request 2", http.StatusBadRequest)
		return
	}
	goodPassword, err := user.CheckPassword()
	if !goodPassword {
		http.Error(w, "bad request 3", http.StatusBadRequest)
		return 
	}
	value := map[string]string{
		"user": user.Username,
	}
	encoded, _ := SecureCookie.Encode("gadgets", value)
	cookie := &http.Cookie{
		Name:  "gadgets",
		Value: encoded,
		Path:  "/",
		HttpOnly: false,
	}
	http.SetCookie(w, cookie)
}

func GetGadgets(w http.ResponseWriter, r *http.Request) {
	controllers.GetGadgets(w, r)
}

func GetMethods(w http.ResponseWriter, r *http.Request) {
	err := controllers.GetMethods(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AddMethod(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromCookie(r)
	if err == nil && user.IsAuthorized() {
		err = controllers.SaveMethod(w, r)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
	}
}

func UpdateMethod(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update method")
	user, err := getUserFromCookie(r)
	if err == nil && user.IsAuthorized() {
		controllers.SaveMethod(w, r)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
	}
}

func GetTimeseries(w http.ResponseWriter, r *http.Request) {
	controllers.GetTimeseries(w, r)
}

func GetSummary(w http.ResponseWriter, r *http.Request) {
	controllers.GetSummary(w, r)
}

func GetSocket(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromCookie(r)
	if err == nil && user.IsAuthorized() {
		controllers.HandleSocket(w, r)
	} else {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
	}
}

func getUserFromCookie(r *http.Request) (*models.User, error) {
	user := &models.User{}
	cookie, err := r.Cookie("gadgets")
	if err == nil {
		m := map[string]string{}
		err = SecureCookie.Decode("gadgets", cookie.Value, &m)
		if err == nil {
			user.Username = m["user"]
		}
	}
	return user, err
}
