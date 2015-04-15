package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/cswank/gadgetsweb/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

type controller func(w http.ResponseWriter, r *http.Request, u *models.User, vars map[string]string) error

var (
	hashKey      = []byte(os.Getenv("GADGETS_HASH_KEY"))
	blockKey     = []byte(os.Getenv("GADGETS_BLOCK_KEY"))
	SecureCookie = securecookie.New(hashKey, blockKey)
)

func CheckAuth(w http.ResponseWriter, r *http.Request, ctrl controller, permission string) {
	user, err := getUserFromCookie(r)
	if err == nil && user.IsAuthorized(permission) {
		vars := mux.Vars(r)
		err = ctrl(w, r, user, vars)
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
		Name:   "gadgets",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("login")
	user := &models.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request 1", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, user)
	fmt.Println("user", user)
	if err != nil {

		http.Error(w, "bad request 2", http.StatusBadRequest)
		return
	}
	goodPassword, err := user.CheckPassword()
	fmt.Println("goodPw", goodPassword)
	if !goodPassword {
		log.Println(err)
		http.Error(w, "bad request 3", http.StatusBadRequest)
		return
	}
	value := map[string]string{
		"user": user.Username,
	}

	encoded, err := SecureCookie.Encode("gadgets", value)
	fmt.Println("user:", user, err)
	cookie := &http.Cookie{
		Name:     "gadgets",
		Value:    encoded,
		Path:     "/",
		HttpOnly: false,
	}
	fmt.Println("cookie", cookie, encoded)
	http.SetCookie(w, cookie)
}
