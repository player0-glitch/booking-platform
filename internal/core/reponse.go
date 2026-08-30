package response

import (
	"encoding/json"
	"net/http"
)

// JSON writes a json repsonse with a given status code and
// encode the data payload.
// 'data' being any allows it to be 'templated'
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		errEncoding := json.NewEncoder(w).Encode(data)
		if errEncoding != nil {
			//fallback onto this error if json encoding fails
			http.Error(w, errEncoding.Error(), http.StatusInternalServerError)
		}
	}
}

// Should json decoding of a request fail, provide an error as a message
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{"error": message})
}
