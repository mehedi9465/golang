package DBManager

import (
	"context"
	"log"

	// get an object type
	// "encoding/json"

	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var config_err = godotenv.Load()
var dbURL string = os.Getenv("DB_SOURCE_URL")
var SystemCollections CMSCollections

type CMSCollections struct {
	User     *mongo.Collection
	Product  *mongo.Collection
	Ticket   *mongo.Collection
	Channel  *mongo.Collection
	Campaign *mongo.Collection
	Setting  *mongo.Collection
}

func InitCMSCollections() bool {
	if config_err != nil {
		return false
	}
	var err error
	SystemCollections.User, err = GetMongoDbCollection(os.Getenv("DB_NAME"), "user")
	if err != nil {
		return false
	}
	SystemCollections.Product, err = GetMongoDbCollection(os.Getenv("DB_NAME"), "product")
	if err != nil {
		return false
	}
	SystemCollections.Ticket, err = GetMongoDbCollection(os.Getenv("DB_NAME"), "ticket")
	if err != nil {
		return false
	}
	SystemCollections.Channel, err = GetMongoDbCollection(os.Getenv("DB_NAME"), "channel")
	if err != nil {
		return false
	}
	SystemCollections.Campaign, err = GetMongoDbCollection(os.Getenv("DB_NAME"), "campaign")
	if err != nil {
		return false
	}
	SystemCollections.Setting, err = GetMongoDbCollection(os.Getenv("DB_NAME"), "setting")
	if err != nil {
		return false
	}

	return err == nil
}

// GetMongoDbConnection get connection of mongodb
func getMongoDbConnection() (*mongo.Client, error) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dbURL))

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}

func GetMongoDbCollection(DbName string, CollectionName string) (*mongo.Collection, error) {
	client, err := getMongoDbConnection()

	if err != nil {
		return nil, err
	}

	collection := client.Database(DbName).Collection(CollectionName)

	return collection, nil
}
