// Handlers for auth endpoints
package auth

import (
	response "booking-platform/internal/core"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	service *AuthService
	store   *StoreManager
}

func NewAuthHandler(service *AuthService, store *StoreManager) *AuthHandler {
	return &AuthHandler{
		service: service,
		store:   store,
	}
}

/*
This is how to define parameters for an api endpoint
@Param <name of the parameter> <location in the request> <type> <required?>
<description for swagger ui>
*/

// Login godoc
//
// @Summary			Login and return jwt
// @Description		Logs in a user, adds them to a 2hr session and return jwt
// @Tags			auth
// @Param request body loginRequest true "Credentials Used For Logging In"
// @Produce			json
// @Success			200 {object} loginResponse
// Failure			400 {object} ErrorResponse
// Failure			401 {object} ErrorResponse
// @Router			/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	errJsonDecoder := json.NewDecoder(r.Body).Decode(&req)

	if errJsonDecoder != nil {
		response.Error(w, http.StatusBadGateway, "invalid login request body")
		return
	}
	token, err := h.service.Authenticate(r.Context(), req.Email, req.Password)

	if err != nil {
		http.Error(w, "Invalid Token On Login", http.StatusUnauthorized)
		return
	}

	err = h.store.SaveToken(w, r, token)

	if err != nil {
		http.Error(w, "Failed To Persist Token", http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK,
		map[string]string{"token": token})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.store.Clear(w, r); err != nil {
		http.Error(w, "Failed To Clear Session", http.StatusInternalServerError)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Logged Out Successfully"})
}
