package app

import (
	"github.com/gin-gonic/gin"
	"github.com/mfirmanakbar/bookstore_users-api/logger"
)

var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()
	logger.Info("about to start the Application ...")
	router.Run(":8282")
}
