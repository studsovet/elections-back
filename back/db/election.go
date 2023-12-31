package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func (e *Election) Save() (*Election, error) {
	_, err := DB.Database("public").Collection("elections").DeleteMany(context.TODO(), bson.D{{Key: "id", Value: e.ID}})

	if err != nil {
		return nil, err
	}

	_, err = DB.Database("public").Collection("elections").InsertOne(context.TODO(), e)

	return e, err
}

func GetElection(id string) (Election, error) {
	coll := DB.Database("public").Collection("elections")
	filter := bson.D{{Key: "id", Value: id}}

	var election Election
	err := coll.FindOne(context.TODO(), filter).Decode(&election)

	if err != nil {
		return Election{}, err
	}
	return election, nil
}

func GetElections(status string) ([]Election, error) {
	coll := DB.Database("public").Collection("elections")
	var filter bson.D
	if status == "all" {
		filter = bson.D{{}}
	} else {
		filter = bson.D{{Key: "status", Value: status}}
	}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return []Election{}, err
	}

	var elections []Election
	err = cursor.All(context.TODO(), &elections)
	if err != nil {
		return []Election{}, err
	}
	if elections == nil {
		return []Election{}, nil
	}
	return elections, nil
}

func ElectionUpdateStatus(id string, new_status string) error {
	coll := DB.Database("public").Collection("elections")
	filter := bson.D{{Key: "id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: new_status}}}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	return err
}
