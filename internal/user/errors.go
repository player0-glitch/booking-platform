package user

import (
	"errors"
)

// Handler/Controller errors
var (
	errJsonDecoder    = errors.New("Failed To Decode Json")
	errParse          = errors.New("Failed to parse JSON data")
	errRecordNotFound = errors.New("No Record Found In The Database")
)
