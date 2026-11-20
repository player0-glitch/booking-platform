package user /*provision guest*/

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	json "booking-platform/internal/core/json"

	chi "github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserServiceInterface interface {
	CreateUser(ctx context.Context, name, lastName, email, password string) (*userRequest, error)
	GetAll(ctx context.Context) ([]userRequest, error)
	GetByID(ctx context.Context, id uuid.UUID) (*userRequest, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
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

// @Summary		Create A New User
// @Tags		users
// @Param request body createUserRequest true  "User Object type"
// @Produce		json
// @Success 	201
// @Failure		403
// @Router		/users	[post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	// parse the json request
	if errJSONDecoder = json.Reader(r, &req); errJSONDecoder != nil {
		json.RequestErrorHandler(w, errJSONDecoder)
		return
	}

	if err := req.Validate(); err != nil {
		json.RequestErrorHandler(w, err)
		return
	}

	user, errSvc := h.userService.CreateUser(r.Context(), req.FirstName, req.LastName, req.Email, req.Password)

	if errSvc != nil {
		json.InternalErrorHandler(w, errors.New("failed to create user"))
		return
	}
	json.Write(w, http.StatusCreated, &user)
}

// @Summary		Retrieve All Users
// @Tags		users
// @Produce		json
// @Success 	200
// @Failure		403
// @Router		/users/all	[get]
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAll(r.Context())
	if err != nil {
		json.RequestErrorHandler(w, errors.New("server failed to get all users"))
		return
	}

	if len(users) == 0 {
		json.Write(w, http.StatusNoContent, "No Users")
		return
	}
	json.Write(w, http.StatusOK, users)
}

// @Security	BearerAuth
// @Summary		Retrieve user by id
// @Tags		users
// @Param id path string true "User Id" format(uuid)
// @Produce		json
// @Success 	202	{object} userRequest
// @Failure		403
// @Failure		404
// @Router		/users/{id}	[get]
func (c *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, errParse := uuid.Parse(id)

	if errParse != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request, Check user id")
		return
	}

	user, errSvc := c.userService.GetByID(r.Context(), userID)

	if errSvc != nil {
		json.WriteError(w, http.StatusNotFound,
			fmt.Sprintf("User with id = %s not found",
				id))
		return
	}

	json.Write(w, http.StatusAccepted, user)
}

// @Security	BearerAuth
// @Summary		Delete user by id
// @Tags		users
// @Param id path string true "User Id" format(uuid)
// @Produce		json
// @Sucess 		204
// @Failure		403
// @Failure 	404
// @Router		/users/{id}	[delete]
func (c *UserHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Printf("Id before parse: %s\n", id)
	userID, errParse := uuid.Parse(id)

	fmt.Printf("Parsed Id: %s\n", userID)
	if errParse != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request, Check user id "+errParse.Error())
		return
	}
	err := c.userService.DeleteByID(r.Context(), userID)
	if err != nil {
		json.WriteError(w, http.StatusNotFound, "User Not found")
		return
	}
	json.Write(w, http.StatusAccepted, "User Deleted")
}

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
func (c *UserHandler) SoftDeleteByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, errParse := uuid.Parse(id)

	if errParse != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request, Check user id")
		return
	}
	err := c.userService.SoftDelete(r.Context(), userID)
	if err != nil {
		json.WriteError(w, http.StatusNotFound, "User Not found")
		return
	}
	json.Write(w, http.StatusAccepted, "User Deleted Softly")
}
