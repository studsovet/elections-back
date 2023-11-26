package db

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Email2ID(email string) (string, error) {
	coll := DB.Database("private").Collection("IDs")
	filter := bson.D{{Key: "email", Value: email}}
	var user_id UserId
	err := coll.FindOne(context.TODO(), filter).Decode(&user_id)
	if err == mongo.ErrNoDocuments {
		id := uuid.NewString()
		coll.InsertOne(context.TODO(), UserId{email, id})
	} else if err != nil {
		return "", err
	}
	return user_id.ID, nil
}

func ID2Email(id string) (string, error) {
	coll := DB.Database("private").Collection("IDs")
	filter := bson.D{{Key: "id", Value: id}}
	var user_id UserId
	err := coll.FindOne(context.TODO(), filter).Decode(&user_id)
	if err != nil {
		return "", err
	}
	return user_id.Email, nil
}
