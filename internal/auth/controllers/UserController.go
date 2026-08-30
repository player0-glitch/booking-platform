package controllers

import (
	"booking-platform/internal/auth/services"
	"booking-platform/internal/core"
	// "booking-platform/internal/core/response"
	"encoding/json"
	"errors"
	"net/http"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userSvc *services.UserService) *UserController {
	return &UserController{
		userService: userSvc,
	}
}

type createUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"passsword"`
}

var (
	errJsonDecoder = errors.New("Failed To Decode Json")
	errJsonEncoder = errors.New("Failed To Encode Json")
	errUserCreate  = errors.New("Failed To Store User")
)

// End Point
func (c *UserController) Create(w http.ResponseWriter, r *http.Request) {
	//create a request type
	var userRequest createUserRequest
	//parse the json request
	errJsonDecoder := json.NewDecoder(r.Body).Decode(&userRequest)

	if errJsonDecoder != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	user, errSvc := c.userService.CreateUser(userRequest.FirstName,
		userRequest.LastName, userRequest.Email, userRequest.Password)

	if errSvc != nil {
		response.Error(w, http.StatusInternalServerError, "Failed To Create User")
		return
	}
	response.JSON(w, http.StatusCreated, &user)
}
