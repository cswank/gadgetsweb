package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bitbucket.com/cswank/gadgetsweb/controllers"
	"bitbucket.com/cswank/gadgetsweb/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var (
	hashKey        = []byte("very-secret")
	blockKey       = []byte("a-lot-secrettttt")
	SecureCookie   = securecookie.New(hashKey, blockKey)
)

func main() {
	
	r := mux.NewRouter()
	r.HandleFunc("/", GetHome).Methods("GET")
	r.HandleFunc("/login", DoLogin).Methods("POST")
	r.HandleFunc("/history/locations/summary", GetSummary).Methods("GET")
	r.HandleFunc("/history/locations/{location}/directions/{direction}/devices/{device}", GetTimeseries).Methods("GET")
	r.HandleFunc("/socket", GetSocket)

	http.Handle("/", r)
	fmt.Println("listening on 0.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	
}

func DoLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("logging in")
	w.Header().Set("Access-Control-Allow-Origin", "http://gadgets.dyndns-ip.com")
	w.Header().Set("Access-Control-Allow-Headers","X-Requested-With");
	w.Header().Set("Access-Control-Allow-Methods","GET, POST");
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	user := &models.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, user)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	goodPassword, err := user.CheckPassword()
	if !goodPassword {
		http.Error(w, "bad request", http.StatusBadRequest)
		return 
	}
	value := map[string]string{
		"user": user.Username,
	}
	fmt.Println(value)
	encoded, _ := SecureCookie.Encode("gadgets", value)
	cookie := &http.Cookie{
		Name:  "gadgets",
		Value: encoded,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}

func GetTimeseries(w http.ResponseWriter, r *http.Request) {
	controllers.GetTimeseries(w, r)
}

func GetSummary(w http.ResponseWriter, r *http.Request) {
	controllers.GetSummary(w, r)
}

func GetSocket(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromCookie(r)
	fmt.Println("get socket", user, err)
	if err == nil && user.IsAuthorized() {
		controllers.HandleSocket(w, r)
	} else {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
	}
}

func getUserFromCookie(r *http.Request) (*models.User, error) {
	user := &models.User{}
	cookie, err := r.Cookie("gadgets")
	fmt.Println("getuserfromcookie", cookie, err)
	m := map[string]string{}
	if err == nil {
		err = SecureCookie.Decode("gadgets", cookie.Value, &m)
		fmt.Println("getuserfromcookie", m, err)
		if err == nil {
			user.Username = m["user"]
		}
	}
	return user, err
}
