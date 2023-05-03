package api

import (
	_ "github.com/burxondv/new-services/api-gateway/api/docs" // swag
	"github.com/burxondv/new-services/api-gateway/api/handlers/token"
	v1 "github.com/burxondv/new-services/api-gateway/api/handlers/v1"
	"github.com/burxondv/new-services/api-gateway/api/middleware"
	"github.com/burxondv/new-services/api-gateway/config"
	"github.com/burxondv/new-services/api-gateway/pkg/logger"
	"github.com/burxondv/new-services/api-gateway/services"
	"github.com/burxondv/new-services/api-gateway/storage/repo"
	"github.com/casbin/casbin/v2"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// Option ...
type Option struct {
	Conf            config.Config
	Logger          logger.Logger
	ServiceManager  services.IServiceManager
	InMemoryStorage repo.RedisRepo
	CasbinEnforcer  *casbin.Enforcer
}

// Swagger...
// @title Microservices
// @version 1.0
// @description user post and comment services in REST API
// @host localhost:8080
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	jwtHandler := token.JWTHandler{
		SigninKey: config.Load().SigningKey,
		Log:       option.Logger,
	}

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Conf,
		Redis:          option.InMemoryStorage,
		JWTHandler:     jwtHandler,
		Enforcer:       *option.CasbinEnforcer,
	})

	router.Use(gin.Recovery())
	router.Use(middleware.NewAuthorizer(option.CasbinEnforcer, jwtHandler, option.Conf))

	api := router.Group("/v1")

	// register ...
	api.POST("/register", handlerV1.Register)
	api.GET("/verify/:email/:code", handlerV1.Verify)
	api.GET("/login/:email/:password", handlerV1.Login)

	// role ...
	api.POST("/rbac/add-policy", handlerV1.AddPolicy)
	api.POST("/rbac/remove-policy", handlerV1.RemovePolicy)
	api.POST("/rbac/add-role-user", handlerV1.AddRoleForUser)
	api.POST("rbac/delete-role-user", handlerV1.DeleteRoleForUser)
	api.GET("/rbac/get-policy", handlerV1.GetPolicy)
	api.PUT("/rbac/change-role", handlerV1.ChangeRoleUser)
	api.GET("/rbac/same-role/:role", handlerV1.GetSameRoleUsers)

	// users ...
	api.POST("/users/create", handlerV1.CreateUser)
	api.GET("/users/get-profile", handlerV1.GetProfile)
	api.GET("/users/:id", handlerV1.GetUserById)
	api.GET("/users", handlerV1.GetAllUsers)
	api.PUT("/users", handlerV1.UpdateUser)
	api.DELETE("/users/:id", handlerV1.DeleteUser)

	// posts ...
	api.POST("/posts", handlerV1.CreatePost)
	api.GET("/posts/:id", handlerV1.GetPost)
	api.GET("/posts/profile", handlerV1.GetProfilePosts)
	api.GET("/posts/users/:id", handlerV1.GetPostsUser)
	api.PUT("/posts/:id", handlerV1.UpdatePost)
	api.DELETE("/posts/:id", handlerV1.DeletePost)

	// comment ...
	api.POST("/comments", handlerV1.WriteComment)
	api.GET("/comments/:id", handlerV1.GetComments)
	api.DELETE("/comments/:id", handlerV1.DeleteComment)

	// swagger
	url := ginSwagger.URL("swagger/doc.json")
	api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
