package models

type SearchUsers struct {
	FirstName string `json:"first_name"`
}

type IdUserRequest struct {
	Id string `json:"id"`
}

type GetAllUsersRequest struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserType  string `json:"user_type"`
	Email     string `json:"email"`
	Posts     int64  `json:"posts"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type DeletedUser struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserType  string `json:"user_type"`
	Email     string `json:"email"`
	Posts     int64  `json:"posts"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}
