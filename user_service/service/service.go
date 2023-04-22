package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"

	p "github.com/new-york-services/user_service/genproto/post"
	u "github.com/new-york-services/user_service/genproto/user"
	"github.com/new-york-services/user_service/pkg/logger"
	grpcclient "github.com/new-york-services/user_service/service/grpc_client"
	"github.com/new-york-services/user_service/storage"
	"github.com/new-york-services/user_service/storage/repo"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	storage storage.IStorage
	Logger  logger.Logger
	Client  grpcclient.Clients
}

func NewUserService(db *sqlx.DB, log logger.Logger, client grpcclient.Clients) *UserService {
	return &UserService{
		storage: storage.NewStoragePg(db),
		Logger:  log,
		Client:  client,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *u.UserResponse) (*u.UserResponse, error) {
	res, err := s.storage.User().CreateUser(repo.User{
		Id:           req.Id,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Password:     req.Password,
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		log.Println("failed to create user in service: ", err)
		return &u.UserResponse{}, err
	}

	return &u.UserResponse{
		Id:           res.Id,
		FirstName:    res.FirstName,
		LastName:     res.LastName,
		UserType:     res.UserType,
		Email:        res.Email,
		Password:     res.Password,
		RefreshToken: res.RefreshToken,
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}, nil
}

func (s *UserService) GetUserById(ctx context.Context, req *u.Request) (*u.UserResponse, error) {
	userResp := u.UserResponse{}
	res, err := s.storage.User().GetUserById(req.Str)
	if err != nil {
		log.Println("failed to get user in service: ", err)
		return &u.UserResponse{}, err
	}
	userResp.Id = res.Id
	userResp.FirstName = res.FirstName
	userResp.LastName = res.LastName
	userResp.UserType = res.UserType
	userResp.Email = res.Email
	userResp.RefreshToken = res.RefreshToken
	userResp.CreatedAt = res.CreatedAt
	userResp.UpdatedAt = res.UpdatedAt

	postRes, err := s.Client.Post().GetPostForUser(ctx, &p.Request{Str: req.Str})
	if err != nil {
		log.Println("failed to get post in user service: ", err)
		return &u.UserResponse{}, err
	}

	userResp.Posts = int64(len(postRes.Posts))

	return &userResp, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, req *u.Request) (*u.UserResponse, error) {
	userResp := u.UserResponse{}
	res, err := s.storage.User().GetUserByEmail(req.Str)
	if err != nil {
		log.Println("failed to get user by email: ", err)
		return &u.UserResponse{}, err
	}

	userResp.Id = res.Id
	userResp.FirstName = res.FirstName
	userResp.LastName = res.LastName
	userResp.UserType = res.UserType
	userResp.Email = res.Email
	userResp.RefreshToken = res.RefreshToken
	userResp.CreatedAt = res.CreatedAt
	userResp.UpdatedAt = res.UpdatedAt

	posts, err := s.Client.Post().GetPostForUser(ctx, &p.Request{Str: userResp.Id})
	if err != nil {
		log.Println("failed to get post in user email service: ", err)
		return &u.UserResponse{}, err
	}

	userResp.Posts = int64(len(posts.Posts))

	return &userResp, nil
}

func (s *UserService) GetUserForClient(ctx context.Context, req *u.Request) (*u.UserResponse, error) {
	userResp := u.UserResponse{}
	res, err := s.storage.User().GetUserById(req.Str)
	if err != nil {
		log.Println("failed to get user for clients in service: ", err)
		return &u.UserResponse{}, err
	}

	userResp.Id = res.Id
	userResp.FirstName = res.FirstName
	userResp.LastName = res.LastName
	userResp.UserType = res.UserType
	userResp.Email = res.Email
	userResp.RefreshToken = res.RefreshToken
	userResp.CreatedAt = res.CreatedAt
	userResp.UpdatedAt = res.UpdatedAt

	return &userResp, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, req *u.GetUsersRequest) (*u.UsersResponse, error) {
	usersResp := u.UsersResponse{}
	res, err := s.storage.User().GetAllUsers(req.Page, req.Limit)
	if err != nil {
		log.Println("failed to get all users in service: ", err)
		return &u.UsersResponse{}, err
	}

	for _, user := range res {
		userResp := u.UserResponse{
			Id:        user.Id,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			UserType:  user.UserType,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		post, err := s.Client.Post().GetPostForUser(ctx, &p.Request{Str: user.Id})
		if err != nil {
			log.Println("failed to get post in user servcise: ", err)
			return &u.UsersResponse{}, err
		}

		userResp.Posts = int64(len(post.Posts))

		usersResp.Users = append(usersResp.Users, &userResp)
	}

	return &usersResp, nil
}

func (s *UserService) SearchUsers(ctx context.Context, req *u.Request) (*u.UsersResponse, error) {
	usersResp := u.UsersResponse{}
	res, err := s.storage.User().SearchUsers(req.Str)
	if err != nil {
		log.Println("failed to searching user by name: ", err)
		return &u.UsersResponse{}, err
	}

	for _, user := range res {
		userResp := u.UserResponse{
			Id:        user.Id,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			UserType:  user.UserType,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		post, err := s.Client.Post().GetPostForUser(ctx, &p.Request{Str: user.Id})
		if err != nil {
			log.Println("failed to get post in user service: ", err)
			return &u.UsersResponse{}, err
		}

		userResp.Posts = int64(len(post.Posts))

		usersResp.Users = append(usersResp.Users, &userResp)
	}

	return &usersResp, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *u.UpdateUserRequest) (*u.UserResponse, error) {
	userResp := u.UserResponse{}
	res, err := s.storage.User().UpdateUser(repo.User{
		Id:        req.Id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	})
	if err != nil {
		log.Println("failed to update user in service: ", err)
		return &u.UserResponse{}, err
	}
	fmt.Println(req)
	userResp.Id = req.Id
	userResp.FirstName = res.FirstName
	userResp.LastName = res.LastName
	userResp.UserType = res.UserType
	userResp.Email = res.Email
	userResp.CreatedAt = res.CreatedAt
	userResp.UpdatedAt = res.UpdatedAt

	post, err := s.Client.Post().GetPostForUser(ctx, &p.Request{Str: req.Id})
	if err != nil {
		log.Println("failed to get post in user delete service: ", err)
		return &u.UserResponse{}, err
	}

	userResp.Posts = int64(len(post.Posts))

	return &userResp, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *u.Request) (*u.UserResponse, error) {
	userResp := u.UserResponse{}
	res, err := s.storage.User().DeleteUser(req.Str)
	if err != nil {
		log.Println("failed to delete user: ", err)
		return &u.UserResponse{}, err
	}

	userResp.Id = req.Str
	userResp.FirstName = res.FirstName
	userResp.LastName = res.LastName
	userResp.UserType = res.UserType
	userResp.Email = res.Email
	userResp.CreatedAt = res.CreatedAt
	userResp.UpdatedAt = res.UpdatedAt

	post, err := s.Client.Post().GetPostForUser(ctx, &p.Request{Str: req.Str})
	if err != nil {
		log.Println("failed to get post in user service: ", err)
		return &u.UserResponse{}, err
	}

	userResp.Posts = int64(len(post.Posts))

	return &userResp, nil
}

func (s *UserService) CheckField(ctx context.Context, req *u.CheckFieldRequest) (*u.CheckFieldResponse, error) {
	res, err := s.storage.User().CheckField(req)
	if err != nil {
		log.Println("failed to check field: ", err)
		return &u.CheckFieldResponse{}, err
	}

	return res, nil
}

func (s *UserService) Login(ctx context.Context, req *u.LoginRequest) (*u.LoginResponse, error) {
	req.Email = strings.ToLower(req.Email)
	req.Email = strings.TrimSpace(req.Email)

	user, err := s.storage.User().GetUserByEmail(req.Email)

	if err == sql.ErrNoRows {
		log.Println("failed to get user by email, not found: ", err)
		return &u.LoginResponse{}, err
	} else if err != nil {
		log.Println("failed to get user by email, internal server error: ", err)
		return &u.LoginResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.Println("failed to compare password: ", err)
		return &u.LoginResponse{}, err
	}

	return &u.LoginResponse{
		Id:           user.Id,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		UserType:     user.UserType,
		Email:        user.Email,
		Password:     user.Password,
		RefreshToken: user.RefreshToken,
	}, nil
}

func (s *UserService) UpdateUserTokens(ctx context.Context, req *u.UpdateUserTokensRequest) (*u.UserResponse, error) {
	res, err := s.storage.User().UpdateUserTokens(req)
	if err != nil {
		log.Println("failed to update user tokens in user service: ", err)
		return nil, err
	}

	return res, nil
}

func (s *UserService) ChangeRoleUser(ctx context.Context, req *u.ChangeRoleRequest) (*u.UserResponse, error) {
	res, err := s.storage.User().ChangeRoleUser(req)
	if err != nil {
		log.Println("failed to update user role in user service: ", err)
		return nil, err
	}

	return res, nil
}

func (s *UserService) GetSameRoleUsers(ctx context.Context, req *u.Request) (*u.UsersResponse, error) {
	res, err := s.storage.User().GetSameRoleUsers(req.Str)
	if err != nil {
		log.Println("failed to get the same role users in service: ", err)
		return nil, err
	}

	return res, nil
}
