package repo

import (
	u "github.com/burxondv/new-services/user-service/genproto/user"
)

type UserStoreI interface {
	// methods...
	CreateUser(User) (User, error)
	GetUserById(string) (User, error)
	GetUserByEmail(string) (User, error)
	GetAllUsers(page, limit int64) ([]User, error)
	SearchUsers(string) ([]User, error)
	UpdateUser(User) (User, error)
	DeleteUser(string) (User, error)

	// check...
	CheckField(*u.CheckFieldRequest) (*u.CheckFieldResponse, error)

	// Register...
	UpdateUserTokens(*u.UpdateUserTokensRequest) (*u.UserResponse, error)

	// for Client...
	GetUserForClient(string) (User, error)

	// Casbin...
	ChangeRoleUser(*u.ChangeRoleRequest) (*u.UserResponse, error)
	GetSameRoleUsers(string) (*u.UsersResponse, error)
}
