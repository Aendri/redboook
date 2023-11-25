package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Q is a {string->string} map for holding HTTP query parameters when setting up routes.
type Q map[string]string

type handlerFunc func(http.ResponseWriter, *http.Request)

// RegisterHandler function registers a provided handler function to handle requests to the given path
// with the given HTTP method. Additionally it registers any optional middleware functions with the
// created sub-router.
// Return the new gorilla Route (which can be specified further).
func RegisterHandler(
	router *mux.Router, // parent Router
	method string, // HTTP method
	path string, // URL path
	query Q, // URL query parameters (can be nil)
	handler handlerFunc, // handler function
	middlewares ...mux.MiddlewareFunc, // optional middlewares
) *mux.Route {
	// new route for the path.
	route := router.Path(path)

	// specify the method unless it's an empty string or "*".
	if method != "" && method != "*" {
		route = route.Methods(method)
	}

	// specify the query params (if any)
	for key, value := range query {
		route.Queries(key, value)
	}

	// specify middlewares (if any) on a dedicated sub-Router
	if len(middlewares) > 0 {
		router = route.Subrouter()
		router.Use(middlewares...)
		route = router.NewRoute()
	}

	// setup the handler which gets us a new route
	if handler != nil {
		route.HandlerFunc(handler)
	}

	return route
}
