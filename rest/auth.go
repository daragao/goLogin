package rest

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/daragao/goLogin/auth"
	"github.com/daragao/goLogin/logger"
	"github.com/daragao/goLogin/session"
	"github.com/daragao/goLogin/users"
	"net/http"
)

type Authentication struct {
}

func (authObj *Authentication) Login(writer rest.ResponseWriter, request *rest.Request) {
	userPostData := &Users{}
	err := request.DecodeJsonPayload(userPostData)
	if err != nil {
		rest.Error(writer, "Could not register user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := users.GetUserByUsername(userPostData.Username)
	isLoggedIn := auth.EqualPass(user.GetPassword(), userPostData.Password)
	if err != nil {
		rest.Error(writer, "Invalid login", http.StatusInternalServerError)
		return
	}
	if isLoggedIn {
		newSession, _ := session.Get(request.Request)
		newSession.Values["userId"] = user.Id
		newSession.Save(request.Request, writer.(http.ResponseWriter)) //writer.ResponseWriter)
		userRow, err := users.GetUserByID(user.Id)
		if err == nil {
			writer.WriteJson(userRow)
		} else {
			rest.Error(writer, "User not logged in: "+err.Error(), http.StatusNotFound)
		}
		//writer.WriteJson(newSession)
	} else {
		logger.ERRO.Println("Invalid login: ", user)
		rest.Error(writer, "Invalid login", http.StatusUnauthorized)
	}
}

func (authObj *Authentication) Logout(writer rest.ResponseWriter, request *rest.Request) {
	userSession, _ := session.Get(request.Request)
	session.Delete(userSession)
	userSession.Save(request.Request, writer.(http.ResponseWriter)) //.ResponseWriter)
	realm := "Administration"
	//writer.WriteJson(`{"status": "off"}`)
	Unauthorized(writer, realm)
}

/* PRE ROUTING METHODS */
type AuthError struct {
	isInvalid    bool
	errorMessage string
}

func (e *AuthError) Error() string {
	return e.errorMessage
}

func BasicAuthenticationLogin(writer rest.ResponseWriter,
	request *rest.Request) (authError *AuthError) {

	realm := "Administration"
	userId := "admin"
	password := "admin"

	authHeader := request.Header.Get("Authorization")
	if authHeader == "" {
		Unauthorized(writer, realm)
		authError = &AuthError{false, "Unauthorized user"}
		return
	}

	providedUserId, providedPassword, err := auth.DecodeBasicAuthHeader(authHeader)

	if err != nil {
		rest.Error(writer, "Invalid authentication", http.StatusBadRequest)
		authError = &AuthError{false, "Unauthorized user"}
		return
	}

	if !(providedUserId == userId && providedPassword == password) {
		Unauthorized(writer, realm)
		authError = &AuthError{false, "Invalid authentication"}
		return
	}

	return
}

func Unauthorized(writer rest.ResponseWriter, realm string) {
	//writer.Header().Set("WWW-Authenticate", "Basic realm="+realm)
	rest.Error(writer, "Not Authorized", http.StatusUnauthorized)
}
