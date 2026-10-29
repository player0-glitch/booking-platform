package user /*provision guest*/

import (
	json "booking-platform/internal/core/json"
	"errors"
	"fmt"
	"net/http"

	"context"

	chi "github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserServiceInterface interface {
	CreateUser(ctx context.Context, name, lastName, email, password string) (*userRequest, error)
	GetAll(ctx context.Context) ([]userRequest, error)
	GetById(ctx context.Context, id uuid.UUID) (*userRequest, error)
	DeleteById(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type UserHandler struct {
	userService UserServiceInterface
}

func NewUserController(userSvc UserServiceInterface) *UserHandler {
	return &UserHandler{
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
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	//parse the json request
	errJsonDecoder = json.Reader(r, &req)

	if errJsonDecoder != nil {
		json.RequestErrorHandler(w, errJsonDecoder)
		return
	}
	user, errSvc := h.userService.CreateUser(r.Context(), req.FirstName, req.LastName, req.Email, req.Password)

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
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAll(r.Context())
	if err != nil {
		json.RequestErrorHandler(w, errors.New("Server Failed To Get All Users"))
		return
	}

	if len(users) == 0 {
		json.Write(w, http.StatusNoContent, "No Users")
		return
	}
	json.Write(w, http.StatusOK, users)
}

// User godoc
// @Security	BearerAuth
// @Summary		Retrieve user by id
// @Tags		users
// @Param id path string true "User Id" format(uuid)
// @Produce		json
// @Success 	202	{object} userRequest
// @Failure		403
// @Failure		404
// @Router		/users/{id}	[get]
func (c *UserHandler) GetById(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	userId, errParse := uuid.Parse(id)

	if errParse != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request, Check user id")
		return
	}

	user, errSvc := c.userService.GetById(r.Context(), userId)

	if errSvc != nil {
		json.WriteError(w, http.StatusNotFound,
			fmt.Sprintf("User with id = %s not found",
				id))
		return
	}

	json.Write(w, http.StatusAccepted, user)
}

// User godoc
// @Security	BearerAuth
// @Summary		Delete user by id
// @Tags		users
// @Param id path string true "User Id" format(uuid)
// @Produce		json
// @Sucess 		204
// @Failure		403
// @Failure 	404
// @Router		/users/{id}	[delete]
func (c *UserHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Printf("Id before parse: %s\n", id)
	userId, errParse := uuid.Parse(id)

	fmt.Printf("Parsed Id: %s\n", userId)
	if errParse != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request, Check user id "+errParse.Error())
		return
	}
	err := c.userService.DeleteById(r.Context(), userId)
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
// @Param id path string true "User Id" format (uuid)
// @Produce		json
// @Success 	200	{object} userRequest
// @Success 	204
// @Failure		403
// @Failure 	404
// @Router		/users/{id}/soft	[delete]
func (c *UserHandler) SoftDeleteById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userId, errParse := uuid.Parse(id)

	if errParse != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request, Check user id")
		return
	}
	err := c.userService.SoftDelete(r.Context(), userId)
	if err != nil {
		json.WriteError(w, http.StatusNotFound, "User Not found")
		return
	}
	json.Write(w, http.StatusAccepted, "User Deleted Softly")
}

// func (h *UserHandler) GetUsersByRole(w http.ResponseWriter, r *http.Request) {
// 	rolename := chi.URLParam(r, "rolename")

// 	users, err := h.userService.GetUsersByRole(r.Context(), rolename)
// 	if err != nil {
// 		json.WriteError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	if len(users) == 0 {
// 		json.Write(w, http.StatusNoContent, "No Users For Give Role")
// 		return
// 	}

// 	json.Write(w, http.StatusOK, users)
// }
