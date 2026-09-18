package main

import "errors"

var (
	errInitModule = errors.New("Failed To Plug in Module")
	errAppStart   = errors.New("Failed To Start Application With Initialised Modules")
)
