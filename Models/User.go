package Models

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" binding:"required"`
}

func (obj User) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required))

}

type UserSearch struct {
	IDIsUsed   bool               `json:"idisused,omitempty" bson:"idisused,omitempty"`
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	NameIsUsed bool               `json:"nameisused,omitempty" binding:"required"`
	Name       string             `json:"name,omitempty" binding:"required"`
}

func (obj UserSearch) GetUserSearchBSONObj() bson.M {
	self := bson.M{}

	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}

	if obj.NameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Name)
		self["name"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	return self
}
