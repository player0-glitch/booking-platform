package auth

import "net/http"

// Private to the auth package
//
//	type authenticator interface {
//		Authenticate(email, password string) (string, error)
//	}
type storeManager interface {
	SaveToken(w http.ResponseWriter, r *http.Request, token string) error
	Clear(w http.ResponseWriter, r *http.Request) error
}
