package db

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

type Vote struct {
	Data        string `bson:"data" json:"data" bindings:"required"`
	BallotBoxID int    `bson:"ballotid" json:"ballotid" bindings:"required"`
}

type VoteUserInfo struct {
	UserId      string `bson:"userid" json:"userid" bindings:"required"`
	BallotBoxID int    `bson:"ballotid" json:"ballotid" bindings:"required"`
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

func ClearVotes() {
	DB.Database("public").Collection("votes").Drop(context.TODO())
}

func AddUserVote(userId string, ballotBox int) error {
	var data VoteUserInfo

	data.UserId = userId
	data.BallotBoxID = ballotBox

	coll := DB.Database("public").Collection("users_votes")
	_, err := coll.InsertOne(context.TODO(), data)

	return err
}

func CheckUserVotes(userId string) ([]int, error) {
	coll := DB.Database("public").Collection("users_votes")

	cursor, err := coll.Find(context.TODO(), bson.D{{"userid", userId}})

	if err != nil {
		return []int{}, err
	}

	var results []VoteUserInfo
	if err = cursor.All(context.TODO(), &results); err != nil {
		return []int{}, err
	}

	res := []int{}

	for _, info := range results {
		res = append(res, info.BallotBoxID)
	}

	return res, nil
}

func IsUserVoted(userId string, ballotId int) (bool, error) {
	coll := DB.Database("public").Collection("users_votes")

	n, err := coll.CountDocuments(context.TODO(), bson.D{{"userid", userId}, {"ballotid", ballotId}})

	return n > 0, err
}

func ClearUserVotes() {
	DB.Database("public").Collection("users_votes").Drop(context.TODO())
}
