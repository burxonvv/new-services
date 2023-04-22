package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	r "github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/new-york-services/api_gateway/api/handlers/models"
	"github.com/new-york-services/api_gateway/api/handlers/token"
	pu "github.com/new-york-services/api_gateway/genproto/user"
	"github.com/new-york-services/api_gateway/pkg/email"
	"github.com/new-york-services/api_gateway/pkg/etc"
	l "github.com/new-york-services/api_gateway/pkg/logger"
)

// unauthorized
// @Summary register user api
// @Description this api for registers new user
// @Tags Sign-in | Sign-up
// @Accept json
// @Produce json
// @Param body body models.UserRegister true "register user"
// @Success 200 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/register [post]
func (h *handlerV1) Register(c *gin.Context) {
	var (
		body models.RegisterModel
	)

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	existsFirstName, err := h.serviceManager.UserService().CheckField(context.Background(), &pu.CheckFieldRequest{
		Field: "first_name",
		Value: body.FirstName,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed check first name uniques ", l.Error(err))
	}

	if existsFirstName.Exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"info":  "please enter another first name",
		})
		h.log.Error("this first name already exists ", l.Error(err))
	}

	existsEmail, err := h.serviceManager.UserService().CheckField(context.Background(), &pu.CheckFieldRequest{
		Field: "email",
		Value: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed check email uniques ", l.Error(err))
	}

	if existsEmail.Exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"info":  "please enter another email",
		})
		h.log.Error("this email already exists ", l.Error(err))
	}

	code := etc.GenerateCode(6)
	msg := "Subject: Exam email verification\n Your verification code: " + code
	err = email.SendEmail([]string{body.Email}, []byte(msg))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"eroor": err.Error(),
		})
		return
	}
	body.Code = code

	userBodyByte, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"eroor": err.Error(),
		})
		h.log.Error("failed while marshal user body", l.Error(err))
		return
	}

	err = h.redis.SetWithTTL(string(body.Email), string(userBodyByte), 300)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error set to redis user body", l.Error(err))
		return
	}

	c.JSON(http.StatusAccepted, code)
}

// unauthorized
// @Summary verify user api
// @Description this api verifies
// @Tags Sign-in | Sign-up
// @Accept json
// @Produce json
// @Param email path string true "email"
// @Param code path string true "code"
// @Succes 200{object} models.RegisterModel
// @Router /v1/verify/{email}/{code} [get]
func (h *handlerV1) Verify(c *gin.Context) {
	var (
		code  = c.Param("code")
		email = c.Param("email")
		body  models.RegisterModel
	)

	userBody, err := h.redis.Get(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error get from redis by email", l.Error(err))
		return
	}

	byteData, err := r.String(userBody, err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error string redis", l.Error(err))
		return
	}

	err = json.Unmarshal([]byte(byteData), &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while unmarshalling user data", l.Error(err))
		return
	}

	if body.Code != code {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while checking code ", l.Error(err))
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to generating uuid", l.Error(err))
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
	accessTokenString, refreshTokenString, err := h.jwtHandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to generating access token", l.Error(err))
		return
	}

	// Create hash of a password
	hashedPassword, err := etc.GeneratePasswordHash(body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to generating hash password", l.Error(err))
		return
	}

	user, err := h.serviceManager.UserService().CreateUser(context.Background(), &pu.UserResponse{
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
		h.log.Error("error while creating user to db", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.RegisterResponseModel{
		Id:           user.Id,
		AccessToken:  accessTokenString,
		RefreshToken: user.RefreshToken,
	})
}
