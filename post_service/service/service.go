package service

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"

	c "github.com/burxondv/new-services/post-service/genproto/comment"
	p "github.com/burxondv/new-services/post-service/genproto/post"
	u "github.com/burxondv/new-services/post-service/genproto/user"
	"github.com/burxondv/new-services/post-service/pkg/logger"
	grpcclient "github.com/burxondv/new-services/post-service/service/grpc_client"
	"github.com/burxondv/new-services/post-service/storage"
	"github.com/burxondv/new-services/post-service/storage/repo"
)

type PostService struct {
	storage storage.IStorage
	Logger  logger.Logger
	Client  grpcclient.Clients
}

func NewPostService(db *sqlx.DB, log logger.Logger, client grpcclient.Clients) *PostService {
	return &PostService{
		storage: storage.NewStoragePg(db),
		Logger:  log,
		Client:  client,
	}
}

func (s *PostService) CreatePost(ctx context.Context, req *p.PostRequest) (*p.PostResponse, error) {
	postResp := p.PostResponse{}
	res, err := s.storage.Post().CreatePost(repo.Post{
		Id:          req.Id,
		Title:       req.Title,
		Description: req.Description,
		UserId:      req.UserId,
	})
	if err != nil {
		log.Println("failed to create post in service: ", err)
		return &p.PostResponse{}, err
	}

	user, err := s.Client.User().GetUserForClient(ctx, &u.Request{Str: res.UserId})
	if err != nil {
		log.Println("failed to get user for create post: ", err)
		return &p.PostResponse{}, err
	}
	postResp.Id = res.Id
	postResp.Title = res.Title
	postResp.Description = res.Description
	postResp.UserId = res.UserId
	postResp.Likes = res.Likes
	postResp.UserName = user.FirstName + " " + user.LastName
	postResp.CreatedAt = res.CreatedAt
	postResp.UpdatedAt = res.UpdatedAt

	return &postResp, nil
}

func (s *PostService) GetPostById(ctx context.Context, req *p.Request) (*p.PostResponse, error) {
	postResp := p.PostResponse{}
	res, err := s.storage.Post().GetPostById(req.Str)
	if err != nil {
		log.Println("failed to get post by id: ", err)
		return &p.PostResponse{}, err
	}

	postUser, err := s.Client.User().GetUserForClient(ctx, &u.Request{Str: res.UserId})
	if err != nil {
		log.Println("failed to get user for get post by id: ", err)
		return &p.PostResponse{}, err
	}

	comments, err := s.Client.Comment().GetCommentsForPost(ctx, &c.Request{Str: res.Id})
	if err != nil {
		log.Println("failed to get comments for get post by id: ", err)
		return &p.PostResponse{}, err
	}

	postResp.Id = res.Id
	postResp.Title = res.Title
	postResp.Description = res.Description
	postResp.UserId = res.UserId
	postResp.Likes = res.Likes
	postResp.UserName = postUser.FirstName + " " + postUser.LastName
	postResp.Comments = int64(len(comments.Comments))
	postResp.CreatedAt = res.CreatedAt
	postResp.UpdatedAt = res.UpdatedAt

	return &postResp, nil
}

func (s *PostService) GetPostByUserId(ctx context.Context, req *p.Request) (*p.PostsResponse, error) {
	postsResp := p.PostsResponse{}
	res, err := s.storage.Post().GetPostByUserId(req.Str)
	if err != nil {
		log.Println("failed to get post by user id: ", err)
		return &p.PostsResponse{}, err
	}

	postUser, err := s.Client.User().GetUserForClient(ctx, &u.Request{Str: req.Str})
	if err != nil {
		log.Println("failed to get user for get post by user id: ", err)
		return &p.PostsResponse{}, err
	}

	for _, ps := range res {
		post := p.PostResponse{}
		post.Id = ps.Id
		post.Title = ps.Title
		post.Description = ps.Description
		post.Likes = ps.Likes
		post.UserId = postUser.Id
		post.CreatedAt = ps.CreatedAt
		post.UpdatedAt = ps.UpdatedAt
		postsResp.Posts = append(postsResp.Posts, &post)
	}

	for _, p := range postsResp.Posts {
		p.UserName = postUser.FirstName + " " + postUser.LastName
	}

	for _, post := range postsResp.Posts {
		comments, err := s.Client.Comment().GetCommentsForPost(ctx, &c.Request{Str: post.Id})
		if err != nil {
			log.Println("failed to get comments for get post by user id: ", err)
			return &p.PostsResponse{}, err
		}

		post.Comments = int64(len(comments.Comments))
	}

	return &postsResp, nil
}

func (s *PostService) GetPostForUser(ctx context.Context, req *p.Request) (*p.PostsResponse, error) {
	postsResp := p.PostsResponse{}
	res, err := s.storage.Post().GetPostForUser(req.Str)
	if err != nil {
		log.Println("failed to get post for user: ", err)
		return &p.PostsResponse{}, err
	}

	for _, ps := range res {
		post := p.PostResponse{}
		post.Id = ps.Id
		post.Title = ps.Title
		post.Description = ps.Description
		post.Likes = ps.Likes
		post.UserId = ps.UserId
		post.CreatedAt = ps.CreatedAt
		post.UpdatedAt = ps.UpdatedAt
		postsResp.Posts = append(postsResp.Posts, &post)
	}

	return &postsResp, nil
}

