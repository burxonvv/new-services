package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/burxondv/new-services/api-gateway/api/handlers/models"
	pp "github.com/burxondv/new-services/api-gateway/genproto/post"
	l "github.com/burxondv/new-services/api-gateway/pkg/logger"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// User
// @Summary Create Post
// @Tags Post
// @Descrtiption Create new Post. User Id from Claims
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param PostInfo body models.PostRequest true "Create Post"
// @Success 201 {object} models.Post
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/posts [post]
func (h *handlerV1) CreatePost(c *gin.Context) {
	var (
		body        models.PostRequest
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

	// random uuid...
	id, err := uuid.NewRandom()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to generating uuid", l.Error(err))
		return
	}

	claims := GetClaims(h, c)
	reqId := claims["sub"].(string)

	response, err := h.serviceManager.PostService().CreatePost(context.Background(), &pp.PostRequest{
		Id:          id.String(),
		Title:       body.Title,
		Description: body.Description,
		UserId:      reqId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create post", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, models.Post{
		Id:          response.Id,
		Title:       response.Title,
		Description: response.Description,
		Likes:       response.Likes,
		UserId:      response.UserId,
		UserName:    response.UserName,
		CreatedAt:   response.CreatedAt,
		UpdatedAt:   response.UpdatedAt,
	})
}

// Super-Admin | Admin | User
// @Summary Get Post
// @Tags Post
// @Descrtiption Get Post by Id
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.Post
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/posts/{id} [get]
func (h *handlerV1) GetPost(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	response, err := h.serviceManager.PostService().GetPostById(context.Background(), &pp.Request{Str: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get post by id", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, models.Post{
		Id:          response.Id,
		Title:       response.Title,
		Description: response.Description,
		Likes:       response.Likes,
		UserId:      response.UserId,
		UserName:    response.UserName,
		Comments:    response.Comments,
		CreatedAt:   response.CreatedAt,
		UpdatedAt:   response.UpdatedAt,
	})
}

// User
// @Summary Get Profile posts
// @Tags Post
// @Descrtiption Get own Profile
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} models.Posts
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/posts/profile [get]
func (h *handlerV1) GetProfilePosts(c *gin.Context) {
	claims := GetClaims(h, c)
	reqId := claims["sub"].(string)

	response, err := h.serviceManager.PostService().GetPostByUserId(context.Background(), &pp.Request{Str: reqId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get own posts", l.Error(err))
		return
	}

	posts := models.Posts{}
	for _, val := range response.Posts {
		pt := models.Post{}
		pt.Id = val.Id
		pt.Title = val.Title
		pt.Description = val.Description
		pt.UserId = val.UserId
		pt.UserName = val.UserName
		pt.Likes = val.Likes
		pt.CreatedAt = val.CreatedAt
		pt.UpdatedAt = val.UpdatedAt
		pt.Comments = val.Comments

		posts.Posts = append(posts.Posts, pt)
	}

	c.JSON(http.StatusOK, posts)
}

// Super-Admin | Admin | User
// @Summary Get posts by user Id
// @Tags Post
// @Descrtiption Get post by User Id
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "User Id"
// @Success 200 {object} models.Posts
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/posts/users/{id} [get]
func (h *handlerV1) GetPostsUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true
	Id := c.Param("id")

	response, err := h.serviceManager.PostService().GetPostByUserId(context.Background(), &pp.Request{Str: Id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get posts by user id", l.Error(err))
		return
	}

	posts := models.Posts{}
	for _, val := range response.Posts {
		pt := models.Post{}
		pt.Id = val.Id
		pt.Title = val.Title
		pt.Description = val.Description
		pt.UserId = val.UserId
		pt.UserName = val.UserName
		pt.Likes = val.Likes
		pt.Comments = val.Comments
		pt.CreatedAt = val.CreatedAt
		pt.UpdatedAt = val.UpdatedAt

		posts.Posts = append(posts.Posts, pt)
	}

	c.JSON(http.StatusOK, posts)
}

// User
// @Summary Update User /
// @Tags Post
// @Descrtiption Update user by Id
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param UpdatePost body models.UpdatePostRequest true "Update Post"
// @Success 200 string Success models.Post
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/posts/{id} [put]
func (h *handlerV1) UpdatePost(c *gin.Context) {
	var (
		body        pp.UpdatePostRequest
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

	response, err := h.serviceManager.PostService().UpdatePost(context.Background(), &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update post", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, models.Post{
		Id:          response.Id,
		Title:       response.Title,
		Description: response.Description,
		Likes:       response.Likes,
		UserId:      response.UserId,
		UserName:    response.UserName,
		Comments:    response.Comments,
		CreatedAt:   response.CreatedAt,
		UpdatedAt:   response.UpdatedAt,
	})
}

// Super-Admin | Admin | User
// @Summary delete post
// @Tags Post
// @Descrtiption this method for delete post by ID
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Id"
// @Success 200 string models.DeletedPost
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/posts/{id} [delete]
func (h *handlerV1) DeletePost(c *gin.Context) {
	jspbMarshal := protojson.MarshalOptions{}
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	response, err := h.serviceManager.PostService().DeletePost(context.Background(), &pp.Request{Str: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete post", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, models.DeletedPost{
		Id:          response.Id,
		Title:       response.Title,
		Description: response.Description,
		Likes:       response.Likes,
		UserId:      response.UserId,
		UserName:    response.UserName,
		Comments:    response.Comments,
		CreatedAt:   response.CreatedAt,
		UpdatedAt:   response.UpdatedAt,
		DeletedAt:   time.Now().Format("2006-01-02 15:04:05"),
	})
}
