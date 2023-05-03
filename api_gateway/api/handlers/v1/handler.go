package v1

import (
	"errors"
	"net/http"

	//"strings"

	"github.com/burxondv/new-services/api-gateway/api/handlers/models"
	"github.com/burxondv/new-services/api-gateway/api/handlers/token"
	"github.com/burxondv/new-services/api-gateway/config"
	"github.com/burxondv/new-services/api-gateway/pkg/logger"
	"github.com/burxondv/new-services/api-gateway/services"
	"github.com/burxondv/new-services/api-gateway/storage/repo"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type handlerV1 struct {
	log            logger.Logger
	serviceManager services.IServiceManager
	cfg            config.Config
	redis          repo.RedisRepo
	jwtHandler     token.JWTHandler
	enforcer       casbin.Enforcer
}

type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	Redis          repo.RedisRepo
	JWTHandler     token.JWTHandler
	Enforcer       casbin.Enforcer
}

func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
		redis:          c.Redis,
		jwtHandler:     c.JWTHandler,
		enforcer:       c.Enforcer,
	}
}

func GetClaims(h *handlerV1, c *gin.Context) jwt.MapClaims {
	var (
		ErrUnauthorized = errors.New("unauthorized")
		authorization   = models.GetProfileByJWTRequest{}
		claims          jwt.MapClaims
		err             error
	)

	authorization.Token = c.GetHeader("Authorization")
	if c.Request.Header.Get("Authorization") == "" {
		c.JSON(http.StatusUnauthorized, models.StandardErrorModel{
			Error: models.Error{Message: "error unauthorized in get header"},
		})
		h.log.Error("Unauthorized request: ", logger.Error(ErrUnauthorized))
		return nil
	}

	// this need only in Postman, because in Postman's token have Bearer word, and we need to trim this word ...
	//authorization.Token = strings.TrimSpace(strings.Trim(authorization.Token, "Bearer"))

	h.jwtHandler.Token = authorization.Token
	claims, err = h.jwtHandler.ExtractClaims()
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.StandardErrorModel{
			Error: models.Error{Message: "error unmarshalling in extract claims"},
		})
		h.log.Error("Unauthorized request: ", logger.Error(ErrUnauthorized))
		return nil
	}

	return claims
}
