package models

import (
	"shoppy/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Name     string             `json:"name"`
	Category types.Category     `json:"category"`
	Price    uint               `json:"price"`
}
