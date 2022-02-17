package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"jwtgolang-mongodb/database"
	"jwtgolang-mongodb/helpers"
	"jwtgolang-mongodb/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("user_id")
		defer cancel()

		// from claims token will a check uid in here
		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			response := helpers.ApiResponse("Failed to retrieve data", http.StatusBadRequest, "error", err)
			c.JSON(http.StatusBadRequest, response)
			return
		}
		var user models.ResponseUserLogin

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			response := helpers.ApiResponse("Failed to retrieve data", http.StatusInternalServerError, "error", err)
			c.JSON(http.StatusInternalServerError, response)
		}

		response := helpers.ApiResponse("Successfully", http.StatusOK, "success", user)
		c.JSON(http.StatusOK, response)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			response := helpers.ApiResponse("Failed to retrieve data", http.StatusBadRequest, "error", gin.H{"error": err.Error()})
			c.JSON(http.StatusBadRequest, response)
			return
		}
		cursor, err := userCollection.Find(ctx, bson.M{})
		if err != nil {
			log.Fatal(err)
			response := helpers.ApiResponse("Failed to retrieve data", http.StatusInternalServerError, "error", err)
			c.JSON(http.StatusInternalServerError, response)
		}
		var users []models.ResponseUserLogin
		if err = cursor.All(ctx, &users); err != nil {
			log.Fatal(err)
		}

		response := helpers.ApiResponse("Successfully", http.StatusOK, "success", users)
		c.JSON(http.StatusOK, response)
	}
}
