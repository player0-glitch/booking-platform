package auth

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type loginResponse struct {
	Token string `json:"token"`
}
type ErrorResponse struct {
	Code    int64  `json:"status"`
	Message string `json:"message"`
}
