package repo

type CommentStorageI interface {
	WriteComment(Comment) (Comment, error)
	GetComments(id string) ([]Comment, error)
	DeleteComment(id string) (Comment, error)
}
