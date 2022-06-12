package Controllers

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"seentech/RECR/DBManager"
	"seentech/RECR/Models"
	"seentech/RECR/Utils"
	"seentech/RECR/Utils/Responses"
)

func isChannelExisting(name string) (bool, interface{}) {
	// Initiate the connection
	collection := DBManager.SystemCollections.Channel
	var filter bson.M = bson.M{
		"name": name,
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

func ChannelNew(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.Channel

	// Fill the received data inside an obj
	image := Utils.UploadImage(c)
	var self Models.Channel
	c.BodyParser(&self)
	if image != "Error" {
		self.Icon = image
	}
	// Validate the obj
	err := self.Validate()
	if err != nil {
		return Responses.NotValid(c, err.Error())
	}

	// Check if this obj is already existing
	existing, id := isChannelExisting(self.Name)
	if existing {
		return Responses.ResourceAlreadyExist(c, "Channel", fiber.Map{"id": id})
	}

	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}

	Responses.Created(c, "Channel", res)
	return nil
}

func ChannelGet(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.Channel

	var self Models.ChannelSearch
	c.QueryParser(&self)

	var results []bson.M

	b, results := Utils.FindByFilter(collection, self.GetChannelSearchBSONObj())
	if !b {
		return Responses.NotFound(c, "Channel")
	}
	Responses.Get(c, "Channel", results)
	return nil
}

func ChannelGetById(objID primitive.ObjectID) (Models.Channel, error) {

	var self Models.Channel

	var filter bson.M = bson.M{}
	filter = bson.M{"_id": objID}

	collection := DBManager.SystemCollections.Channel

	var results []bson.M
	b, results := Utils.FindByFilter(collection, filter)
	if !b || len(results) == 0 {
		return self, errors.New("obj not found")
	}

	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode

	return self, nil
}

func ChannelSetStatus(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.Channel

	if c.Params("id") == "" || c.Params("new_status") == "" {
		return Responses.NotValid(c, "Missing required parameter [id, new_status]")
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	var newValue = true
	if c.Params("new_status") == "inactive" {
		newValue = false
	}

	updateData := bson.M{
		"$set": bson.M{
			"status": newValue,
		},
	}

	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	c.Set("Content-Type", "application/json")
	if updateErr != nil {
		return Responses.ModifiedFail(c, "Channel", updateErr.Error())
	} else {
		Responses.ModifiedSuccess(c, "Channel")
	}
	return nil
}

func ChannelModify(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.Channel
	// Fill the received data inside an obj
	var self Models.Channel
	c.BodyParser(&self)

	// Validate the obj
	err := self.Validate()
	if err != nil {
		return Responses.NotValid(c, err.Error())
	}

	updateData := bson.M{
		"$set": self.GetModifcationBSONObj(),
	}

	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": self.GetId()}, updateData)

	if updateErr != nil {
		return Responses.ModifiedFail(c, "Channel", updateErr.Error())
	} else {
		Responses.ModifiedSuccess(c, "Channel")
	}
	return nil
}
