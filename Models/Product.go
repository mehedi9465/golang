package Models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name,omitempty" binding:"required"`
	Price float64            `json:"price,omitempty" binding:"required"`
}

func (obj Product) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
		validation.Field(&obj.Price, validation.Required))
}
