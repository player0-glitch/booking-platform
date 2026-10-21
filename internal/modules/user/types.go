package user

import (
	"github.com/google/uuid"
)

type createUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type rolesRequest struct {
	//Do not allow any request to explicitly provide role id
	Id       int    `json:"-"`
	RoleName string `json:"role_name"`
}

type userRequest struct {
	Id        uuid.UUID    `json:"id" swaggertype:"string" format:"uuid"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Email     string       `json:"email"`
	Role      rolesRequest `json:"roles"`
}

func newUserRequest(u *User) userRequest {

	role := rolesRequest{
		Id:       u.RoleId,
		RoleName: u.Role.RoleName,
	}
	return userRequest{
		Id:        u.Base.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Role:      role,
	}
}
