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

// NewCounterHandler - builds main http handler for api.CounterService implementation.
// If there is no a plan to log requests and responses, pass Logger as nil,
// otherwise make an adapter to expose rest.Logger interface.
func NewCounterHandler(baseURI string, service api.CyclicCounterService, l Logger) http.Handler {
	if service == nil {
		panic(errors.New("rest.NewHandler: CounterService is not implemented"))
	}
	r := mux.NewRouter()
	r.NotFoundHandler = handleNotFound(l)
	r.MethodNotAllowedHandler = handleMethodNotAllowed(l)

	{
		v1 := r.PathPrefix(baseURI).Subrouter()
		// v1.Use(logMiddleware())

		v1.NewRoute().
			Path("/getnumber/").
			Methods("GET").
			HandlerFunc(handleGetCounterValue(service, l))

		v1.NewRoute().
			Path("/incrementnumber/").
			Methods("POST").
			HandlerFunc(handleIncreaseCounter(service, l))

		v1.NewRoute().
			Path("/setsettings/{increment:[0-9]+}/{upper:[0-9]+}/").
			Methods("PUT").
			HandlerFunc(handleSetSettings(service, l))
	}

	// return logRequestMiddleware(l, r)
	return r
}

func httpStatusFactory(err error) int {
	switch e := err.(type) {
	case nil:
		return http.StatusOK
	case *api.Error:
		if e == nil {
			return http.StatusOK
		}
		if e.IsInternal() {
			return http.StatusInternalServerError
		}
		return http.StatusBadRequest
	default:
		return http.StatusServiceUnavailable
	}
}

func handleGetCounterValue(service api.CyclicCounterService, l Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, apiErr := service.GetCounterValue()
		status := httpStatusFactory(apiErr)
		if apiErr != nil {
			logError(l, status, formatRequest(r), formatError(apiErr.ExposeError()))
			response.HandleJSON(
				status,
				&response.Fail{response.ErrorDescription{0, fmt.Sprint(apiErr)}},
			)(w, r)
			return
		}
		logInfo(l, status, formatRequest(r))
		response.HandleJSON(status, &response.Success{Result: result})(w, r)
	}
}

func handleIncreaseCounter(service api.CyclicCounterService, l Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, apiErr := service.IncreaseCounter()
		status := httpStatusFactory(apiErr)
		if apiErr != nil {
			logError(l, status, formatRequest(r), formatError(apiErr.ExposeError()))
			response.HandleJSON(
				status,
				&response.Fail{response.ErrorDescription{0, fmt.Sprint(apiErr)}},
			)(w, r)
			return
		}
		logInfo(l, status, formatRequest(r))
		response.HandleJSON(status, &response.Success{Result: result})(w, r)
	}
}

func handleSetSettings(service api.CyclicCounterService, l Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		increment, err := strconv.Atoi(mux.Vars(r)["increment"])
		if err != nil {
			logError(l, http.StatusBadRequest, formatRequest(r), formatError(err))
			response.HandleJSON(
				http.StatusBadRequest,
				&response.Fail{response.ErrorDescription{0, fmt.Sprint("Invalid or bad increment")}},
			)(w, r)
			return
		}
		upper, err := strconv.Atoi(mux.Vars(r)["upper"])
		if err != nil {
			logError(l, http.StatusBadRequest, formatRequest(r), formatError(err))
			response.HandleJSON(
				http.StatusBadRequest,
				&response.Fail{response.ErrorDescription{0, fmt.Sprint("Invalid or bad upper limit value")}},
			)(w, r)
			return
		}
		// TODO Change URI to allow 3 parameters
		result, apiErr := service.SetCounterSettings(increment, 0, upper)
		status := httpStatusFactory(apiErr)
		if apiErr != nil {
			logError(l, status, formatRequest(r), formatError(apiErr.ExposeError()))
			response.HandleJSON(
				status,
				&response.Fail{response.ErrorDescription{0, fmt.Sprint(apiErr)}},
			)(w, r)
			return
		}
		logInfo(l, status, formatRequest(r))
		response.HandleJSON(status, &response.Success{Result: result})(w, r)
	}
}

func handleNotFound(l Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusNotFound
		logError(l, status, formatRequest(r))
		response.HandleJSON(
			status,
			&response.Fail{response.ErrorDescription{0, "Not Found"}},
		)(w, r)
	}
}

func handleMethodNotAllowed(l Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusMethodNotAllowed
		logError(l, status, formatRequest(r))
		// NOTE If the reason for this handler is HEAD request - gorilla.mux will not send response body to client!
		response.HandleJSON(
			status,
			&response.Fail{response.ErrorDescription{0, fmt.Sprintf("Method Not Allowed (%s)", r.Method)}},
		)(w, r)
	}
}

func logRequestMiddleware(l Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logInfo(l, formatRequest(r))
		next.ServeHTTP(w, r)
	})
}
