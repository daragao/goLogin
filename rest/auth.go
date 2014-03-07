package rest

import (
	//"../logger"
	"../auth"
	"github.com/ant0ine/go-json-rest"
	"net/http"
)

type Authentication struct {
}

func (authObj *Authentication) Logout(writer *rest.ResponseWriter, request *rest.Request) {
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
