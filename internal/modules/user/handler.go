package user

import (
	response "booking-platform/internal/core"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	chi "github.com/go-chi/chi/v5"
)

type UserController struct {
	userService *UserService
}

func NewUserController(userSvc *UserService) *UserController {
	return &UserController{
		userService: userSvc,
	}
}

type createUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (c *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	//parse the json request
	errJsonDecoder = json.NewDecoder(r.Body).Decode(&req)

	if errJsonDecoder != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	user, errSvc := c.userService.CreateUser(r.Context(), req.FirstName, req.LastName, req.Email, req.Password)

	if errSvc != nil {
		response.Error(w, http.StatusInternalServerError, "Failed To Create User")
		return
	}
	response.JSON(w, http.StatusCreated, &user)
}

func (c *UserController) FindAll(w http.ResponseWriter, r *http.Request) {
	users, err := c.userService.FindAll(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Server Failed To Get All Users")
		return
	}

	if len(users) == 0 {
		response.JSON(w, http.StatusNotFound, "No Users")
		return
	}
	response.JSON(w, http.StatusFound, users)
}
func (c *UserController) GetById(w http.ResponseWriter, r *http.Request) {

	id, errParse := strconv.Atoi(chi.URLParam(r, "id"))
	if errParse != nil {
		response.Error(w, http.StatusBadRequest, "invalid request, Check user id")
		return
	}
	user, errSvc := c.userService.GetById(r.Context(), id)
	if errSvc != nil {
		response.Error(w, http.StatusNotFound, fmt.Sprintf("User with id=%d not found", id))
		return
	}

	response.JSON(w, http.StatusFound, user)
}

func (c *UserController) DeleteById(w http.ResponseWriter, r *http.Request) {
	id, errParse := strconv.Atoi(chi.URLParam(r, "id"))
	if errParse != nil {
		response.Error(w, http.StatusBadRequest, "invalid request, Check user id")
		return
	}
	err := c.userService.DeleteById(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "User Not found")
		return
	}
	response.JSON(w, http.StatusAccepted, "User Deleted")
}

func (c *UserController) SoftDeleteById(w http.ResponseWriter, r *http.Request) {
	id, errParse := strconv.Atoi(chi.URLParam(r, "id"))
	if errParse != nil {
		response.Error(w, http.StatusBadRequest, "invalid request, Check user id")
		return
	}
	err := c.userService.SoftDelete(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "User Not found")
		return
	}
	response.JSON(w, http.StatusAccepted, "User Deleted Softly")
}
