package models

type IdPostRequest struct {
	Id string `json:"id"`
}

type PostRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdatePostRequest struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Post struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Likes       int64  `json:"likes"`
	UserId      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Comments    int64  `json:"comments"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"update_at"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}

type DeletedPost struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Likes       int64  `json:"likes"`
	UserId      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Comments    int64  `json:"comments"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"update_at"`
	DeletedAt   string `json:"deleted_at"`
}
