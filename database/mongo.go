package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client
var UserCollection *mongo.Collection
var VaultCollection *mongo.Collection

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI is not set in the environment variables")
	}

	clientOptions := options.Client().ApplyURI(uri)

	Client, err = mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	db := Client.Database("password_manager")
	UserCollection = db.Collection("users")
	log.Println("connection to mongoDB successful")

	VaultCollection = db.Collection("vault")

}
