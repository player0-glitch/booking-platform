// Handlers for authentication endpoints
package authentication

import (
	"net/http"

	json "booking-platform/internal/core/json"
	"booking-platform/internal/sessions"

	gorillaSessions "github.com/gorilla/sessions"
)

type handlerSessionInterface interface {
	Rotate(w http.ResponseWriter, r *http.Request) (*gorillaSessions.Session, error)
	Clear(w http.ResponseWriter, r *http.Request) error
}

type Handler struct {
	service *Service
	sess    handlerSessionInterface
}

func NewHandler(service *Service, sessions handlerSessionInterface) *Handler {
	return &Handler{
		service: service,
		sess:    sessions,
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
// @Tags		authentication
// @Param 	request body loginRequest true "Credentials Used For Logging In"
// @Produce		json
// @Success		200 {object} tokenResponse "Returns token object"
// @Failure		400  "Bad Request"
// @Failure		401  "Unauthorised"
// @Router		/authentication/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
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
	// TODO: this should eventually be in a database transaction. This needs to be atomic
	session, err := h.sess.Rotate(w, r)
	if err != nil {
		http.Error(w, "session failed to be created or rotated", http.StatusInternalServerError)
		return
	}
	session.Values[sessions.JWTToken] = token
	//deleting a session
	delete(session.Values, sessions.SessionType)

	if err != nil {
		http.Error(w, "failed to save token", http.StatusInternalServerError)
		return
	}

	res.Token = token
	res.Session = session
	json.Write(w, http.StatusNoContent, res)
}

// godoc
// Logout
// @Summary		Logout and invalidate session and token
// @Tags		authentication
// @Produce		json
// @Success		200 "Successfully Logged out"
// @Failure 	401 "Unauthenticated"
// @Failure		500
// @Router 		/authentication/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.sess.Clear(w, r); err != nil {
		http.Error(w, "Failed To Clear Session", http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusOK, map[string]string{"message": "Logged Out Successfully"})
}

type tokenResponse struct {
	Token   string `json:"token"`
	Session any    `json:"-"`
}

type emailRequest struct {
	Email string `json:"email"`
}

// godoc
// Request Password Reset
// @Summary Request password reset token
// @Description		Request for any role but 'Guest' to request password reset token
// @Tags		authentication
// @Param	request body emailRequest true "Password reset credentials"
// @Produce		json
// @Success		200 {object} tokenResponse "Returns Raw Reset Token"
// @Failure		400 "Bad Request"
// @Failure 	403 "Forbidden"
// @Router 		/authentication/password/reset [post]
func (h *Handler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
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
// @Tags		authentication
// @Param	request body passwordResetRequest true "Reset Credentials"
// @Produce		json
// @Success		200 "Successfully Updated User Password"
// @Failure		403 "Forbidden"
// @Router 		/authentication/password [post]
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
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
