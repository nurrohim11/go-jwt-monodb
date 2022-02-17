package main

import (
	"log"
	"net/http"
	"os"

	"jwtgolang-mongodb/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
	port := os.Getenv("PORT")

	if port == "" {
		port = "9090"
	}

	// testMongo()

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "API is running and accessible"})
	})

	router.GET("api-1", func(c *gin.Context) { c.JSON(200, gin.H{"success": "access granted for api-1"}) })
	router.GET("api-2", func(c *gin.Context) { c.JSON(200, gin.H{"success": "access granted for api-2 "}) })

	routes.Auth(router)
	routes.User(router)

	router.Run("localhost:" + port)
}