func (s *PostService) GetPostForComment(ctx context.Context, req *p.Request) (*p.PostResponse, error) {
	postResp := p.PostResponse{}
	res, err := s.storage.Post().GetPostForComment(req.Str)
	if err != nil {
		log.Println("failed to get post for comment: ", err)
		return &p.PostResponse{}, err
	}

	postResp.Id = res.Id
	postResp.Title = res.Title
	postResp.Description = res.Description
	postResp.Likes = res.Likes
	postResp.UserId = res.UserId
	postResp.CreatedAt = res.CreatedAt
	postResp.UpdatedAt = res.UpdatedAt

	return &postResp, nil
}

func (s *PostService) SearchPosts(ctx context.Context, req *p.Request) (*p.PostsResponse, error) {
	postsResp := p.PostsResponse{}
	res, err := s.storage.Post().SearchPosts(req.Str)
	if err != nil {
		log.Println("failed to get posts by search title: ", err)
		return &p.PostsResponse{}, err
	}

	for _, ps := range res {
		post := p.PostResponse{}
		post.Id = ps.Id
		post.Title = ps.Title
		post.Description = ps.Description
		post.Likes = ps.Likes
		post.UserId = ps.UserId
		post.CreatedAt = ps.CreatedAt
		post.UpdatedAt = ps.UpdatedAt
		postsResp.Posts = append(postsResp.Posts, &post)
	}

	for _, post := range postsResp.Posts {
		postUser, err := s.Client.User().GetUserForClient(ctx, &u.Request{Str: post.UserId})
		if err != nil {
			log.Println("failed to get user for get posts by search title: ", err)
			return &p.PostsResponse{}, err
		}

		post.UserName = postUser.FirstName + " " + postUser.LastName

		comments, err := s.Client.Comment().GetCommentsForPost(ctx, &c.Request{Str: post.Id})
		if err != nil {
			log.Println("failed to get comments for get posts by search title: ", err)
			return &p.PostsResponse{}, err
		}

		post.Comments = int64(len(comments.Comments))
	}

	return &postsResp, nil
}

func (s *PostService) LikePost(ctx context.Context, req *p.LikeRequest) (*p.PostResponse, error) {
	postResp := p.PostResponse{}
	res, err := s.storage.Post().LikePost(req.PostId, req.IsLiked)
	if err != nil {
		log.Println("failed to like post: ", err)
		return &p.PostResponse{}, err
	}
	if !req.IsLiked {
		res.Likes -= 1
	}

	postUser, err := s.Client.User().GetUserForClient(ctx, &u.Request{Str: res.UserId})
	if err != nil {
		log.Println("failed to get user for like post: ", err)
		return &p.PostResponse{}, err
	}

	postResp.Id = res.Id
	postResp.Title = res.Title
	postResp.Description = res.Description
	postResp.Likes = res.Likes
	postResp.UserId = res.UserId
	postResp.UserName = postUser.FirstName + " " + postUser.LastName
	postResp.CreatedAt = res.CreatedAt
	postResp.UpdatedAt = res.UpdatedAt

	comments, err := s.Client.Comment().GetCommentsForPost(ctx, &c.Request{Str: postResp.Id})
	if err != nil {
		log.Println("failed to get comments for like post: ", err)
		return &p.PostResponse{}, err
	}

	postResp.Comments = int64(len(comments.Comments))

	return &postResp, nil
}

func (s *PostService) UpdatePost(ctx context.Context, req *p.UpdatePostRequest) (*p.PostResponse, error) {
	res, err := s.storage.Post().UpdatePost(repo.Post{
		Title:       req.Title,
		Description: req.Description,
		Id:          req.Id,
	})
	if err != nil {
		log.Println("failed to update post: ", err)
		return &p.PostResponse{}, err
	}

	user, err := s.Client.User().GetUserForClient(ctx, &u.Request{Str: res.UserId})
	if err != nil {
		log.Println("failed to get user for update post: ", err)
		return &p.PostResponse{}, err
	}

	return &p.PostResponse{
		Id:          res.Id,
		Title:       res.Title,
		Description: res.Description,
		Likes:       res.Likes,
		UserId:      res.UserId,
		UserName:    user.FirstName + " " + user.LastName,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *p.Request) (*p.PostResponse, error) {
	postResp := p.PostResponse{}
	res, err := s.storage.Post().DeletePost(req.Str)
	if err != nil {
		log.Println("failed to delete post: ", err)
		return &p.PostResponse{}, err
	}

	postUser, err := s.Client.User().GetUserForClient(ctx, &u.Request{Str: res.UserId})
	if err != nil {
		log.Println("failed to get user for delete post: ", err)
		return &p.PostResponse{}, err
	}

	postResp.Id = res.Id
	postResp.Title = res.Title
	postResp.Description = res.Description
	postResp.Likes = res.Likes
	postResp.UserId = res.UserId
	postResp.UserName = postUser.FirstName + " " + postUser.LastName
	postResp.CreatedAt = res.CreatedAt
	postResp.UpdatedAt = res.UpdatedAt

	comments, err := s.Client.Comment().GetCommentsForPost(ctx, &c.Request{Str: res.Id})
	if err != nil {
		log.Println("failed to get comments for delete post: ", err)
		return &p.PostResponse{}, err
	}

	postResp.Comments = int64(len(comments.Comments))

	return &postResp, err
}
