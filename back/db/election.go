package db

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

func (e *Election) Save() (*Election, error) {
	_, err := DB.Database("public").Collection("elections").InsertOne(context.TODO(), e)

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

func ElectionNext(id string) error {
	election, err := GetElection(id)
	if err != nil {
		return err
	}
	status := election.Status
	status_num := -1
	for i, s := range Statuses {
		if s == status {
			status_num = i
			break
		}
	}
	if status_num == -1 {
		return errors.New("not such status: " + status)
	}
	if status_num+1 == len(Statuses) {
		return errors.New("Cannot move next last status: " + status)
	}
	new_status := Statuses[status_num+1]
	coll := DB.Database("public").Collection("elections")
	filter := bson.D{{Key: "id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: new_status}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	return err
}
