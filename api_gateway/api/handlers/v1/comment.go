package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/burxondv/new-services/api-gateway/api/handlers/models"
	pc "github.com/burxondv/new-services/api-gateway/genproto/comment"
	l "github.com/burxondv/new-services/api-gateway/pkg/logger"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// Super-Admin | Admin | User
// @Summary Write comment
// @Tags Comment
// @Description Writing new comment for Post
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param CommentInfo body models.CommentRequest true "write comment"
// @Success 201 {object} models.Comment
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/comments [post]
func (h *handlerV1) WriteComment(c *gin.Context) {
	var (
		body        models.CommentRequest
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

	response, err := h.serviceManager.CommentService().WriteComment(context.Background(), &pc.CommentRequest{
		Id:     id.String(),
		PostId: body.PostId,
		UserId: reqId,
		Text:   body.Text,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to write comment", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, models.Comment{
		Id:           response.Id,
		PostId:       response.PostId,
		PostTitle:    response.PostTitle,
		UserId:       response.UserId,
		UserName:     response.UserName,
		UserType:     response.UserType,
		PostUserName: response.PostUserName,
		Text:         response.Text,
		CreatedAt:    response.CreatedAt,
	})
}

// Super-Admin | Admin | User
// @Summary Get Comments by Post ID
// @Tags Comment
// @Description Getting comments by Post Id
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 string models.Comments
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/comments/{id} [get]
func (h *handlerV1) GetComments(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true
	comments := models.Comments{}

	id := c.Param("id")

	response, err := h.serviceManager.CommentService().GetComments(context.Background(), &pc.Request{Str: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get comment by post id", l.Error(err))
		return
	}

	for _, val := range response.Comments {
		com := models.Comment{}
		com.Id = val.Id
		com.PostId = val.PostId
		com.PostTitle = val.PostTitle
		com.PostUserName = val.PostUserName
		com.UserId = val.UserId
		com.UserName = val.UserName
		com.UserType = val.UserType
		com.Text = val.Text
		com.CreatedAt = val.CreatedAt

		comments.Comments = append(comments.Comments, com)
	}

	c.JSON(http.StatusOK, comments)
}

// Super-Admin | Admin | User
// @Summary Delete Comment
// @Tags Comment
// @Description Delete Comment by Id
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Comment Id"
// @Success 200 {object} models.DeletedComment
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/comments/{id} [delete]
func (h *handlerV1) DeleteComment(c *gin.Context) {
	jspbMarshal := protojson.MarshalOptions{}
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	response, err := h.serviceManager.CommentService().DeleteComment(context.Background(), &pc.Request{Str: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete comment", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, models.DeletedComment{
		Id:           response.Id,
		PostId:       response.PostId,
		PostTitle:    response.PostTitle,
		UserId:       response.UserId,
		UserName:     response.UserName,
		UserType:     response.UserType,
		PostUserName: response.PostUserName,
		Text:         response.Text,
		CreatedAt:    response.CreatedAt,
		DeletedAt:    time.Now().Format("2006-01-02 15:04:05"),
	})
}
