package database

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrDBNotConnected = errors.New("OOPS !!!! sorry to say but you found some issue while data base connection ")
	ErrDBNotPinnged   = errors.New("not able to connect the ping the connected database ")
)

var Client *mongo.Client = DBSetup()

// func DBSetup() *mongo.Client {

// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Panic("error while env file load")
// 	}
// 	clientoption := options.Client().ApplyURI("mongodb://development:testpassword@localhost:27017")

// 	client, err := mongo.Connect(context.Background(), clientoption)

// 	if err != nil {
// 		log.Panicln(customerrors.ErrDBNotConnected)
// 	}

// 	err = client.Ping(context.Background(), nil)

// 	if err != nil {
// 		log.Panicln(customerrors.ErrDBNotPinnged)
// 	}
// 	log.Println("connected to database")

// 	return client
// }

func DBSetup() *mongo.Client {
    err := godotenv.Load()
    if err != nil {
        log.Panic("error while loading environment variables")
    }

    // Print environment variables for debugging
    log.Println("MONGO_INITDB_ROOT_USERNAME:", os.Getenv("MONGO_INITDB_ROOT_USERNAME"))
    log.Println("MONGO_INITDB_ROOT_PASSWORD:", os.Getenv("MONGO_INITDB_ROOT_PASSWORD"))

    // Construct MongoDB connection string
    uri := "mongodb://" + os.Getenv("MONGO_INITDB_ROOT_USERNAME") + ":" + os.Getenv("MONGO_INITDB_ROOT_PASSWORD") + "@mongo:27017"
    log.Println("Connection String:", uri)

    // Set client options
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

    // Connect to MongoDB
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Panicln(ErrDBNotConnected)
    }

    // Check if connected to MongoDB
    err = client.Ping(context.Background(), nil)
    if err != nil {
        log.Panicln(ErrDBNotPinnged)
    }

    log.Println("Connected to MongoDB")
    return client
}

// fetching the user collection

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {

	collection := client.Database(os.Getenv("DB")).Collection(collectionName)
	return collection
}
