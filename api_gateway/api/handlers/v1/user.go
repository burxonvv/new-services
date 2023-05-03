package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/burxondv/new-services/api-gateway/api/handlers/models"
	"github.com/burxondv/new-services/api-gateway/api/handlers/token"
	pu "github.com/burxondv/new-services/api-gateway/genproto/user"
	"github.com/burxondv/new-services/api-gateway/pkg/etc"
	l "github.com/burxondv/new-services/api-gateway/pkg/logger"
	"github.com/burxondv/new-services/api-gateway/pkg/utils"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// Super-Admin | Admin
// @Summary create user
// @Tags User
// @Descrtiption this method for create a new user
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param UserInfo body models.UserRegister true "Create User"
// @Success 201 {object} models.User
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/users/create [post]
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		body        models.UserRegister
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	// random uuid...
	id, err := uuid.NewRandom()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to generating uuid", l.Error(err))
		return
	}

	// password hashing...
	hashedPassword, err := etc.GeneratePasswordHash(body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to generating hash password", l.Error(err))
		return
	}

	h.jwtHandler = token.JWTHandler{
		SigninKey: h.cfg.SigningKey,
		Sub:       id.String(),
		Iss:       "user",
		Role:      "user",
		Aud: []string{
			"bnnfav_token",
		},
		Log: h.log,
	}

	// Create access and refresh tokens
	_, refreshTokenString, err := h.jwtHandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to generating access token", l.Error(err))
		return
	}

	response, err := h.serviceManager.UserService().CreateUser(context.Background(), &pu.UserResponse{
		Id:           id.String(),
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		Email:        body.Email,
		Password:     string(hashedPassword),
		RefreshToken: refreshTokenString,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create user", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, models.User{
		Id:        response.Id,
		FirstName: response.FirstName,
		LastName:  response.LastName,
		UserType:  response.UserType,
		Email:     response.Email,
		CreatedAt: response.CreatedAt,
		UpdatedAt: response.UpdatedAt,
	})
}

// User
// @Summary get own user profile
// @Tags User
// @Descrtiption this method for get own user profile
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} models.User
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/users/get-profile [get]
func (h *handlerV1) GetProfile(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	claims := GetClaims(h, c)
	reqId := claims["sub"].(string)

	response, err := h.serviceManager.UserService().GetUserById(context.Background(), &pu.Request{Str: reqId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user by id", l.Error(err))
		return
	}

	user := models.User{
		Id:        response.Id,
		FirstName: response.FirstName,
		LastName:  response.LastName,
		UserType:  response.UserType,
		Email:     response.Email,
		Posts:     response.Posts,
		CreatedAt: response.CreatedAt,
		UpdatedAt: response.UpdatedAt,
	}

	c.JSON(http.StatusOK, user)
}

// Super-Admin | Admin | User
// @Summary get user by id
// @Tags User
// @Descrtiption this method for get user by id
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Id"
// @Success 200 {object} models.User
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/users/{id} [get]
func (h *handlerV1) GetUserById(c *gin.Context) {
	var (
		jspbMarshal = protojson.MarshalOptions{}
	)
	jspbMarshal.UseProtoNames = true

	res, err := h.serviceManager.UserService().GetUserById(context.Background(), &pu.Request{Str: c.Param("id")})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user by id", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, models.User{
		Id:        res.Id,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		UserType:  res.UserType,
		Email:     res.Email,
		Posts:     res.Posts,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	})
}

// Super-Admin | Admin | User
// @Summary get all users
// @Tags User
// @Descrtiption this method for get all users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param limit query int true "Limit"
// @Param page query int true "Page"
// @Success 200 {object} models.Users
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/users [get]
func (h *handlerV1) GetAllUsers(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	users := models.Users{}

	params, errStr := utils.ParseQueryParams(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errStr[0],
		})
		h.log.Error("failed to parse query params to json: " + errStr[0])
		return
	}

	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	response, err := h.serviceManager.UserService().GetAllUsers(context.Background(), &pu.GetUsersRequest{Limit: params.Limit, Page: params.Page})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get all users", l.Error(err))
		return
	}

	for _, val := range response.Users {
		newUser := models.User{}
		newUser.Id = val.Id
		newUser.FirstName = val.FirstName
		newUser.LastName = val.LastName
		newUser.UserType = val.UserType
		newUser.Email = val.Email
		newUser.Posts = val.Posts
		newUser.CreatedAt = val.CreatedAt
		newUser.UpdatedAt = val.UpdatedAt

		users.Users = append(users.Users, newUser)
	}

	c.JSON(http.StatusOK, users)
}

// User
// @Summary update user
// @Tags User
// @Descrtiption this method for update a user
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param NewUser body models.UpdateUserRequest true "Update User"
// @Success 200 string Success models.Success
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/users [put]
func (h *handlerV1) UpdateUser(c *gin.Context) {
	var (
		body        models.UpdateUserRequest
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind JSON", l.Error(err))
		return
	}

	claims := GetClaims(h, c)
	reqId := claims["sub"].(string)

	res, err := h.serviceManager.UserService().UpdateUser(context.Background(), &pu.UpdateUserRequest{
		Id:        reqId,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, models.User{
		Id:        res.Id,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		UserType:  res.UserType,
		Email:     res.Email,
		Posts:     res.Posts,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	})
}

// Super-Admin | Admin
// @Summary delete user
// @Tags User
// @Descrtiption this method for delete user by ID
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Id"
// @Success 200 string models.DeletedUser
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/users/{id} [delete]
func (h *handlerV1) DeleteUser(c *gin.Context) {
	jspbMarshal := protojson.MarshalOptions{}
	jspbMarshal.UseProtoNames = true

	response, err := h.serviceManager.UserService().DeleteUser(context.Background(), &pu.Request{Str: c.Param("id")})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete user", l.Error(err))
		return
	}

	now := time.Now()
	user := models.DeletedUser{
		Id:        response.Id,
		FirstName: response.FirstName,
		LastName:  response.LastName,
		UserType:  response.UserType,
		Email:     response.Email,
		Posts:     response.Posts,
		CreatedAt: response.CreatedAt,
		UpdatedAt: response.UpdatedAt,
		DeletedAt: now.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, user)
}
