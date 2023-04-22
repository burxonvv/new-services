package repo

type User struct {
	Id           string
	FirstName    string
	LastName     string
	UserType     string
	Email        string
	Password     string
	Posts        int64
	RefreshToken string
	CreatedAt    string
	UpdatedAt    string
}