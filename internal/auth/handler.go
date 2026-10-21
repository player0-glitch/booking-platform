// Handlers for auth endpoints
package auth

import (
	jsonApi "booking-platform/internal/json"
	"net/http"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

/*
This is how to define parameters for an api endpoint
@Param <name of the parameter> <location in the request> <type> <required?>
{
  "status": 400,
  "message": "json: Unmarshal(non-pointer auth.loginRequest)"
}

<description for swagger ui>
*/

// Login godoc
//
// @Summary			Login and return jwt
// @Description		Logs in a user, adds them to a 2hr session and return jwt
// @Tags			auth
// @Param request body loginRequest true "Credentials Used For Logging In"
// @Produce			json
// @Success			200 {object} loginRequest
// Failure			400 {object} ErrorResponse
// Failure			401 {object} ErrorResponse
// @Router			/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	errJsonDecoder := jsonApi.Reader(r, &req)

	if errJsonDecoder != nil {
		jsonApi.RequestErrorHandler(w, errJsonDecoder)
		return
	}

	token, err := h.service.Authenticate(r.Context(), req.Email, req.Password)

	if err != nil {
		http.Error(w, "Invalid Token On Login", http.StatusUnauthorized)
		return
	}

	err = h.service.store.SaveToken(w, r, token)

	if err != nil {
		http.Error(w, "Failed To Persist Token", http.StatusInternalServerError)
		return
	}

	jsonApi.Write(w, http.StatusOK,
		map[string]string{"token": token})
}

// Logout godoc
// @Summary		Logout and invalidate session and token
// @Tags		auth
// @Produce		json
// @Success 	200
// @Failure 	500
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.service.store.Clear(w, r); err != nil {
		http.Error(w, "Failed To Clear Session", http.StatusInternalServerError)
		return
	}
	jsonApi.Write(w, http.StatusOK, map[string]string{"message": "Logged Out Successfully"})
}
