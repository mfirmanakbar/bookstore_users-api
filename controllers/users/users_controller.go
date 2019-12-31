package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Implement CreateUser")
}

func GetUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Implement GetUser")
}
