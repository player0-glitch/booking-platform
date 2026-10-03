package user

type createUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type rolesRequest struct {
	Id       uint   `json:"id"`
	RoleName string `json:"role_name" example:"Admin"`
}

type userRequest struct {
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Email     string         `json:"email"`
	Roles     []rolesRequest `json:"roles"`
}

func newUserRequest(u *User) *userRequest {
	roles := make([]rolesRequest, len(u.Roles))

	for i, role := range u.Roles {
		roles[i] = rolesRequest{
			Id:       role.Model.ID,
			RoleName: role.RoleName,
		}
	}
	return &userRequest{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Roles:     roles,
	}
}
