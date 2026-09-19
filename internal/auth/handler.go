// Handlers for auth endpoints
package auth

import (
	response "booking-platform/internal/core"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	service authenticator
	store   storeManager
}

func NewAuthHandler(service authenticator, store storeManager) *AuthHandler {
	return &AuthHandler{
		service: service,
		store:   store,
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	//parse the request
	var req loginRequest
	errJsonDecoder := json.NewDecoder(r.Body).Decode(&req)

	if errJsonDecoder != nil {
		response.Error(w, http.StatusBadGateway, "invalid login request body")
		return
	}
	token, err := h.service.Authenticate(req.Email, req.Password)

	if err != nil {
		http.Error(w, "Invalid Token On Login", http.StatusUnauthorized)
		return
	}

	err = h.store.SaveToken(w, r, token)

	if err != nil {
		http.Error(w, "Failed To Persist Token", http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Logged In Successfully"})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.store.Clear(w, r); err != nil {
		http.Error(w, "Failed To Clear Session", http.StatusInternalServerError)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Logged Out Successfully"})
}
