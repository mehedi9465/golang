package Utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"mime/multipart"
	"os"
	"seentech/RECR/Utils/Responses"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CheckIfObjExistingByObjId(collection *mongo.Collection, objID primitive.ObjectID) error {
	filter := bson.M{"_id": objID}

	var results []bson.M
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())

	cur.All(context.Background(), &results)
	fmt.Println("Count : ", len(results))

	if len(results) == 0 {
		return errors.New("obj not found")
	}

	return nil
}

func AdaptCurrentTimeByUnit(unit string, period int) time.Time {
	now := time.Now()
	if unit == "Month" {
		now = now.AddDate(0, period, 0)
	} else if unit == "Week" {
		now = now.AddDate(0, 0, period*7)
	} else if unit == "Day" {
		now = now.AddDate(0, 0, period)
	} else if unit == "Year" {
		now = now.AddDate(period, 0, 0)
	}
	return now
}

func AdaptRefernceTimeByUnit(refernceTime time.Time, unit string, period int) time.Time {
	if unit == "Month" {
		refernceTime = refernceTime.AddDate(0, period, 0)
	} else if unit == "Week" {
		refernceTime = refernceTime.AddDate(0, 0, period*7)
	} else if unit == "Day" {
		refernceTime = refernceTime.AddDate(0, 0, period)
	} else if unit == "Year" {
		refernceTime = refernceTime.AddDate(period, 0, 0)
	}
	return refernceTime
}

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func UploadImage(c *fiber.Ctx) string {
	file, err := c.FormFile("image")
	if err != nil {
		//fmt.Println("Failed in saving Image")
		//c.Status(500).Send([]byte("Invalid data sent for uploading"))
		return "Error"
	}

	// Save file to root directory
	var filePath = fmt.Sprintf("Resources/Images/img_%d_%d.png", rand.Intn(1024), MakeTimestamp())
	saveing_err := c.SaveFile(file, "./"+filePath)
	if saveing_err != nil {
		//c.Status(500).Send([]byte("Failed to save the uploaded image"))
		return "Error"
	} else {
		//c.Status(200).Send([]byte("Saved Successfully"))
		return filePath
	}
}

func UploadFiles(files map[string][]*multipart.FileHeader, c *fiber.Ctx) ([]string, error) {

	var filePaths = []string{}
	// Loop through files:
	for _, file := range files {
		var filePath = fmt.Sprintf("Resources/Files/img_%d_%d.png", rand.Intn(1024), MakeTimestamp())
		saveing_err := c.SaveFile(file[0], "./"+filePath)
		if saveing_err != nil {
			c.Status(500).Send([]byte("Failed to save the uploaded image"))
			return nil, saveing_err
		} else {
			filePaths = append(filePaths, filePath)
		}
	}
	return filePaths, nil
}

func FindByFilter(collection *mongo.Collection, filter bson.M) (bool, []bson.M) {
	var results []bson.M

	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return false, results
	}
	defer cur.Close(context.Background())

	cur.All(context.Background(), &results)

	return true, results
}

func FindByFilterSorted(collection *mongo.Collection, filter bson.M, sort bson.M) (bool, []bson.M) {
	var results []bson.M
	findOptions := options.Find()
	findOptions.SetSort(sort)

	cur, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return false, results
	}
	defer cur.Close(context.Background())

	cur.All(context.Background(), &results)

	return true, results
}

func Contains(arr []primitive.ObjectID, elem primitive.ObjectID) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

func ArrayStringContains(arr []string, elem string) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

func DecodeArrData(inStructArr, outStructArr interface{}) error {
	in := struct{ Data interface{} }{Data: inStructArr}
	inStructArrData, err := bson.Marshal(in)
	if err != nil {
		return err
	}
	var out struct{ Data bson.Raw }
	if err := bson.Unmarshal(inStructArrData, &out); err != nil {
		return err
	}
	return bson.Unmarshal(out.Data, &outStructArr)
}

func SendTextResponseAsJSON(c *fiber.Ctx, msg string) {
	response, _ := json.Marshal(
		bson.M{"result": msg},
	)
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
}

func CollectionProjected(c *fiber.Ctx, collection *mongo.Collection, fields bson.M) error {
	var results []bson.M
	opts := options.FindOptions{Projection: fields}
	cur, err := collection.Find(context.Background(), bson.M{}, &opts)
	if err != nil {
		return Responses.NotFound(c, "Data not found")
	}
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results)
	Responses.Get(c, "Data", results)
	return nil
}

func FindByFilterProjected(collection *mongo.Collection, filter bson.M, fields bson.M) (bool, []bson.M) {
	var results []bson.M
	opts := options.FindOptions{Projection: fields}
	cur, err := collection.Find(context.Background(), filter, &opts)
	if err != nil {
		return false, results
	}
	defer cur.Close(context.Background())

	cur.All(context.Background(), &results)

	return true, results
}

func UtilsUploadFilesBase64(stringBase64 string, filesDocType string, filesName string) (string, error) {
	i := strings.Index(stringBase64, ",")

	if i != -1 {
		file, _ := base64.StdEncoding.DecodeString(stringBase64[i+1:])
		var filePath = fmt.Sprintf("Public/Files/"+filesName+"_files_%d_%d.%s", rand.Intn(1024), MakeTimestamp(), filesDocType)

		f, err := os.Create("./" + filePath)
		if err != nil {
			return "", err
		}

		defer f.Close()

		if _, err := f.Write(file); err != nil {
			return "", err
		}

		f.Sync()

		return filePath, nil
	}

	return "", nil
}
