package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type Vote struct {
	Vote   string `bson:"vote" json:"vote" bindings:"required"`
	UserId string `bson"user_id" json:"user_id"`
}

func (v *Vote) SaveVote() (*Vote, error) {
	_ = v.UserId
	var err error
	coll := DB.Database("public").Collection("votes")
	_, err = coll.InsertOne(context.TODO(), v)
	if err != nil {
		return &Vote{}, err
	}
	return v, nil
}

func GetVotes() ([]Vote, error) {
	coll := DB.Database("public").Collection("votes")
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return []Vote{}, err
	}
	var votes []Vote
	err = cursor.All(context.TODO(), &votes)
	if err != nil {
		return []Vote{}, err
	}
	return votes, nil
}
