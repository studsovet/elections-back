package db

import (
	"context"
	"errors"
)

type Vote struct {
	Data        string `bson:"data" json:"data" bindings:"required"`
	BallotBoxID int    `bson"ballotid" json:"ballotid" bindings:"required"`
}

func (v *Vote) SaveVote() (*Vote, error) {
	var err error

	coll := DB.Database("public").Collection("votes")
	_, err = coll.InsertOne(context.TODO(), v)
	if err != nil {
		return &Vote{}, err
	}

	return v, nil
}

func GetVotes() ([]Vote, error) {
	/*
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
	*/

	return []Vote{}, errors.New("GetVotes is unavalible")
}
