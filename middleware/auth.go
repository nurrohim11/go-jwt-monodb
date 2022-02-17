package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"jwtgolang-mongodb/helpers"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {

		clientToken := c.Request.Header.Get("Authorization")
		// check token
		// check if the first string starting with bearer or not
		if clientToken == "" {
			clientToken = c.Request.Header.Get("token")
		} else if strings.HasPrefix(clientToken, "Bearer ") {
			reqToken := c.Request.Header.Get("Authorization")
			splitToken := strings.Split(reqToken, "Bearer ")
			clientToken = splitToken[1]
		} else {
			response := helpers.ApiResponse("Authentication failed", http.StatusInternalServerError, "error", gin.H{"error": "invalid authorization token"})
			c.JSON(http.StatusInternalServerError, response)
			c.Abort()
			return
		}

		if clientToken == "" {
			response := helpers.ApiResponse("Authentication failed", http.StatusInternalServerError, "error", gin.H{"error": "no Authorization header provided"})
			c.JSON(http.StatusInternalServerError, response)
			c.Abort()
			return
		}
		// handle access token
		claims, err := helpers.ValidateToken(clientToken)

		if err != "" {
			response := helpers.ApiResponse("Authentication failed", http.StatusInternalServerError, "error", gin.H{"error": err})
			c.JSON(http.StatusInternalServerError, response)
			c.Abort()
			return
		}

		fmt.Println(claims)
		if claims.Token_type == "access_token" {
			c.Set("email", claims.Email)
			c.Set("first_name", claims.First_name)
			c.Set("last_name", claims.Last_name)
			c.Set("user_id", claims.Uid)
			c.Set("user_type", claims.User_type)
			c.Next()
		} else if claims.Token_type == "refresh_token" {
			c.Set("token_type", claims.Token_type)
			c.Set("user_id", claims.Uid)
			c.Next()
		}

	}
}
