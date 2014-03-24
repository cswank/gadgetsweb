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
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var (
	hashKey        = []byte(os.Getenv("HASH_KEY"))
	blockKey       = []byte(os.Getenv("BLOCK_KEY"))
	SecureCookie   = securecookie.New(hashKey, blockKey)
)

type controller func(w http.ResponseWriter, r *http.Request) error

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", Login).Methods("POST")
	r.HandleFunc("/logout", Logout).Methods("POST")
	r.HandleFunc("/socket", GetSocket)
	r.HandleFunc("/recipes/{name}", GetRecipe).Methods("GET")
	r.HandleFunc("/gadgets", GetGadgets).Methods("GET")
	r.HandleFunc("/gadgets/{name}/methods", GetMethods).Methods("GET")
	r.HandleFunc("/gadgets/{name}/methods", AddMethod).Methods("POST")
	r.HandleFunc("/gadgets/{name}/methods/{methodId}", UpdateMethod).Methods("PUT")
	r.HandleFunc("/gadgets/{name}/methods/{methodId}", DeleteMethod).Methods("DELETE")
	r.HandleFunc("/history/gadgets/{gadget}/devices", GetDevices).Methods("GET")
	r.HandleFunc("/history/gadgets/{gadget}/locations/{location}/devices/{device}", GetTimeseries).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(os.Getenv("GADGETS_STATIC"))))
	
	http.Handle("/", r)
	fmt.Println("listening on 0.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}

func GetGadgets(w http.ResponseWriter, r *http.Request) {
	checkAuth(w, r , controllers.GetGadgets)
}

func GetMethods(w http.ResponseWriter, r *http.Request) {
	checkAuth(w, r , controllers.GetMethods)
}

func AddMethod(w http.ResponseWriter, r *http.Request) {
	checkAuth(w, r, controllers.SaveMethod)
}

func UpdateMethod(w http.ResponseWriter, r *http.Request) {
	checkAuth(w, r, controllers.SaveMethod)
}

func DeleteMethod(w http.ResponseWriter, r *http.Request) {
	checkAuth(w, r, controllers.DeleteMethod)
}

func GetRecipe(w http.ResponseWriter, r *http.Request) {
	checkAuth(w, r, controllers.GetRecipe)
}

func GetTimeseries(w http.ResponseWriter, r *http.Request) {
	checkAuth(w, r, controllers.GetTimeseries)
}

func GetDevices(w http.ResponseWriter, r *http.Request) {
	checkAuth(w, r, controllers.GetDevices)
}

func GetSocket(w http.ResponseWriter, r *http.Request) {
	checkAuth(w, r, controllers.HandleSocket)
}

func checkAuth(w http.ResponseWriter, r *http.Request, ctrl controller) {
	user, err := getUserFromCookie(r)
	if err == nil && user.IsAuthorized() {
		err = ctrl(w, r)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
	fmt.Println(goodPassword, err)
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
