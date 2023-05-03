package v1

import (
	"context"
	"net/http"

	"github.com/burxondv/new-services/api-gateway/api/handlers/models"
	"github.com/burxondv/new-services/api-gateway/api/handlers/token"
	pu "github.com/burxondv/new-services/api-gateway/genproto/user"
	l "github.com/burxondv/new-services/api-gateway/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// unauthorized
// @Summary Login
// @Tags Sign-in | Sign-up
// @Description If you have an account, you need to Login
// @Accept json
// @Produce json
// @Param email path string true "email"
// @Param password path string true "password"
// @Success 200 {object} models.LoginResponseModel
// @Router /v1/login/{email}/{password} [get]
func (h *handlerV1) Login(c *gin.Context) {
	var (
		loginResponse             models.LoginResponseModel
		accessToken, refreshToken string

		email    = c.Param("email")
		password = c.Param("password")
	)

	res, err := h.serviceManager.UserService().Login(
		context.Background(), &pu.LoginRequest{
			Email:    email,
			Password: password,
		},
	)

	st, ok := status.FromError(err)
	if !ok || st.Code() == codes.Internal {
		c.JSON(http.StatusInternalServerError, models.StandardErrorModel{
			Error: models.Error{
				Message: st.Message(),
			},
		})
		h.log.Error("failed get client by email", l.Error(err))
		return
	}

	h.jwtHandler = token.JWTHandler{
		SigninKey: h.cfg.SigningKey,
		Sub:       res.Id,
		Iss:       "user",
		Role:      res.UserType,
		Aud: []string{
			"bnnfav_token",
		},
		Log: h.log,
	}

	accessToken, refreshToken, err = h.jwtHandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to access and refresh token", l.Error(err))
		return
	}

	ucReq := &pu.UpdateUserTokensRequest{
		Id:           res.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	newRes, err := h.serviceManager.UserService().UpdateUserTokens(context.Background(), &pu.UpdateUserTokensRequest{
		Id:           ucReq.Id,
		RefreshToken: ucReq.RefreshToken,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user tokens", l.Error(err))
		return
	}

	loginResponse.Id = newRes.Id
	loginResponse.FirstName = newRes.FirstName
	loginResponse.LastName = newRes.LastName
	loginResponse.UserType = newRes.UserType
	loginResponse.Email = newRes.Email
	loginResponse.AccessToken = accessToken
	loginResponse.RefreshToken = refreshToken

	c.JSON(http.StatusOK, loginResponse)
}
