package rest

import (
	"../db"
	"../logger"
	"github.com/ant0ine/go-json-rest"
	"net/http"
)

func Start() {
	logger.Init()

	StartDB()
	StartServer()
}

func StartDB() {
	logger.INFO.Println("Starting DB...")
	db.CreateTables()
}

func StartServer() {
	logger.INFO.Println("Starting server...")

	rootUri := "/api/v1"

	handler := rest.ResourceHandler{
		PreRoutingMiddleware:     PreRoutingMiddleware,
		EnableRelaxedContentType: true,
	}

	users := Users{}
	auth := Authentication{}

	handler.SetRoutes(
		//login and create session!
		rest.RouteObjectMethod("GET", rootUri+"/users", &users, "GetAllUsers"),
		rest.RouteObjectMethod("GET", rootUri+"/users/:id", &users, "GetUserByID"),
		rest.RouteObjectMethod("POST", rootUri+"/users", &users, "RegisterUser"),
		rest.RouteObjectMethod("GET", rootUri+"/logout", &auth, "Logout"),
	)

	http.ListenAndServe(":4000", &handler)
	/*  // IF WE WANT TO USE NGINX
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		logger.ERRO.Println(err)
	}
	err = fcgi.Serve(listener, &handler)
	if err != nil {
		logger.ERRO.Println(err)
	}
	*/
}

func PreRoutingMiddleware(handler rest.HandlerFunc) rest.HandlerFunc {

	return func(writer *rest.ResponseWriter, request *rest.Request) {

		//authErr := BasicAuthenticationLogin(writer, request)
		//if authErr != nil {
		//logger.ERRO.Println(authErr)
		//	return
		//}

		handler(writer, request)
	}
}
