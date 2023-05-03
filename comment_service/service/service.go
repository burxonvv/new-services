package service

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"

	c "github.com/burxondv/new-services/comment-service/genproto/comment"
	p "github.com/burxondv/new-services/comment-service/genproto/post"
	u "github.com/burxondv/new-services/comment-service/genproto/user"
	"github.com/burxondv/new-services/comment-service/pkg/logger"
	grpcclient "github.com/burxondv/new-services/comment-service/service/grpc_client"
	"github.com/burxondv/new-services/comment-service/storage"
	"github.com/burxondv/new-services/comment-service/storage/repo"
)

type CommentService struct {
	storage storage.IStorage
	Logger  logger.Logger
	Client  grpcclient.Clients
}

func NewCommentService(db *sqlx.DB, log logger.Logger, client grpcclient.Clients) *CommentService {
	return &CommentService{
		storage: storage.NewStoragePg(db),
		Logger:  log,
		Client:  client,
	}
}

func (s *CommentService) WriteComment(ctx context.Context, req *c.CommentRequest) (*c.CommentResponse, error) {
	comRes := c.CommentResponse{}
	res, err := s.storage.Comment().WriteComment(repo.Comment{
		Id:     req.Id,
		PostId: req.PostId,
		UserId: req.UserId,
		Text:   req.Text,
	})
	if err != nil {
		log.Println("failed to write comment in service: ", err)
		return &c.CommentResponse{}, err
	}

	comRes.Id = res.Id
	comRes.PostId = res.PostId
	comRes.UserId = res.UserId
	comRes.Text = req.Text
	comRes.CreatedAt = res.CreatedAt

	post, err := s.Client.Post().GetPostForComment(ctx, &p.Request{Str: res.PostId})
	if err != nil {
		log.Println("failed to get post in write comment in service: ", err)
		return &c.CommentResponse{}, err
	}

	comRes.PostTitle = post.Title

	user, err := s.Client.User().GetUserForClient(ctx, &u.Request{Str: res.UserId})
	if err != nil {
		log.Println("failed to get user in write comment in service")
		return &c.CommentResponse{}, err
	}
	comRes.UserName = user.FirstName + " " + user.LastName
	comRes.UserType = user.UserType

	postUser, err := s.Client.User().GetUserForClient(ctx, &u.Request{Str: post.UserId})
	if err != nil {
		log.Println("failed to get post's user in write comment in service: ", err)
		return &c.CommentResponse{}, err
	}
	comRes.PostUserName = postUser.FirstName + " " + postUser.LastName

	return &comRes, nil
}

func (s *CommentService) GetComments(ctx context.Context, req *c.Request) (*c.CommentsResponse, error) {
	coms := c.CommentsResponse{}

	res, err := s.storage.Comment().GetComments(req.Str)
	if err != nil {
		log.Println("failed to get comments in service: ", err)
		return &c.CommentsResponse{}, err
	}

	for _, val := range res {
		coms.Comments = append(coms.Comments, &c.CommentResponse{Id: val.Id, PostId: val.PostId, UserId: val.UserId, Text: val.Text, CreatedAt: val.CreatedAt})
	}

	post, err := s.Client.Post().GetPostForComment(ctx, &p.Request{Str: req.Str})
	if err != nil {
		log.Println("failed to get post in get comments in service: ", err)
		return &c.CommentsResponse{}, err
	}

	for _, comment := range coms.Comments {
		user, err := s.Client.User().GetUserForClient(ctx, &u.Request{Str: comment.UserId})
		if err != nil {
			log.Println("failed to get user in get comments in service: ", err)
			return &c.CommentsResponse{}, err
		}
		comment.UserName = user.FirstName + " " + user.LastName
		comment.UserType = user.UserType
	}

	postUser, err := s.Client.User().GetUserForClient(ctx, &u.Request{Str: post.UserId})
	if err != nil {
		log.Println("failed to get post user in get comments in service: ", err)
		return &c.CommentsResponse{}, err
	}

	for _, comment := range coms.Comments {
		comment.PostTitle = post.Title
		comment.PostUserName = postUser.FirstName + " " + postUser.LastName
	}

	return &coms, nil
}

func (s *CommentService) GetCommentsForPost(ctx context.Context, req *c.Request) (*c.CommentsResponse, error) {
	coms := c.CommentsResponse{}

	res, err := s.storage.Comment().GetComments(req.Str)
	if err != nil {
		log.Println("failed to get comments for post in service: ", err)
		return &c.CommentsResponse{}, err
	}

	for _, val := range res {
		coms.Comments = append(coms.Comments, &c.CommentResponse{Id: val.Id, PostId: val.PostId, UserId: val.UserId, Text: val.Text, CreatedAt: val.CreatedAt})
	}

	return &coms, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, id *c.Request) (*c.CommentResponse, error) {
	comRes := c.CommentResponse{}
	res, err := s.storage.Comment().DeleteComment(id.Str)
	if err != nil {
		log.Println("failed to delete comment service: ", err)
		return &c.CommentResponse{}, err
	}

	comRes.Id = res.Id
	comRes.PostId = res.PostId
	comRes.UserId = res.UserId
	comRes.Text = res.Text
	comRes.CreatedAt = res.CreatedAt

	post, err := s.Client.Post().GetPostForComment(ctx, &p.Request{Str: res.PostId})
	if err != nil {
		log.Println("failed to get post in delete comment service: ", err)
		return &c.CommentResponse{}, err
	}
	comRes.PostTitle = post.Title

	user, err := s.Client.User().GetUserForClient(ctx, &u.Request{Str: res.UserId})
	if err != nil {
		log.Println("failed to get user in delete comment service: ", err)
		return &c.CommentResponse{}, err
	}
	comRes.UserName = user.FirstName + " " + user.LastName
	comRes.UserType = user.UserType

	postUser, err := s.Client.User().GetUserForClient(ctx, &u.Request{Str: post.UserId})
	if err != nil {
		log.Println("failed to get post user in delete comment service: ", err)
		return &c.CommentResponse{}, err
	}
	comRes.PostUserName = postUser.FirstName + " " + postUser.LastName

	return &comRes, nil
}
