package goyave

import (
	"net/http"

	"github.com/System-Glitch/goyave/config"
	"github.com/System-Glitch/goyave/helpers/response"
	"github.com/gorilla/mux"
)

// Router registers routes to be matched and dispatches a handler.
type Router struct {
	muxRouter *mux.Router
}

func newRouter() *Router {
	muxRouter := mux.NewRouter()
	muxRouter.Schemes(config.Get("protocol").(string))
	// TODO recover middleware
	return &Router{muxRouter: muxRouter}
}

// Subrouter create a new sub-router from this router.
// Use subrouters to create route groups and to apply middlewares to multiple routes.
func (r *Router) Subrouter(prefix string) *Router {
	return &Router{muxRouter: r.muxRouter.PathPrefix(prefix).Subrouter()}
}

// Middleware apply one or more middleware(s) to the route group.
func (r *Router) Middleware(middlewares ...func(http.Handler) http.Handler) {
	// TODO implement middleware
}

// Route register a new route.
func (r *Router) Route(method string, endpoint string, handler func(http.ResponseWriter, *Request), requestGenerator func() *Request) {
	r.muxRouter.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		// TODO handle url params
		req, ok := requestHandler(w, r, requestGenerator)
		if ok {
			handler(w, req)
		}
	}).Methods(method)
}

func requestHandler(w http.ResponseWriter, r *http.Request, requestGenerator func() *Request) (*Request, bool) {
	var request *Request
	if requestGenerator != nil {
		request = requestGenerator()
	} else {
		request = &Request{}
	}
	request.httpRequest = r
	errsBag := request.validate()
	if errsBag == nil {
		return request, true
	}

	var code int
	if isRequestMalformed(errsBag) {
		code = http.StatusBadRequest
	} else {
		code = http.StatusUnprocessableEntity
	}
	response.JSON(w, code, errsBag)

	return nil, false
}
