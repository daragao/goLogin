package session

import (
	"github.com/gorilla/sessions"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))
var session_name = "geteonSession"

func Get(r *http.Request) (*sessions.Session, error) {
	thisSession, err := store.Get(r, session_name)
	thisSession.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	return thisSession, err
}

func Delete(thisSession *sessions.Session) {
	thisSession.Options = &sessions.Options{
		Path:   "/",
		MaxAge: -1,
	}
}
