package users

import (
	"github.com/gin-gonic/gin"
	"github.com/mfirmanakbar/bookstore_users-api/domain/users"
	"github.com/mfirmanakbar/bookstore_users-api/services"
	"github.com/mfirmanakbar/bookstore_users-api/utils/erros"
	"net/http"
)

func CreateUser(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		//TODO : return bad request for the caller
		restErr := erros.RestErr{
			Message: "invalid json body",
			Status:  http.StatusBadRequest,
			Error:   "Bad_Request",
		}
		c.JSON(restErr.Status, restErr)
		return
	}

	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		//TODO : Handle user creation error
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func GetUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Implement GetUser")
}
