package cleanapi

import (
	"clean-rest-api/types"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type CleanApi struct {
	Context types.Context
	Router  *mux.Router
}

func New() *CleanApi {

	return &CleanApi{
		Router: mux.NewRouter(),
	}
}

func (c *CleanApi) Get(path string, httpHandlerFunc types.HttpHandlerFunc) {
	c.Router.HandleFunc(path, WrapHandler(httpHandlerFunc)).Methods("GET")
}

func WrapHandler(hf types.HttpHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &types.Context{
			Writer:  w,
			Request: r,
		}

		New().Context = *ctx

		if err := hf(ctx); err != nil {
			log.Println("error handled")
		}
	}
}
