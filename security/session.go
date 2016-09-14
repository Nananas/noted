package security

import (
	"net/http"
	"strconv"
	"time"
)

var ACCOUNTS_FILEPATH = "./accounts"

var cache map[string]string

func Validate(user string, pass string) bool {
	if cache[user] != "" {
		return cache[user] == SaltHash(pass)
	}

	return false
}

// A session value is created using the username and the current day of the year
//
func getSessionValue(user string) string {
	d := time.Now().YearDay()
	s := strconv.Itoa(d) + user

	return SaltHash(s)
}

// Returns the session value of yesterday
//
func oldSessionValue(user string) string {
	d := time.Now().Add(-24 * time.Hour).YearDay()
	s := strconv.Itoa(d) + user

	return SaltHash(s)
}

func ValidSession(req *http.Request) bool {
	cookie, err1 := req.Cookie("session")
	user, err2 := req.Cookie("user")

	if err1 != nil || err2 != nil || (cookie.Value != getSessionValue(user.Value) && cookie.Value != oldSessionValue(user.Value)) {
		return false
	}

	return true
}

// A session is valid for two days
//
func StartSession(rw http.ResponseWriter, user string) {

	expire := time.Now().AddDate(0, 0, 1)

	sessioncookie := http.Cookie{
		Name:    "session",
		Value:   getSessionValue(user),
		Expires: expire,
	}

	usercookie := http.Cookie{
		Name:    "user",
		Value:   user,
		Expires: expire,
	}

	http.SetCookie(rw, &sessioncookie)
	http.SetCookie(rw, &usercookie)

}

func UpdateSession(rw http.ResponseWriter, req *http.Request) {

	user := req.Header.Get("user")
	if user == "" {
		return
	}

	expire := time.Now().AddDate(0, 0, 1)

	sessioncookie := http.Cookie{
		Name:    "session",
		Value:   getSessionValue(user),
		Expires: expire,
	}

	usercookie := http.Cookie{
		Name:    "user",
		Value:   user,
		Expires: expire,
	}

	http.SetCookie(rw, &sessioncookie)
	http.SetCookie(rw, &usercookie)
}
