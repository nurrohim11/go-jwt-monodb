package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"time"

	"jwtgolang-mongodb/helpers"
	"jwtgolang-mongodb/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPass string, providedPass string) (passIsValid bool, msg string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPass), []byte(providedPass))
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(userPass, providedPass)
		return false, "email or password is incorrect"
	} else {
		return true, ""
	}

}
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			response := helpers.ApiResponse("Register is failed", http.StatusBadRequest, "error", gin.H{"error": err.Error()})
			c.JSON(http.StatusBadRequest, response)
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			response := helpers.ApiResponse("Register is failed", http.StatusBadRequest, "error", gin.H{"error": validationErr.Error()})
			c.JSON(http.StatusBadRequest, response)
			return
		}

		_, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			fmt.Println(err)
			response := helpers.ApiResponse("Register is failed", http.StatusInternalServerError, "error", gin.H{"error": "error occured while chacking email existance"})
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		count, err := userCollection.CountDocuments(ctx, bson.M{"phone_number": user.Phone_number})
		if err != nil {
			fmt.Println(err)
			response := helpers.ApiResponse("Register is failed", http.StatusInternalServerError, "error", gin.H{"error": "error occured while chacking phone number existance"})
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		if count > 0 {
			response := helpers.ApiResponse("Register is failed", http.StatusInternalServerError, "error", gin.H{"error": "This email or phone number already exist"})
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		local, _ := time.LoadLocation("Asia/Jakarta")
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().In(local).Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().In(local).Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		uuid := user.User_id
		if user.User_id == "" {
			uuid = user.ID.Hex()
		}
		user.User_id = uuid
		accessToken, refreshToken, _ := helpers.GenerateAllToken(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)
		user.Access_token = &accessToken
		user.Refresh_token = &refreshToken

		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := "User item was not created"
			response := helpers.ApiResponse("Register is failed", http.StatusInternalServerError, "error", gin.H{"error": msg})
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		defer cancel()

		response := helpers.ApiResponse("Register is successfully", http.StatusOK, "success", user)
		c.JSON(http.StatusOK, response)
	}
}

