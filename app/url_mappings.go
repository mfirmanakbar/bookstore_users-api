package app

import (
	"github.com/mfirmanakbar/bookstore_users-api/controllers/ping"
	"github.com/mfirmanakbar/bookstore_users-api/controllers/users"
)

// 1. PUT (Full) --> all field on json must be defined, else it will update all column by empty or 0 on database
// 2. PATCH (Partial) --> just only field defined on json will update on database
func mapUrls() {
	router.GET("/ping", ping.Ping)

	router.POST("/users", users.CreateUser)
	router.GET("/users/:user_id", users.GetUser)
	router.PUT("/users/:user_id", users.UpdateUser)
	router.PATCH("/users/:user_id", users.UpdateUser)
}
