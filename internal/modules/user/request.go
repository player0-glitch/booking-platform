package user

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"

	"booking-platform/internal/core"

	"github.com/google/uuid"
)

// Reponse types for User Handler
type deletedUserResponse struct {
	ID uuid.UUID `json:"id" swaggertype:"string" format:"uuid"`
}

// Request types for User Handler
type createUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (r createUserRequest) Validate() error {
	var err error
	if r.FirstName == "" {
		err = errors.Join(err, errors.New("first name expected"))
	}
	if r.Email == "" {
		err = errors.Join(err, errors.New("email is expected"))
	}
	if r.LastName == "" {
		err = errors.Join(err, errors.New("last name expected"))
	}

	var hasUpper, hasLower, hasSpecial, hasSpace, hasDigit bool
	if utf8.RuneCountInString(r.Password) < 8 {
		err = errors.Join(err, errors.New("password must be atleast 8 characters long"))
	}
	for _, c := range r.Password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	if !hasUpper {
		err = errors.New("password must have atleast 1 uppercase letter")
	}
	if !hasLower {
		err = errors.New("password must have atleast 1 lowercase letter")
	}
	if !hasDigit {
		err = errors.New("password must have atleast 1 digit")
	}
	if hasSpace {
		err = errors.New("password cannot contain space")
	}
	if !hasSpecial {
		err = errors.New("password must have atleast 1 special character")
	}
	if strings.Contains(strings.ToLower(r.Password),
		strings.ToLower(r.Email)) ||
		strings.Contains(strings.ToLower(r.Password),
			strings.ToLower(r.FirstName)) ||
		strings.Contains(strings.ToLower(r.Password),
			strings.ToLower(r.LastName)) {
		err = errors.Join(err, errors.New("personal infomation cannot be used in password"))
	}
	return err
}

type rolesRequest struct {
	// Do not allow any request to explicitly provide role id
	ID      int    `json:"-"`
	RoleKey string `json:"role_key"`
}

func (r rolesRequest) Validate() error {
	var err error
	role := core.IDToRole[core.RoleID(r.ID)]
	if !role.IsValid() {
		err = errors.Join(err, errors.New("role is invalid"))
	}
	return err
}

type userRequest struct {
	ID        uuid.UUID    `json:"id" swaggertype:"string" format:"uuid"`
	FirstName *string      `json:"first_name"`
	LastName  *string      `json:"last_name"`
	Email     *string      `json:"email"`
	Role      rolesRequest `json:"roles"`
}

func newUserRequest(u *User) userRequest {
	return userRequest{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
}