func LoginByFirebase() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		// var user models.User
		var userFirebase models.UserFirebase
		local, _ := time.LoadLocation("Asia/Jakarta")
		defer cancel()

		if err := c.BindJSON(&userFirebase); err != nil {
			response := helpers.ApiResponse("Login is failed", http.StatusBadRequest, "error", gin.H{"error": err.Error()})
			c.JSON(http.StatusBadRequest, response)
			return
		}

		validationErr := validate.Struct(userFirebase)
		if validationErr != nil {
			response := helpers.ApiResponse("Login is failed", http.StatusBadRequest, "error", gin.H{"error": validationErr.Error()})
			c.JSON(http.StatusBadRequest, response)
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"user_id": userFirebase.User_id}).Decode(&userFirebase)
		if err == nil {

			userFirebase.Created_at, _ = time.Parse(time.RFC3339, time.Now().In(local).Format(time.RFC3339))
			userFirebase.Updated_at, _ = time.Parse(time.RFC3339, time.Now().In(local).Format(time.RFC3339))
			userFirebase.ID = primitive.NewObjectID()
			// userFirebase.User_id = user.ID.Hex()
			accessToken, refreshToken, _ := helpers.GenerateAllToken(*userFirebase.Email, *userFirebase.First_name, *userFirebase.Last_name, userFirebase.User_type, userFirebase.User_id)
			userFirebase.Access_token = &accessToken
			userFirebase.Refresh_token = &refreshToken

			response := helpers.ApiResponse("Login is successfully", http.StatusOK, "success", userFirebase)
			c.JSON(http.StatusOK, response)
			return
		}

		err = userCollection.FindOne(ctx, bson.M{"email": userFirebase.Email}).Decode(&userFirebase)
		if err == nil {

			userFirebase.Created_at, _ = time.Parse(time.RFC3339, time.Now().In(local).Format(time.RFC3339))
			userFirebase.Updated_at, _ = time.Parse(time.RFC3339, time.Now().In(local).Format(time.RFC3339))
			userFirebase.ID = primitive.NewObjectID()
			// userFirebase.User_id = user.ID.Hex()
			accessToken, refreshToken, _ := helpers.GenerateAllToken(*userFirebase.Email, *userFirebase.First_name, *userFirebase.Last_name, userFirebase.User_type, userFirebase.User_id)
			userFirebase.Access_token = &accessToken
			userFirebase.Refresh_token = &refreshToken

			response := helpers.ApiResponse("Login is successfully", http.StatusOK, "success", userFirebase)
			c.JSON(http.StatusOK, response)
			return
		}

		userFirebase.Created_at, _ = time.Parse(time.RFC3339, time.Now().In(local).Format(time.RFC3339))
		userFirebase.Updated_at, _ = time.Parse(time.RFC3339, time.Now().In(local).Format(time.RFC3339))
		userFirebase.ID = primitive.NewObjectID()
		userFirebase.User_type = "User"
		accessToken, refreshToken, _ := helpers.GenerateAllToken(*userFirebase.Email, *userFirebase.First_name, *userFirebase.Last_name, userFirebase.User_type, userFirebase.User_id)
		userFirebase.Access_token = &accessToken
		userFirebase.Refresh_token = &refreshToken

		_, insertErr := userCollection.InsertOne(ctx, userFirebase)
		if insertErr != nil {
			msg := "User item was not created"
			response := helpers.ApiResponse("Login is failed", http.StatusInternalServerError, "error", gin.H{"error": msg})
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		defer cancel()

		response := helpers.ApiResponse("Login is successfully", http.StatusOK, "success", userFirebase)
		c.JSON(http.StatusOK, response)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		var foundUser models.User
		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			response := helpers.ApiResponse("Login is failed", http.StatusBadRequest, "error", err.Error())
			c.JSON(http.StatusBadRequest, response)
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			response := helpers.ApiResponse("Email or password is incorrect", http.StatusBadRequest, "error", "Email not found")
			c.JSON(http.StatusBadRequest, response)
			defer cancel()
			return
		}

		passIsValid, msg := VerifyPassword(*foundUser.Password, *user.Password)
		defer cancel()
		if !passIsValid {
			resp := helpers.ApiResponse("Login is failed", http.StatusInternalServerError, "error", gin.H{"error": msg})
			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		if foundUser.Email == nil {
			resp := helpers.ApiResponse("Login is failed", http.StatusInternalServerError, "error", gin.H{"error": "user not found"})
			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		accessToken, refreshToken, _ := helpers.GenerateAllToken(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)
		helpers.UpdateAllTokens(accessToken, refreshToken, foundUser.User_id)

		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
		if err != nil {
			resp := helpers.ApiResponse("Login is failed", http.StatusInternalServerError, "error", gin.H{"error": gin.H{"error": err.Error()}})
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		response := helpers.ApiResponse("Login is successfully", http.StatusOK, "success", foundUser)
		c.JSON(http.StatusOK, response)
	}
}

func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("token_type") == "refresh_token" {
			var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
			var user models.User
			user_id := c.GetString("user_id")
			id, _ := primitive.ObjectIDFromHex(user_id)
			err := userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
			if err != nil {
				resp := helpers.ApiResponse("refresh token is failed", http.StatusBadRequest, "error", gin.H{"error": err.Error()})
				c.JSON(http.StatusBadRequest, resp)
				defer cancel()
				return
			}

			accessToken, refreshToken, _ := helpers.GenerateAllToken(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)
			helpers.UpdateAllTokens(accessToken, refreshToken, user.User_id)

			var tokenModel = models.Token{Access_token: accessToken, Refresh_token: refreshToken}
			response := helpers.ApiResponse("Refresh token is successfully", http.StatusOK, "success", tokenModel)
			c.JSON(http.StatusOK, response)
			defer cancel()

		} else {
			resp := helpers.ApiResponse("refresh token is failed", http.StatusUnauthorized, "error", gin.H{"error": "invald refresh token"})
			c.JSON(http.StatusUnauthorized, resp)
		}
	}
}
