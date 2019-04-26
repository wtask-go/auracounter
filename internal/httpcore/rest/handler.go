package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/wtask-go/auracounter/internal/httpcore/response"

	"github.com/wtask-go/auracounter/internal/api"

	"github.com/gorilla/mux"
)

// NewCounterHandler - builds main http handler for api.CounterService implementation
func NewCounterHandler(baseURI string, service api.CyclicCounterService) http.Handler {
	if service == nil {
		panic(errors.New("rest.NewHandler: CounterService is not implemented"))
	}
	r := mux.NewRouter()
	r.NotFoundHandler = handleNotFound()
	r.MethodNotAllowedHandler = handleMethodNotAllowed()

	{
		v1 := r.PathPrefix(baseURI).Subrouter()

		v1.NewRoute().
			Path("/getnumber/").
			Methods("GET").
			HandlerFunc(handleGetCounterValue(service))

		v1.NewRoute().
			Path("/incrementnumber/").
			Methods("POST").
			HandlerFunc(handleIncreaseCounter(service))

		v1.NewRoute().
			Path("/setsettings/{delta:[0-9]+}/{max:[0-9]+}/").
			Methods("PUT").
			HandlerFunc(handleSetSettings(service))
	}

	return r
}

func httpStatusFactory(err error) int {
	switch e := err.(type) {
	case nil:
		return http.StatusOK
	case *api.Error:
		if !e.IsInternal() {
			return http.StatusBadRequest
		}
		return http.StatusInternalServerError
	default:
		return http.StatusServiceUnavailable
	}
}

func handleGetCounterValue(service api.CyclicCounterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := service.GetCounterValue()
		status := httpStatusFactory(err)
		if err != nil {
			response.HandleJSON(status, &response.Fail{response.ErrorDescription{0, fmt.Sprint(err)}})(w, r)
			return
		}
		response.HandleJSON(status, &response.Success{Result: result})(w, r)
	}
}

func handleIncreaseCounter(service api.CyclicCounterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := service.IncreaseCounter()
		status := httpStatusFactory(err)
		if err != nil {
			response.HandleJSON(status, &response.Fail{response.ErrorDescription{0, fmt.Sprint(err)}})(w, r)
			return
		}
		response.HandleJSON(status, &response.Success{Result: result})(w, r)
	}
}

func handleSetSettings(service api.CyclicCounterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		delta, err := strconv.Atoi(mux.Vars(r)["delta"])
		if err != nil {
			response.HandleJSON(
				http.StatusBadRequest,
				&response.Fail{response.ErrorDescription{0, fmt.Sprint("Invalid or bad delta param")}},
			)(w, r)
			return
		}
		max, err := strconv.Atoi(mux.Vars(r)["max"])
		if err != nil {
			response.HandleJSON(
				http.StatusBadRequest,
				&response.Fail{response.ErrorDescription{0, fmt.Sprint("Invalid or bad max param")}},
			)(w, r)
			return
		}
		// TODO Change URI to allow 3 parameters
		result, err := service.SetCounterSettings(delta, 0, max)
		status := httpStatusFactory(err)
		if err != nil {
			response.HandleJSON(
				status,
				&response.Fail{response.ErrorDescription{0, fmt.Sprint(err)}},
			)(w, r)
			return
		}
		response.HandleJSON(status, &response.Success{Result: result})(w, r)
	}
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
			&response.Fail{response.ErrorDescription{0, fmt.Sprintf("Request method is not allowed (%s)", r.Method)}},
		)(w, r)
	}
}
