package response

import (
	"encoding/json"
	"net/http"
)

// HandleJSON - retturn http-handler function to convert response into JSON format
// on the final stage of processing client request.
func HandleJSON(status int, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(&data)
	}
}
