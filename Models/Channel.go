package Models

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Channel struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name,omitempty"`
	Status bool               `json:"status"`
	Link   string             `json:"link"`
	Icon   string             `json:"icon"`
}

func (obj Channel) Validate() error {
	nameErr := validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
	)
	if nameErr != nil {
		return nameErr
	}
	return nil
}

func (obj Channel) GetId() primitive.ObjectID {
	return obj.ID
}

func (obj Channel) GetModifcationBSONObj() bson.M {
	self := bson.M{
		"_id":    obj.ID,
		"name":   obj.Name,
		"status": obj.Status,
		"link":   obj.Link,
		"icon":   obj.Icon,
	}

	return self
}

type ChannelSearch struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	IDIsUsed       bool               `json:"idisused"`
	Name           string             `json:"name,omitempty"`
	NameIsUsed     bool               `json:"nameisused,omitempty"`
	Status         bool               `json:"status,omitempty"`
	StatusIsUsed   bool               `json:"statusisused,omitempty"`
	Link           string             `json:"link,omitempty"`
	LinkIsUsed     bool               `json:"linkisused,omitempty"`
	Icon           string             `json:"icon,omitempty"`
	IconIsUsed     bool               `json:"iconisused,omitempty"`
	TextData       string             `json:"textdata,omitempty"`
	IsTextDataUsed bool               `json:"istextdataused,omitempty"`
}

func (obj ChannelSearch) GetChannelSearchBSONObj() bson.M {
	self := bson.M{}

	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}

	if obj.NameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Name)
		self["name"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.StatusIsUsed {
		self["status"] = obj.Status
	}

	if obj.LinkIsUsed {
		self["link"] = obj.Link
	}

	if obj.IconIsUsed {
		self["icon"] = obj.Icon
	}

	if obj.IsTextDataUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.TextData)

		self["$or"] = []bson.M{
			{"name": bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}},
			{"link": bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}},
			{"icon": bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}},
		}
	}

	return self
}
