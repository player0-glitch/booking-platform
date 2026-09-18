package auth

import "net/http"

type Session interface {
	Get(r *http.Request, key string) (any, error)
	Set(w http.ResponseWriter, r *http.Request, key string, val any) error
	Delete(w http.ResponseWriter, r *http.Request, key string) error
	Clear(w http.ResponseWriter, r *http.Request) error
}
