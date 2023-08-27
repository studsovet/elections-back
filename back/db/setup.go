package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func ConnectDB() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		fmt.Println("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	DB = db
	if err != nil {
		panic(err)
	}
}

func DisconnectDB() {
	if err := DB.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
