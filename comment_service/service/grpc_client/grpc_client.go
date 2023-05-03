package grpcclient

import (
	"fmt"

	"github.com/burxondv/new-services/comment-service/config"
	cp "github.com/burxondv/new-services/comment-service/genproto/post"
	cu "github.com/burxondv/new-services/comment-service/genproto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients interface {
	User() cu.UserServiceClient
	Post() cp.PostServiceClient
}

type ServiceManager struct {
	Config      config.Config
	userService cu.UserServiceClient
	postService cp.PostServiceClient
}

func New(cfg config.Config) (*ServiceManager, error) {
	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.UserServiceHost, cfg.UserServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("user service dial host:%s, port:%s", cfg.UserServiceHost, cfg.UserServicePort)
	}

	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.PostServiceHost, cfg.PostServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("post service dial host:%s, port:%s", cfg.PostServiceHost, cfg.PostServicePort)
	}

	return &ServiceManager{
		Config:      cfg,
		userService: cu.NewUserServiceClient(connUser),
		postService: cp.NewPostServiceClient(connPost),
	}, nil
}

func (s *ServiceManager) User() cu.UserServiceClient {
	return s.userService
}

func (s *ServiceManager) Post() cp.PostServiceClient {
	return s.postService
}
