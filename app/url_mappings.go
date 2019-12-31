package app

import (
	"github.com/mfirmanakbar/bookstore_users-api/controllers"
)

func mapUrls() {
	router.GET("/ping", controllers.Ping)
}
