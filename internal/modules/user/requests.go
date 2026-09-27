package user

type createUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type ErrorResponse struct {
	Code    int64  `json:"status"`
	Message string `json:"message"`
}
