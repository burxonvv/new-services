package models

type IdCommentRequest struct {
	Id string `json:"id"`
}

type CommentRequest struct {
	PostId string `json:"post_id"`
	Text   string `json:"text"`
}

type Comment struct {
	Id           string `json:"id"`
	PostId       string `json:"post_id"`
	PostTitle    string `json:"post_title"`
	UserId       string `json:"user_id"`
	UserName     string `json:"user_name"`
	UserType     string `json:"user_type"`
	PostUserName string `json:"post_user_name"`
	Text         string `json:"text"`
	CreatedAt    string `json:"created_at"`
}

type Comments struct {
	Comments []Comment `json:"comments"`
}

type DeletedComment struct {
	Id           string `json:"id"`
	PostId       string `json:"post_id"`
	PostTitle    string `json:"post_title"`
	UserId       string `json:"user_id"`
	UserName     string `json:"user_name"`
	UserType     string `json:"user_type"`
	PostUserName string `json:"post_user_name"`
	Text         string `json:"text"`
	CreatedAt    string `json:"created_at"`
	DeletedAt    string `json:"deleted_at"`
}
