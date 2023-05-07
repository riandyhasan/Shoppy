package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Email    string             `json:"email" validate:"required"`
	Password string             `json:"password" validate:"required"`
	Name     string             `json:"name,omitempty" validate:"required"`
}
