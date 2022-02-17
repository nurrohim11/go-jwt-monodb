package routes

import (
	"jwtgolang-mongodb/controllers"

	"github.com/gin-gonic/gin"
)

func Auth(router *gin.Engine) {
	router.POST("/auth/signup", controllers.Signup())
	router.POST("/auth/login", controllers.Login())
	router.POST("/auth/login_firebase", controllers.LoginByFirebase())
}
