package rest

import (
	"fmt"
	"net/http"

	"github.com/wtask-go/auracounter/internal/httpcore/response"

	"github.com/wtask-go/auracounter/internal/api"

	"github.com/gorilla/mux"
)

// NewCounterHandler - builds main http handler for api.CounterService implementation
func NewCounterHandler(service api.CounterService) http.Handler {
	// if service == nil {
	// 	panic(errors.New("rest.NewHandler: CounterService is not implemented"))
	// }
	r := mux.NewRouter()
	r.NotFoundHandler = handleNotFound()
	r.MethodNotAllowedHandler = handleMethodNotAllowed()

	{
		v1 := r.PathPrefix("/counter/v1/").Subrouter()

		v1.NewRoute().
			Path("/getnumber/").
			Methods("GET").
			HandlerFunc(handleGetNumber(service))

		v1.NewRoute().
			Path("/incrementnumber/").
			Methods("POST").
			HandlerFunc(handleIncrementNumber(service))

		v1.NewRoute().
			Path("/setsettings/{delta:[0-9]+}/{max:[0-9]+}/").
			Methods("PUT").
			HandlerFunc(handleSetSettings(service))
	}

	return r
}

func handleGetNumber(service api.CounterService) http.HandlerFunc {
	return response.HandleJSON(
		http.StatusNotImplemented,
		&response.Fail{
			response.ErrorDescription{0, "Not implemented"},
		},
	)

	// return func(w http.ResponseWriter, r *http.Request) {

	// }
}

func handleIncrementNumber(service api.CounterService) http.HandlerFunc {
	return response.HandleJSON(
		http.StatusNotImplemented,
		&response.Fail{
			response.ErrorDescription{0, "Not implemented"},
		},
	)
	// return func(w http.ResponseWriter, r *http.Request) {

	// }
}

func handleSetSettings(service api.CounterService) http.HandlerFunc {
	return response.HandleJSON(
		http.StatusNotImplemented,
		&response.Fail{
			response.ErrorDescription{0, "Not implemented"},
		},
	)

	// return func(w http.ResponseWriter, r *http.Request) {

	// }
}

func handleNotFound() http.HandlerFunc {
	return response.HandleJSON(
		http.StatusNotFound,
		&response.Fail{
			response.ErrorDescription{0, "Resource not found"},
		},
	)
}

func handleMethodNotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.HandleJSON(
			http.StatusMethodNotAllowed,
			&response.Fail{
				response.ErrorDescription{0, fmt.Sprintf("Request method is not allowed (%s)", r.Method)},
			},
		)(w, r)
	}
}
