package Controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"seentech/RECR/DBManager"
	"seentech/RECR/Models"
	"seentech/RECR/Utils"
	"seentech/RECR/Utils/Responses"
)

func CampaignGetById(objID primitive.ObjectID) (Models.Campaign, error) {
	collection := DBManager.SystemCollections.Campaign
	var self Models.Campaign

	filter := bson.M{"_id": objID}

	var results []bson.M
	b, results := Utils.FindByFilter(collection, filter)
	if !b || len(results) == 0 {
		return self, errors.New("Campaign not found")
	}

	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode

	return self, nil
}

func CampaignNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Campaign

	var self Models.Campaign
	c.BodyParser(&self)
	self.Status = Models.StatusDraft

	err := self.Validate()
	if err != nil {
		return err
	}

	if self.Type == "Online" {
		err = self.Online.ValidateOnline()
		if err != nil {
			return err
		}

		_, err := ChannelGetById(self.Online.ChannelRef)
		if err != nil {
			return err
		}

	} else if self.Type == "Event" {
		err = self.Event.ValidateEvent()
		if err != nil {
			return err
		}

	}

	if len(self.Online.Files) > 0 {
		if len(self.Online.Files) != len(self.Online.FilesExt) {
			return err
		}

		for i, doc := range self.Online.Files {
			filePath, err := Utils.UtilsUploadFilesBase64(doc, self.Online.FilesExt[i], "file")
			if err != nil {
				c.Status(500)
				return err
			}
			if filePath != "" {
				self.Online.Files[i] = filePath
			}
		}
	} else {
		if len(self.Event.Files) != len(self.Event.FilesExt) {
			return err
		}

		for i, doc := range self.Event.Files {
			filePath, err := Utils.UtilsUploadFilesBase64(doc, self.Event.FilesExt[i], "file")
			if err != nil {
				c.Status(500)
				return err
			}
			if filePath != "" {
				self.Event.Files[i] = filePath
			}
		}
	}

	self.Status = Models.StatusDraft

	err = self.Validate()
	if err != nil {
		return err
	}

	serial, err := SettingGetCampaignSerial()
	if err != nil {
		return err
	}
	self.Serial = fmt.Sprintf("%09d", serial)

	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}
	// increase Campaign Serial
	_ = SettingIncreaseCampaignSerial()

	Responses.Created(c, "Campaign", res)
	return nil
}

func CampaignDelete(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Campaign

	campaignID, err := primitive.ObjectIDFromHex(c.Params("campaignID"))
	if err != nil {
		return err
	}

	self, err := CampaignGetById(campaignID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": campaignID,
	}

	if self.Status == Models.StatusDraft || self.Status == Models.StatusReview {
		res, err := collection.DeleteOne(context.Background(), filter)
		if err != nil {
			return Responses.BadRequest(c, err.Error())
		}
		fmt.Println(res)
		Responses.DeletedSuccess(c, "Campaign")
	} else {
		return Responses.BadRequest(c, err.Error())
	}
	return nil
}

func CampaignStatusModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Campaign

	if c.Params("campaignID") == "" || c.Params("status") == "" {
		return Responses.Empty(c, "Not Valid Parameters")
	}

	campaignID, err := primitive.ObjectIDFromHex(c.Params("campaignID"))
	if err != nil {
		return err
	}
	filter := bson.M{"_id": campaignID}

	self, err := CampaignGetById(campaignID)
	if err != nil {
		return err
	}
	campaignStatus := strings.Title(c.Params("status"))
	if err != nil {
		return err
	}

	var updatedData primitive.M

	switch self.Status {
	case "Draft":
		if campaignStatus == Models.StatusDraft {
			return Responses.StatusUnchanged(c, "Already in "+campaignStatus)
		} else if campaignStatus == Models.StatusReview {

			if self.OwnerRef == primitive.NilObjectID {
				return Responses.NotFound(c, "OwnerRef")
			}
			updatedData = bson.M{
				"$set": bson.M{
					"status":      campaignStatus,
					"reviewerref": self.OwnerRef, //TODO : add user ID
				},
			}
		} else {
			return Responses.StatusChangeFail(c)
		}

	case "Review":
		if campaignStatus == Models.StatusReview {
			return Responses.StatusUnchanged(c, "Already in "+campaignStatus)
		} else if campaignStatus == Models.StatusConfirmed || campaignStatus == Models.StatusDraft {
			updatedData = bson.M{
				"$set": bson.M{
					"status": campaignStatus,
				},
			}
		} else {
			return Responses.StatusChangeFail(c)
		}

	case "Confirmed":
		if campaignStatus == Models.StatusConfirmed {
			return Responses.StatusUnchanged(c, "Already in "+campaignStatus)
		} else if campaignStatus == Models.StatusStarted || campaignStatus == Models.StatusCanceled {
			updatedData = bson.M{
				"$set": bson.M{
					"status": campaignStatus,
				},
			}
		} else {
			return Responses.StatusChangeFail(c)
		}

	case "Started":
		if campaignStatus == Models.StatusStarted {
			return Responses.StatusUnchanged(c, "Already in "+campaignStatus)
		} else if campaignStatus == Models.StatusClosed || campaignStatus == Models.StatusCanceled {
			updatedData = bson.M{
				"$set": bson.M{
					"status": campaignStatus,
				},
			}
		} else {
			return Responses.StatusChangeFail(c)
		}

	case "Closed":
		if campaignStatus == Models.StatusClosed {
			return Responses.StatusUnchanged(c, "Already in "+campaignStatus)
		}
		return Responses.StatusChangeFail(c)

	case "Canceled":
		if campaignStatus == Models.StatusCanceled {
			return Responses.StatusUnchanged(c, "Already in "+campaignStatus)
		}
		return Responses.StatusChangeFail(c)

	default:
		return Responses.StatusChangeFail(c)
	}

	_, err = collection.UpdateOne(context.Background(), filter, updatedData)
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}

	Responses.ModifiedSuccess(c, "Campaign Status")
	return nil
}

func CampaignModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Campaign

	campaignID, err := primitive.ObjectIDFromHex(c.Params("campaignID"))
	if err != nil {
		return err
	}

	var self Models.Campaign
	c.BodyParser(&self)
	er := self.Validate()
	if er != nil {
		c.Status(500)
		return er
	}

	if campaignID == primitive.NilObjectID {
		return Responses.Empty(c, "Not Valid Parameters")
	}

	if len(self.Online.Files) > 0 {
		if len(self.Online.Files) != len(self.Online.FilesExt) {
			return err
		}

		for i, doc := range self.Online.Files {
			filePath, err := Utils.UtilsUploadFilesBase64(doc, self.Online.FilesExt[i], "file")
			if err != nil {
				c.Status(500)
				return err
			}
			if filePath != "" {
				self.Online.Files[i] = filePath
			}
		}
	} else {
		if len(self.Event.Files) != len(self.Event.FilesExt) {
			return err
		}

		for i, doc := range self.Event.Files {
			filePath, err := Utils.UtilsUploadFilesBase64(doc, self.Event.FilesExt[i], "file")
			if err != nil {
				c.Status(500)
				return err
			}
			if filePath != "" {
				self.Event.Files[i] = filePath
			}
		}
	}

	updateData := bson.M{
		"$set": self.GetCampaignModificationBSONObj(),
	}

	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": campaignID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return Responses.BadRequest(c, err.Error())
	}
	Responses.ModifiedSuccess(c, "Campaign")
	return nil
}

func campaignGetAll(self *Models.CampaignSearch) ([]bson.M, error) {
	collection := DBManager.SystemCollections.Campaign
	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetCampaignSearchBSONObj())
	if !b {
		return results, errors.New("Campaign object found")
	}
	return results, nil
}

func CampaignGetAll(c *fiber.Ctx) error {
	var self Models.CampaignSearch
	c.BodyParser(&self)
	results, err := campaignGetAll(&self)
	if err != nil {
		return Responses.NotFound(c, err.Error())
	}
	response, _ := json.Marshal(bson.M{"result": results})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

func CampaignGetAllPopulated(c *fiber.Ctx) error {
	CampaignCollection := DBManager.SystemCollections.Campaign

	var self Models.CampaignSearch
	c.BodyParser(&self)

	b, results := Utils.FindByFilter(CampaignCollection, self.GetCampaignSearchBSONObj())
	if !b {
		c.Status(500)
		return Responses.NotFound(c, "Campaign")
	}

	byteArr, _ := json.Marshal(results)
	var ResultDocs []Models.Campaign
	json.Unmarshal(byteArr, &ResultDocs)

	populatedResult := make([]Models.CampaignPopulated, len(ResultDocs))

	for i, v := range ResultDocs {
		populatedResult[i], _ = CampaignGetByIdPopulated(v.ID, &v)
	}
	allpopulated, _ := json.Marshal(bson.M{"results": populatedResult})
	c.Set("Content-Type", "application/json")
	c.Send(allpopulated)
	return nil
}

func CampaignGetByIdPopulated(objID primitive.ObjectID, ptr *Models.Campaign) (Models.CampaignPopulated, error) {
	var CampaignDoc Models.Campaign
	if ptr == nil {
		CampaignDoc, _ = CampaignGetById(objID)
	} else {
		CampaignDoc = *ptr
	}
	populatedResult := Models.CampaignPopulated{}
	populatedResult.CloneFrom(CampaignDoc)

	populatedResult.OwnerRef, _ = UserGetById(CampaignDoc.OwnerRef)

	populatedResult.ReviewerRef, _ = UserGetById(CampaignDoc.ReviewerRef)

	populatedResult.Products = make([]Models.Product, len(CampaignDoc.Products))
	for i, productID := range CampaignDoc.Products {
		populatedResult.Products[i], _ = ProductGetById(productID)
	}

	populatedResult.Team = make([]Models.User, len(CampaignDoc.Team))
	for i, userID := range CampaignDoc.Team {
		populatedResult.Team[i], _ = UserGetById(userID)
	}

	populatedResult.Online.CloneFrom(CampaignDoc.Online)
	populatedResult.Online.ChannelRef, _ = ChannelGetById(CampaignDoc.Online.ChannelRef)
	populatedResult.Online.Files = CampaignDoc.Online.Files
	populatedResult.Online.Others = CampaignDoc.Online.Others

	return populatedResult, nil
}
