package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	User_ID     primitive.ObjectID `json:"user_id"`
	Total_Price uint               `json:"total_price"`
}
type Transaction_Item struct {
	Transaction_ID primitive.ObjectID `json:"transaction_id"`
	Product_ID     primitive.ObjectID `json:"product_id"`
	Amount         uint               `json:"amount"`
}
