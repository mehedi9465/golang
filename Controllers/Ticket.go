package Controllers

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"seentech/RECR/DBManager"
	"seentech/RECR/Models"
	"seentech/RECR/Utils"
	"seentech/RECR/Utils/Responses"
)

func isTicketExisting(collection *mongo.Collection, name string) (bool, interface{}) {

	var filter bson.M = bson.M{
		"$or": []bson.M{
			{"name": name},
		},
	}
	var results []bson.M

	b, results := Utils.FindByFilter(collection, filter)
	b = (b && len(results) > 0)
	id := ""
	if b {
		id = results[0]["_id"].(primitive.ObjectID).Hex()
	}
	return b, id
}

func TicketGetById(objID primitive.ObjectID) (Models.Ticket, error) {

	var self Models.Ticket

	var filter bson.M = bson.M{}
	filter = bson.M{"_id": objID}

	collection := DBManager.SystemCollections.Ticket

	var results []bson.M
	b, results := Utils.FindByFilter(collection, filter)
	if !b || len(results) == 0 {
		return self, errors.New("obj not found")
	}

	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode

	return self, nil
}

func TicketNew(c *fiber.Ctx) error {
	// TODO use sessions and check if the user exists and authorized
	userID, err := primitive.ObjectIDFromHex("623832d6c29e056d7376cd9c")
	if err != nil {
		return err
	}

	// Initiate the connection
	collection := DBManager.SystemCollections.Ticket

	// Fill the received data inside an obj
	var self Models.Ticket
	c.BodyParser(&self)
	self.IsClosed = false

	// Validate the obj
	err = self.Validate()
	if err != nil {
		return Responses.NotValid(c, err.Error())
	}

	// Check if this obj is already existing
	existing, _ := isTicketExisting(collection, self.Name)
	if existing {
		return errors.New("ticket name is already existing")
	}

	// create the insertion event
	var event Models.Event
	event.Type = Models.CreateTicket
	event.Date = primitive.NewDateTimeFromTime(time.Now())
	event.UserBy = userID

	err = event.Validate()
	if err != nil {
		return Responses.NotValid(c, err.Error())
	}

	self.Events = []Models.Event{event}

	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}

	Responses.Created(c, "Ticket", res)
	return nil
}

func TicketClose(c *fiber.Ctx) error {
	// TODO use sessions and check if the user exists and authorized
	userID, err := primitive.ObjectIDFromHex("623832d6c29e056d7376cd9c")
	if err != nil {
		return err
	}

	// Initiate the connection
	collection := DBManager.SystemCollections.Ticket

	ticketID, err := primitive.ObjectIDFromHex(c.Params("ticketid"))
	if err != nil {
		return err
	}

	// Check if this obj exists
	ticket, existingError := TicketGetById(ticketID)
	if existingError != nil {
		return existingError
	}

	ticket.IsClosed = true

	// create the insertion event
	var event Models.Event
	event.Type = Models.CloseTicket
	event.Date = primitive.NewDateTimeFromTime(time.Now())
	event.UserBy = userID
	event.Reason = c.Params("reason")
	event.Deal, err = strconv.ParseBool(c.Params("deal"))
	if err != nil {
		return err
	}

	err = event.Validate()
	if err != nil {
		return Responses.NotValid(c, err.Error())
	}

	ticket.Events = append(ticket.Events, event)

	updateData := bson.M{
		"$set": ticket.GetModifcationBSONObj(),
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": ticketID}, updateData)
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}

	Responses.ModifiedSuccess(c, "Ticket")
	return nil
}

func TicketsGetAll(c *fiber.Ctx) error {
	TicketCol := DBManager.SystemCollections.Ticket
	cur, err := TicketCol.Find(context.Background(), bson.M{})
	defer cur.Close(context.Background())
	if err != nil {
		return Responses.NotFound(c, "Data not found")
	}

	results := []bson.M{}
	cur.All(context.Background(), &results)

	response, _ := json.Marshal(bson.M{"result": results})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

func TicketsGetAllPopulated(c *fiber.Ctx) error {
	TicketCol := DBManager.SystemCollections.Ticket
	cur, err := TicketCol.Find(context.Background(), bson.M{})
	defer cur.Close(context.Background())
	if err != nil {
		return Responses.NotFound(c, "Data not found")
	}

	var tickets []Models.Ticket
	cur.All(context.Background(), &tickets)

	populatedResult := make([]Models.TicketPopulated, len(tickets))
	for indx, value := range tickets {
		populatedResult[indx].CloneFrom(value)

		user, err := UserGetById(value.UserRef)
		if err == nil {
			populatedResult[indx].User = &user
		}

		product, err := ProductGetById(value.ProductRef)
		if err == nil {
			populatedResult[indx].Product = &product
		}
	}

	response, _ := json.Marshal(bson.M{"result": populatedResult})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}
