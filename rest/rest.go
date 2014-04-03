package rest

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/daragao/goLogin/db"
	"github.com/daragao/goLogin/logger"
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
		//PreRoutingMiddleware:     PreRoutingMiddleware,
		PreRoutingMiddlewares: []rest.Middleware{
			&MyCorsMiddleware{},
		},
		EnableRelaxedContentType: true,
	}

	users := Users{}
	auth := Authentication{}

	handler.SetRoutes(
		//login and create session!
		rest.RouteObjectMethod("GET", rootUri+"/users", &users, "GetAllUsers"),
		rest.RouteObjectMethod("GET", rootUri+"/users/:id", &users, "GetUserByID"),
		rest.RouteObjectMethod("POST", rootUri+"/users", &users, "RegisterUser"),
		rest.RouteObjectMethod("GET", rootUri+"/login", &users, "GetCurrentUser"),
		rest.RouteObjectMethod("POST", rootUri+"/login", &auth, "Login"),
		rest.RouteObjectMethod("PUT", rootUri+"/login", &auth, "Login"),
		rest.RouteObjectMethod("DELETE", rootUri+"/login", &auth, "Logout"),
		rest.RouteObjectMethod("GET", rootUri+"/logout", &auth, "Logout"),
	)

	http.ListenAndServe(":"+os.Getenv("PORT"), &handler)
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

type MyCorsMiddleware struct{}

//func PreRoutingMiddleware(handler rest.HandlerFunc) rest.HandlerFunc {
func (mw *MyCorsMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	return func(writer rest.ResponseWriter, request *rest.Request) {

		corsInfo := request.GetCorsInfo()

		//authErr := BasicAuthenticationLogin(writer, request)
		//if authErr != nil {
		//logger.ERRO.Println(authErr)
		//	return
		//}

		if !corsInfo.IsCors {
			// continure, execute the wrapped middleware
			handler(writer, request)
			return
		}

		// Validate the Origin
		// More sophisticated validations can be implemented, regexps, DB lookups, ...
		if corsInfo.Origin != "http://localhost:9000" {
			//rest.Error(writer, "Invalid Origin", http.StatusForbidden)
			//return
		}

		if corsInfo.IsPreflight {
			// check the request methods
			allowedMethods := map[string]bool{
				"GET":    true,
				"POST":   true,
				"PUT":    true,
				"DELETE": true,
			}
			if !allowedMethods[corsInfo.AccessControlRequestMethod] {
				rest.Error(writer, "Invalid Preflight Request", http.StatusForbidden)
				return
			}
			// check the request headers
			allowedHeaders := map[string]bool{
				"Accept":          true,
				"Content-Type":    true,
				"X-Custom-Header": true,
			}
			for _, requestedHeader := range corsInfo.AccessControlRequestHeaders {
				if !allowedHeaders[requestedHeader] {
					rest.Error(writer, "Invalid Preflight Request", http.StatusForbidden)
					return
				}
			}

			for allowedMethod, _ := range allowedMethods {
				writer.Header().Add("Access-Control-Allow-Methods", allowedMethod)
			}
			for allowedHeader, _ := range allowedHeaders {
				writer.Header().Add("Access-Control-Allow-Headers", allowedHeader)
			}
			writer.Header().Set("Access-Control-Allow-Origin", corsInfo.Origin)
			writer.Header().Set("Access-Control-Allow-Credentials", "true")
			writer.Header().Set("Access-Control-Max-Age", "3600")
			writer.WriteHeader(http.StatusOK)
			return
		} else {
			writer.Header().Set("Access-Control-Expose-Headers", "X-Powered-By")
			writer.Header().Set("Access-Control-Allow-Origin", corsInfo.Origin)
			writer.Header().Set("Access-Control-Allow-Credentials", "true")
			// continure, execute the wrapped middleware
			handler(writer, request)
			return
		}

		//handler(writer, request)
	}
}
