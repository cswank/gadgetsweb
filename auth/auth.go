package auth

import (
	"log"
	"os"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bitbucket.org/cswank/gadgetsweb/models"
	"github.com/gorilla/securecookie"
)

type controller func(w http.ResponseWriter, r *http.Request) error

var (
	hashKey        = []byte(os.Getenv("HASH_KEY"))
	blockKey       = []byte(os.Getenv("BLOCK_KEY"))
	SecureCookie   = securecookie.New(hashKey, blockKey)
)

func CheckAuth(w http.ResponseWriter, r *http.Request, ctrl controller, permission string) {
	user, err := getUserFromCookie(r)
	if err == nil && user.IsAuthorized(permission) {
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
