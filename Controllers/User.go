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

func isUserExisting(collection *mongo.Collection, name string) (bool, interface{}) {

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

func UserGetById(objID primitive.ObjectID) (Models.User, error) {

	var self Models.User

	var filter bson.M = bson.M{}
	filter = bson.M{"_id": objID}

	collection := DBManager.SystemCollections.User

	var results []bson.M
	b, results := Utils.FindByFilter(collection, filter)
	if !b || len(results) == 0 {
		return self, errors.New("obj not found")
	}

	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode

	return self, nil
}

func UserNew(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.User

	// Fill the received data inside an obj
	var self Models.User
	c.BodyParser(&self)

	// Validate the obj
	err := self.Validate()
	if err != nil {
		return Responses.NotValid(c, err.Error())
	}

	// Check if this obj is already existing
	existing, _ := isUserExisting(collection, self.Name)
	if existing {
		return errors.New("user name is already existing")
	}

	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}

	Responses.Created(c, "User", res)
	return nil
}

func userGetAll(self *Models.UserSearch) ([]bson.M, error) {
	collection := DBManager.SystemCollections.User
	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetUserSearchBSONObj())
	if !b {
		return results, errors.New("No User found")
	}
	return results, nil
}

func UserGetAll(c *fiber.Ctx) error {
	var self Models.UserSearch
	c.BodyParser(&self)
	results, err := userGetAll(&self)
	if err != nil {
		return Responses.NotFound(c, err.Error())
	}
	response, _ := json.Marshal(bson.M{"result": results})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}
