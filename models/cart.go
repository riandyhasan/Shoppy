package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Add_Products struct {
	Product_ID primitive.ObjectID `json:"product_id"`
	Amount     uint               `json:"amount"`
}

type Cart_Products struct {
	User_ID primitive.ObjectID `json:"user_id"`
	Product Product            `json:"product"`
	Amount  uint               `json:"amount"`
}
