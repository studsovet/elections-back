package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func IsAdmin(id string) (bool, error) {
	coll := DB.Database("private").Collection("admins")
	filter := bson.D{{Key: "id", Value: id}}
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	return (count > 0), nil
}

func IsObserver(id string) (bool, error) {
	coll := DB.Database("private").Collection("observers")
	filter := bson.D{{Key: "id", Value: id}}
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	return (count > 0), nil
}
