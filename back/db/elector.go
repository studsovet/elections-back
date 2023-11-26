package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func (e *Elector) Save() (*Elector, error) {
	coll := DB.Database("private").Collection("electors")
	filter := bson.D{{Key: "id", Value: e.ID}}
	res, err := coll.ReplaceOne(context.TODO(), filter, e)
	if err != nil {
		return &Elector{}, err
	}
	if res.MatchedCount == 0 {
		_, err := coll.InsertOne(context.TODO(), e)
		if err != nil {
			return &Elector{}, err
		}
	}
	return e, nil
}

func IsElectorSaved(id string) (bool, error) {
	coll := DB.Database("private").Collection("electors")

	filter := bson.D{{Key: "id", Value: id}}
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
