package json

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	//Http Status Codes for errors
	Code int `json:"status"`
	//Associated Http Error message
	Message string `json:"message"`
}

// writeJsonResponse writes a json repsonse with a given status code and
// encode the data payload.
// 'data' being any allows it to be 'templated'
func Write(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

// Should json decoding of a request fail, provide an error as a message
func WriteError(w http.ResponseWriter, status int, message string) {
	error := ErrorResponse{
		Code:    status,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(error)
}
func Reader(r *http.Request, data any) error {
	decoder := json.NewDecoder(r.Body)
	//making requests strict. Server throws error
	//when request body is unexpected
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

// Global HTTP handler errors. These are typically shared between all modules
var (
	//Used to return errors that our API consumer needs to resolve
	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		WriteError(w, http.StatusBadRequest, err.Error())
	}
	//This is for generic errors that may occur in the server and are not a
	// result of the user's input
	InternalErrorHandler = func(w http.ResponseWriter, err error) {

		WriteError(w, http.StatusInternalServerError, "Internal Server Error Occured")
	}
)
