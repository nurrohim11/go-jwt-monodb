package helpers

import (
	"context"

	"log"
	"os"
	"time"

	"jwtgolang-mongodb/database"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetails struct {
	Email        string
	First_name   string
	Last_name    string
	Uid          string
	User_type    string
	Firebase_uid string
	Token_type   string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GenerateAllToken(email string, firstName string, lastName string, userType string, uid string) (signedAccessToken string, signedRefreshToken string, err error) {

	local, _ := time.LoadLocation("Asia/Jakarta")

	accessClaims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		User_type:  userType,
		Uid:        uid,
		Token_type: "access_token",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().In(local).Add(time.Minute * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		Uid:        uid,
		Token_type: "refresh_token",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().In(local).Add(time.Hour * time.Duration(24*7)).Unix(),
		},
	}

	// firebaseClaims := map[string]interface{}{
	// 	"Email":      email,
	// 	"First_name": firstName,
	// 	"Last_name":  lastName,
	// 	"User_type":  userType,
	// 	"Uid":        uid,
	// 	"Token_type": "access_token",
	// }

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Panic(err)
		return
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Panic(err)
		return
	}

	// firebaseAuth := configs.SetupFirebase()
	// firebaseToken, errors := firebaseAuth.CustomTokenWithClaims(context.Background(), uid, firebaseClaims)
	// if errors != nil {
	// 	log.Panic(errors)
	// 	return
	// }

	return accessToken, refreshToken, err

}

func UpdateAllTokens(signedAccessToken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	// var updateObj primitive.D

	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	id, _ := primitive.ObjectIDFromHex(userId)
	_, err := userCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"access_token": signedAccessToken, "refresh_token": signedRefreshToken, "update_at": updated_at}},
	)

	defer cancel()

	if err != nil {
		log.Panic(err)
	}

}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = err.Error()
		return
	}

	return claims, msg
}
