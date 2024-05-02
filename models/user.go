package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	First_Name string             `json:"firstname" bson:"firstname"`
	Last_Name  string             `json:"lastname" bson:"lastname"`
	Username   string             `json:"username" bson:"username" validate:"required,min=5,max=10"`
	Email      string             `json:"email" bson:"email" validate:"required"`
	PhoneNo    string             `json:"phoneno" bson:"phoneno" validate:"min=10 max=14"`
	Password   string             `json:"password" bson:"password" validate:"required,min=6,max=15"`
	Token      string             `json:"token" bson:"token"`
	UserId     string             `json:"userId" bson:"userId"`
	CreatedAt  time.Time          `json:"createdat" bson:"createdat"`
	UpdatedAt  time.Time          `json:"updatedat" bson:"updatedat"`
	Active     bool               `json:"isactive" bson:"isctive" validate:"required"`
}
