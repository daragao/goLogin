package rest

import (
	"../session"
	"../users"
	"github.com/ant0ine/go-json-rest"
	"net/http"
	"strconv"
)

type Users struct {
	Username string
	Password string
	Offset   int
}

func (self *Users) GetCurrentUser(w *rest.ResponseWriter, r *rest.Request) {
	currSession, _ := session.Get(r.Request)
	userId := currSession.Values["userId"]
	if userId != nil {
		userRow, err := users.GetUserByID(userId.(int))
		if err == nil {
			w.WriteJson(userRow)
			return
		} else {
			rest.Error(w, "User not logged in: "+err.Error(), http.StatusNotFound)
			return
		}
	}
	rest.Error(w, "User not logged in: failed to get session", http.StatusNotFound)
}

func (self *Users) GetUserByID(w *rest.ResponseWriter, r *rest.Request) {
	idStr := r.PathParam("id")
	id, err := strconv.Atoi(idStr)
	userRow, err := users.GetUserByID(id)
	if err == nil {
		w.WriteJson(userRow)
	} else {
		rest.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
	}
}

func (self *Users) GetAllUsers(w *rest.ResponseWriter, r *rest.Request) {
	userRows, err := users.GetAllUsers(100, 0)
	if err == nil {
		w.WriteJson(userRows)
	} else {
		rest.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
	}
}

/* TODO: add some security to this thing! */
func (self *Users) RegisterUser(w *rest.ResponseWriter, r *rest.Request) {
	userStruct := Users{}
	err := r.DecodeJsonPayload(&userStruct)
	if err != nil {
		rest.Error(w, "Could not register user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userRow, err := users.RegisterUser(userStruct.Username, userStruct.Password)
	if err == nil {
		w.WriteJson(userRow)
	} else {
		rest.Error(w, "Could not register user: "+err.Error(), http.StatusInternalServerError)
	}
}
