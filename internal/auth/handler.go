// Handlers for auth endpoints
package auth

import (
	json "booking-platform/internal/core/json"
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

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type loginResponse struct {
	Token string `json:"jwt_token"`
}

// godoc
// Login
// @Summary		Login and return jwt
// @Description	Logs in a user, adds them to a 2hr session and return jwt
// @Tags		auth
// @Param 	request body loginRequest true "Credentials Used For Logging In"
// @Produce		json
// @Success		200 {object} tokenResponse "Returns token object"
// @Failure		400  "Bad Request"
// @Failure		401  "Unauthorised"
// @Router		/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	var res tokenResponse
	errJsonDecoder := json.Reader(r, &req)

	if errJsonDecoder != nil {
		json.RequestErrorHandler(w, errJsonDecoder)
		return
	}

	token, err := h.service.Login(r.Context(), req.Email, req.Password)

	if err != nil {
		http.Error(w, "Invalid Token On Login", http.StatusUnauthorized)
		return
	}
	err = h.service.store.SaveToken(w, r, token)

	if err != nil {
		http.Error(w, "Failed To Persist Token", http.StatusInternalServerError)
		return
	}

	res.Token = token
	json.Write(w, http.StatusOK, res)
}

// godoc
// Logout
// @Summary		Logout and invalidate session and token
// @Tags		auth
// @Produce		json
// @Success		200 "Successfully Logged out"
// @Failure 	401 "Unauthenticated"
// @Failure		500
// @Router 		/auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.service.store.Clear(w, r); err != nil {
		http.Error(w, "Failed To Clear Session", http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusOK, map[string]string{"message": "Logged Out Successfully"})
}

type tokenResponse struct {
	Token string `json:"token"`
}

type emailRequest struct {
	Email string `json:"email"`
}

// godoc
// Request Password Reset
// @Summary Request password reset token
// @Description		Request for any role but 'Guest' to request password reset token
// @Tags		auth
// @Param	request body emailRequest true "Password reset credentials"
// @Produce		json
// @Success		200 {object} tokenResponse "Returns Raw Reset Token"
// @Failure		400 "Bad Request"
// @Failure 	403 "Forbidden"
// @Router 		/auth/password/reset [post]
func (h *AuthHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req emailRequest
	err := json.Reader(r, req)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.service.RequestPasswordReset(r.Context(), req.Email)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	json.Write(w, http.StatusAccepted, token)
}

type resetTokenResponse struct {
	Token string `json:"reset_token"`
}

type passwordResetRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

// godoc
// @Description	Reset password for any role but 'Guest'
// @Summary 	Reset Recognised User Password
// @Tags		auth
// @Param	request body passwordResetRequest true "Reset Credentials"
// @Produce		json
// @Success		200 "Successfully Updated User Password"
// @Failure		403 "Forbidden"
// @Router 		/auth/password [post]
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req passwordResetRequest

	errJsonDecode := json.Reader(r, req)
	if errJsonDecode != nil {
		// Sends a default status 400
		json.RequestErrorHandler(w, errJsonDecode)
		return
	}

	err := h.service.ResetPassword(r.Context(), req.Email, req.Password, req.Token)
	if err != nil {
		json.InternalErrorHandler(w, err)
		return
	}

	json.Write(w, http.StatusOK,
		"Successfully Updated User Password")
}
