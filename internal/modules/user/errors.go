package user

import (
	"errors"
)

// Handler/Controller errors
var (
	errJSONDecoder    = errors.New("failed to decode json")
	errParse          = errors.New("failed to parse writejsonresponse data")
	errRecordNotFound = errors.New("no record found in the database")
)
