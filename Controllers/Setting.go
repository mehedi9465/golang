package Controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"seentech/RECR/DBManager"
	"seentech/RECR/Models"
	"seentech/RECR/Utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SettingGetAll(c *fiber.Ctx) error {
	var self Models.SettingSearch
	c.BodyParser(&self)
	results, err := settingGetAll(&self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(bson.M{"result": results})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

func settingGetAll(self *Models.SettingSearch) ([]bson.M, error) {
	collection := DBManager.SystemCollections.Setting
	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetSettingSearchBSONObj())
	if !b || len(results) <= 0 {
		return results, errors.New("no settings object found")
	}
	return results, nil
}

func SettingNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Setting
	var self Models.Setting
	c.BodyParser(&self)
	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(res)
	c.Status(200).Send(response)
	return nil
}

func InitializeSetting() bool {
	collection := DBManager.SystemCollections.Setting
	_, results := Utils.FindByFilter(collection, bson.M{})
	if len(results) <= 0 { //no settings is initialized
		var self Models.Setting
		self.ID = primitive.NewObjectID()
		self.CampaignSerial = 1
		_, err := collection.InsertOne(context.Background(), self)
		if err != nil {
			return false
		}
		fmt.Println("Initializing Setting Is Done")
	}
	return true
}

func SettingIncreaseCampaignSerial() error {
	// get setting value
	settingRes, settingErr := settingGetAll(&Models.SettingSearch{})
	if settingErr != nil {
		return settingErr
	}
	byteArray, _ := json.Marshal(settingRes[0])
	var setting Models.Setting
	json.Unmarshal(byteArray, &setting)

	// set setting value
	collectionSetting := DBManager.SystemCollections.Setting
	updateData := bson.M{
		"$set": bson.M{
			"campaignserial": setting.CampaignSerial + 1,
		},
	}
	_, updateErr := collectionSetting.UpdateOne(context.Background(), bson.M{"_id": setting.ID}, updateData)
	if updateErr != nil {
		return errors.New("an error occurred when Incrementing Serial Number")
	} else {
		return nil
	}
}

func SettingGetCampaignSerial() (int, error) {
	settingRes, settingErr := settingGetAll(&Models.SettingSearch{})
	if settingErr != nil {
		return -1, settingErr
	}

	byteArray, _ := json.Marshal(settingRes[0])
	var setting Models.Setting
	json.Unmarshal(byteArray, &setting)
	return setting.CampaignSerial, nil
}
