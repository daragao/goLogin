package rest

import (
	//"../logger"
	"../auth"
	"../session"
	"../users"
	"github.com/ant0ine/go-json-rest"
	"net/http"
)

type Authentication struct {
}

func (authObj *Authentication) Login(writer *rest.ResponseWriter, request *rest.Request) {
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
		newSession.Save(request.Request, writer.ResponseWriter)
		writer.WriteJson(userPostData)
	} else {
		rest.Error(writer, "Invalid login", http.StatusUnauthorized)
	}
}

func (authObj *Authentication) Logout(writer *rest.ResponseWriter, request *rest.Request) {
	userSession, _ := session.Get(request.Request)
	session.Delete(userSession)
	userSession.Save(request.Request, writer.ResponseWriter)
	realm := "Administration"
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

func BasicAuthenticationLogin(writer *rest.ResponseWriter,
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

func Unauthorized(writer *rest.ResponseWriter, realm string) {
	writer.Header().Set("WWW-Authenticate", "Basic realm="+realm)
	rest.Error(writer, "Not Authorized", http.StatusUnauthorized)
}
