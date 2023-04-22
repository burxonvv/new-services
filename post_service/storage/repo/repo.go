package repo

type PostStorageI interface {
	CreatePost(Post) (Post, error)
	GetPostById(string) (Post, error)
	GetPostByUserId(string) ([]Post, error)
	SearchPosts(string) ([]Post, error)
	LikePost(post_id string, is_liked bool) (Post, error)
	UpdatePost(Post) (Post, error)
	DeletePost(string) (Post, error)

	// for Clients...
	GetPostForUser(string) ([]Post, error)
	GetPostForComment(string) (Post, error)
}