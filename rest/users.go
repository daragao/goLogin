package rest

import (
	//"../logger"
	"github.com/ant0ine/go-json-rest"
)

type Users struct {
	//cache map[string]*users.Users
	Name string
}

func (self *Users) GetAll(w *rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(&Users{"tchcahhsdh"})
}
