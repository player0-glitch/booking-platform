package user

import (
	json "booking-platform/internal/json"
	"fmt"
	"net/http"
	"strconv"

	"errors"

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

// User godoc
// @Summary		Create A New User
// @Tags		users
// @Param request body createUserRequest true  "User Object type"
// @Produce		json
// @Success 	201
// @Failure		403
// @Router		/users	[post]
func (c *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	//parse the json request
	errJsonDecoder = json.Reader(r, &req)

	if errJsonDecoder != nil {
		json.RequestErrorHandler(w, errJsonDecoder)
		return
	}
	user, errSvc := c.userService.CreateUser(r.Context(), req.FirstName, req.LastName, req.Email, req.Password)

	if errSvc != nil {
		json.InternalErrorHandler(w, errors.New("Failed To Create User"))
		return
	}
	json.Write(w, http.StatusCreated, &user)
}

// User godoc
// @Summary		Retrieve All Users
// @Tags		users
// @Produce		json
// @Success 	200
// @Failure		403
// @Router		/users/all	[get]
func (c *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := c.userService.FindAll(r.Context())
	if err != nil {
		json.RequestErrorHandler(w, errors.New("Server Failed To Get All Users"))
		return
	}

	if len(users) == 0 {
		json.Write(w, http.StatusNotFound, "No Users")
		return
	}
	json.Write(w, http.StatusOK, users)
}

// User godoc
// @Security	BearerAuth
// @Summary		Retrieve user by id
// @Tags		users
// @Param id path int true "User Id"
// @Produce		json
// @Success 	202	{object} userRequest
// @Failure		403
// @Failure		404
// @Router		/users/{id}	[get]
func (c *UserController) GetById(w http.ResponseWriter, r *http.Request) {

	id, errParse := strconv.Atoi(chi.URLParam(r, "id"))
	if errParse != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request, Check user id")
		return
	}
	user, errSvc := c.userService.GetById(r.Context(), uint(id))
	if errSvc != nil {
		json.WriteError(w, http.StatusNotFound, fmt.Sprintf("User with id=%d not found", id))
		return
	}

	json.Write(w, http.StatusAccepted, user)
}

// User godoc
// @Security	BearerAuth
// @Summary		Delete user by id
// @Tags		users
// @Param request path int true "User Id"
// @Produce		json
// @Success 	200	{object} userRequest
// @Failure		403
// @Failure 	404
// @Router		/users/{id}	[delete]
func (c *UserController) DeleteById(w http.ResponseWriter, r *http.Request) {
	id, errParse := strconv.Atoi(chi.URLParam(r, "id"))
	if errParse != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request, Check user id")
		return
	}
	err := c.userService.DeleteById(r.Context(), uint(id))
	if err != nil {
		json.WriteError(w, http.StatusNotFound, "User Not found")
		return
	}
	json.Write(w, http.StatusAccepted, "User Deleted")
}

// User godoc
// @Security	BearerAuth
// @Summary		Solf delete user by id
// @Tags		users
// @Param request path int true "User Id"
// @Produce		json
// @Success 	200	{object} userRequest
// @Failure		403
// @Failure 	404
// @Router		/users/{id}/soft	[delete]
func (c *UserController) SoftDeleteById(w http.ResponseWriter, r *http.Request) {
	id, errParse := strconv.Atoi(chi.URLParam(r, "id"))
	if errParse != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request, Check user id")
		return
	}
	err := c.userService.SoftDelete(r.Context(), uint(id))
	if err != nil {
		json.WriteError(w, http.StatusNotFound, "User Not found")
		return
	}
	json.Write(w, http.StatusAccepted, "User Deleted Softly")
}
