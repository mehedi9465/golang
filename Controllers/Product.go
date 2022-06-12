package Controllers

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"seentech/RECR/DBManager"
	"seentech/RECR/Models"
	"seentech/RECR/Utils"
	"seentech/RECR/Utils/Responses"
)

func isPrdocutExisting(collection *mongo.Collection, name string) (bool, interface{}) {

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

func ProductGetById(objID primitive.ObjectID) (Models.Product, error) {

	var self Models.Product

	var filter bson.M = bson.M{}
	filter = bson.M{"_id": objID}

	collection := DBManager.SystemCollections.Product

	var results []bson.M
	b, results := Utils.FindByFilter(collection, filter)
	if !b || len(results) == 0 {
		return self, errors.New("obj not found")
	}

	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode

	return self, nil
}

func ProductNew(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.Product

	// Fill the received data inside an obj
	var self Models.Product
	c.BodyParser(&self)

	// Validate the obj
	err := self.Validate()
	if err != nil {
		return Responses.NotValid(c, err.Error())
	}

	// Check if this obj is already existing
	existing, _ := isPrdocutExisting(collection, self.Name)
	if existing {
		return errors.New("product name is already existing")
	}

	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}

	Responses.Created(c, "Product", res)
	return nil
}

func ProductsGetAll(c *fiber.Ctx) error {
	TicketCol := DBManager.SystemCollections.Product
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
