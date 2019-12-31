package users

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mfirmanakbar/bookstore_users-api/domain/users"
	"github.com/mfirmanakbar/bookstore_users-api/services"
	"io/ioutil"
	"net/http"
)

func CreateUser(c *gin.Context) {
	var user users.User

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		//TODO : Handle error
		return
	}
	if err := json.Unmarshal(bytes, &user); err != nil {
		//TODO : Handle json error
		fmt.Println(err.Error())
		return
	}

	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		//TODO : Handle user creation error
		return
	}
	c.JSON(http.StatusCreated, result)
}

func GetUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Implement GetUser")
}
