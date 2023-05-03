package v1

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/burxondv/new-services/api-gateway/api/handlers/models"
	pu "github.com/burxondv/new-services/api-gateway/genproto/user"
	l "github.com/burxondv/new-services/api-gateway/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Super-Admin
// @Summary Add Policy
// @Tags RBAC
// @Descrtiption Add Policy from enforcer
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param policy body models.Policy true "Policy"
// @Success 200 string Success
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/rbac/add-policy [post]
func (h *handlerV1) AddPolicy(c *gin.Context) {
	body := models.Policy{}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	ok, err := h.enforcer.AddPolicy(body.User, body.Domain, body.Action)
	if err != nil {
		log.Println("failed to add policy: ", err)
	}

	h.enforcer.SavePolicy()
	fmt.Println(ok)

	c.JSON(http.StatusOK, models.Success{
		Message: "successfully added policy",
	})
}

// Super-Admin
// @Summary Remove Policy
// @Tags RBAC
// @Descrtiption Remove Policy by id
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param policy body models.Policy true "Policy"
// @Success 200 string Success
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/rbac/remove-policy [post]
func (h *handlerV1) RemovePolicy(c *gin.Context) {
	body := models.Policy{}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	ok, err := h.enforcer.RemovePolicy(body.User, body.Domain, body.Action)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to remove policy", l.Error(err))
		return
	}

	h.enforcer.SavePolicy()
	fmt.Println(ok)

	c.JSON(http.StatusOK, models.Success{
		Message: "successfully removed policy",
	})
}

// Super-Admin
// @Summary add role
// @Tags RBAC
// @Descrtiption this method for add role for user
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role body models.RoleRequest true "Role"
// @Success 200 string Success
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/rbac/add-role-user [post]
func (h *handlerV1) AddRoleForUser(c *gin.Context) {
	body := models.RoleRequest{}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	ok, err := h.enforcer.AddRoleForUser(body.Id, body.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to add role for user", l.Error(err))
		return
	}

	fmt.Println(ok)

	c.JSON(http.StatusOK, models.Success{
		Message: "successfully added role",
	})
}

// Super-Admin
// @Summary delete role
// @Tags RBAC
// @Descrtiption this method for delete role fom user
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role body models.RoleRequest true "Role"
// @Success 200 string Success
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/rbac/delete-role-user [post]
func (h *handlerV1) DeleteRoleForUser(c *gin.Context) {
	body := models.RoleRequest{}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	ok, err := h.enforcer.DeleteRoleForUser(body.Id, body.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete role for user", l.Error(err))
		return
	}

	fmt.Println(ok)

	c.JSON(http.StatusOK, models.Success{
		Message: "successfully deleted role",
	})
}

// @Summary get policy
// @Description this method for get policy
// @Tags RBAC
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Router /v1/rbac/get-policy [get]
func (h *handlerV1) GetPolicy(c *gin.Context) {
	data := h.enforcer.GetPolicy()
	c.JSON(http.StatusOK, data)
}

// Super-Admin
// @Summary change role
// @Tags RBAC
// @Descrtiption this method for change role user
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param ChangeUserRole body models.RoleRequest true "change Role"
// @Success 200 string models.User
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/rbac/change-role [put]
func (h *handlerV1) ChangeRoleUser(c *gin.Context) {
	body := models.RoleRequest{}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	response, err := h.serviceManager.UserService().ChangeRoleUser(context.Background(), &pu.ChangeRoleRequest{
		Id:   body.Id,
		Role: body.Role,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to add role user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, models.User{
		Id:        response.Id,
		FirstName: response.FirstName,
		LastName:  response.LastName,
		UserType:  response.UserType,
		Email:     response.Email,
		CreatedAt: response.CreatedAt,
		UpdatedAt: response.UpdatedAt,
	})
}

// Super-Admin
// @Summary get the same role users
// @Tags RBAC
// @Descrtiption this method for get the same role users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role path string true "Role"
// @Success 200 {object} models.Users
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/rbac/same-role/{role} [get]
func (h *handlerV1) GetSameRoleUsers(c *gin.Context) {
	users := models.Users{}
	res, err := h.serviceManager.UserService().GetSameRoleUsers(context.Background(), &pu.Request{Str: c.Param("role")})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get the same role users", l.Error(err))
		return
	}

	for _, val := range res.Users {
		newUser := models.User{
			Id:        val.Id,
			FirstName: val.FirstName,
			LastName:  val.LastName,
			UserType:  val.UserType,
			Email:     val.Email,
			CreatedAt: val.CreatedAt,
			UpdatedAt: val.UpdatedAt,
		}

		users.Users = append(users.Users, newUser)
	}

	c.JSON(http.StatusOK, users)
}
