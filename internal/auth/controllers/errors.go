package controllers

import (
	"errors"
)

var (
	errJsonDecoder = errors.New("Failed To Decode Json")
	errParse       = errors.New("Failed to parse JSON data")
)
