package Models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ticket struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" binding:"required"`
	UserRef    primitive.ObjectID `json:"userref,omitempty" bson:"userref,omitempty"`
	ProductRef primitive.ObjectID `json:"productref,omitempty" bson:"productref,omitempty"`
	IsClosed   bool               `json:"isclosed"`
	Events     []Event            `json:"events,omitempty" bson:"events,omitempty"`
}

func (obj Ticket) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required))
}

func (obj Ticket) GetModifcationBSONObj() bson.M {
	self := bson.M{
		"name":     obj.Name,
		"isclosed": obj.IsClosed,
		"events":   obj.Events,
	}
	// TODO CHECK THIS (user ref and product ref is inserted zeros even they omitempty)
	if obj.UserRef.IsZero() == false {
		self["userref"] = obj.UserRef
	}
	if obj.ProductRef.IsZero() == false {
		self["productref"] = obj.ProductRef
	}
	return self
}

type TicketPopulated struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name,omitempty" binding:"required"`
	User    *User              `json:"user,omitempty" bson:"user,omitempty"`
	Product *Product           `json:"product,omitempty" bson:"product,omitempty"`
	Events  []Event            `json:"events,omitempty" bson:"events,omitempty"`
}

// Using User and product as pointers to be easy to marshal (withoud will have this User{} , Product{})
// Note we should do this with bool too(3 states : nil , true , false)
func (obj *TicketPopulated) CloneFrom(other Ticket) {
	obj.ID = other.ID
	obj.Name = other.Name
	obj.User = nil
	obj.Product = nil
	obj.Events = other.Events
}

// possible events that can happen to the ticket
const (
	CreateTicket string = "CreateTicket"
	AssignTicket        = "AssignTicket"
	CloseTicket         = "CloseTicket"
)

// Event struct which will represent our log messages that happens to the ticket
// There is a problem with ObjectID here https://jira.mongodb.org/browse/GODRIVER-1332 so we will use pointers
type Event struct {
	// Event Type which can be (Assign/CloseTicket/Moveing Forward/Moving Backward/Created/Call/Meeting/Call Reminder/Meeting Reminder)
	Type       string              `json:"type"`
	UserBy     primitive.ObjectID  `json:"userby,omitempty" bson:"userby,omitempty"`
	NewUserRef *primitive.ObjectID `json:"newuserref,omitempty" bson:"newuserref,omitempty"`
	OldUserRef *primitive.ObjectID `json:"olduserref,omitempty" bson:"olduserref,omitempty"`
	Date       primitive.DateTime  `json:"date"`
	Reason     string              `json:"reason,omitempty" bson:"reason,omitempty"`
	Deal       bool                `json:"deal" bson:"deal"`
}

func (obj Event) Validate() error {
	err := validation.ValidateStruct(&obj,
		validation.Field(&obj.Type, validation.In(CreateTicket, AssignTicket, CloseTicket)),
		validation.Field(&obj.Date, validation.Required),
		validation.Field(&obj.UserBy, validation.Required))
	if err != nil {
		return err
	}

	switch obj.Type {
	case AssignTicket:
		return validation.ValidateStruct(&obj,
			validation.Field(&obj.NewUserRef, validation.Required),
			validation.Field(&obj.OldUserRef, validation.Required))
	case CloseTicket:
		return validation.ValidateStruct(&obj,
			validation.Field(&obj.Reason, validation.Required))
	}

	return nil
}
