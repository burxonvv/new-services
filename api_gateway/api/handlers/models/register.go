package models

type GetProfileByJWTRequest struct {
	Token string `header:"Authorization"`
}

type LoginModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponseModel struct {
	Id           string
	AccessToken  string
	RefreshToken string
}

type RegisterModel struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Code      string `json:"code"`
}

type UserRegister struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginResponseModel struct {
	Id           string
	FirstName    string
	LastName     string
	UserType     string
	Email        string
	Password     string
	AccessToken  string
	RefreshToken string
}
