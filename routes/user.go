package routes

import (
	"jwtgolang-mongodb/controllers"
	"jwtgolang-mongodb/middleware"

	"github.com/gin-gonic/gin"
)

func User(router *gin.Engine) {
	router.Use(middleware.Authenticate())
	router.GET("/user/refresh_token", controllers.RefreshToken())
	router.GET("/users", controllers.GetUsers())
	router.GET("/user/:user_id", controllers.GetUser())
}
