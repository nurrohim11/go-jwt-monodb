package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	First_name    *string            `json:"first_name" validate:"required,min=2,max=100"`
	Last_name     *string            `json:"last_name" validate:"required"`
	Password      *string            `json:"password" validate:"required,min=6,max=100"`
	Email         *string            `json:"email" validate:"email,required"`
	Phone_number  *string            `json:"phone_number" validate:"required"`
	Access_token  *string            `json:"access_token"`
	Refresh_token *string            `json:"refresh_token"`
	User_type     *string            `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"uuid"`
}

type UserFirebase struct {
	ID            primitive.ObjectID `bson:"_id"`
	First_name    *string            `json:"first_name"`
	Last_name     *string            `json:"last_name"`
	Email         *string            `json:"email" validate:"email,required"`
	Phone_number  *string            `json:"phone_number"`
	Access_token  *string            `json:"access_token"`
	Refresh_token *string            `json:"refresh_token"`
	User_type     string             `json:"user_type" example="USER"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"uuid"`
}

type ResponseUserLogin struct {
	ID           primitive.ObjectID `bson:"_id"`
	First_name   *string            `json:"first_name" validate:"required,min=2,max=100"`
	Last_name    *string            `json:"last_name" validate:"required"`
	Email        *string            `json:"email" validate:"email,required"`
	Phone_number *string            `json:"phone_number" validate:"required"`
	User_type    *string            `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Created_at   time.Time          `json:"created_at"`
	Updated_at   time.Time          `json:"updated_at"`
	User_id      string             `json:"user_id"`
}
