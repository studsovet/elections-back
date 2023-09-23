package db

import (
	"context"
	"crypto"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type Vote struct {
	Data        string `bson:"data" json:"data" bindings:"required"`
	BallotBoxID int    `bson:"ballotid" json:"ballotid" bindings:"required"`
}

type DecryptedVote struct {
	Data        map[string]interface{} `bson:"data" json:"data" bindings:"required"`
	BallotBoxID int                    `bson:"ballotid" json:"ballotid" bindings:"required"`
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

func DecodeVotes() {
	SetStatus(2, "Votes are being decoded now!")

	cursor, err := DB.Database("public").Collection("votes").Find(context.TODO(), bson.D{})
	if err != nil {
		SetStatus(-1, "Votes can't be decoded! Error: "+err.Error())
		return
	}

	privateKey, err := GetParsedPrivateKey()
	if err != nil {
		SetStatus(-1, "Votes can't be decoded! Error: "+err.Error())
		return
	}

	for cursor.Next(context.TODO()) {
		var tmp map[string]interface{}

		var encVote Vote
		cursor.Decode(&encVote)

		encByteData, err := hex.DecodeString(encVote.Data)
		if err != nil {
			log.Println("Error in decripting. Error: " + err.Error())
			continue
		}

		data, err := privateKey.Decrypt(nil, encByteData, &rsa.OAEPOptions{Hash: crypto.SHA256})
		if err != nil {
			log.Println("Error in decripting. Error: " + err.Error())
			continue
		}

		log.Println(string(data))

		json.Unmarshal(data, &tmp)

		log.Println(tmp)

		decVote := DecryptedVote{
			Data:        tmp,
			BallotBoxID: encVote.BallotBoxID,
		}

		DB.Database("public").Collection("decoded_votes").InsertOne(context.Background(), decVote)
	}

	SetStatus(3, "Votes are decoded!")
}

func GetVotes() ([]DecryptedVote, error) {
	status := GetLastStatus()

	if status.Code != 3 {
		return []DecryptedVote{}, errors.New("Votes are encoded!")
	}

	coll := DB.Database("public").Collection("decoded_votes")
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return []DecryptedVote{}, err
	}
	var votes []DecryptedVote
	err = cursor.All(context.TODO(), &votes)
	if err != nil {
		return []DecryptedVote{}, err
	}
	return votes, nil
}

func ClearVotes() {
	DB.Database("public").Collection("votes").Drop(context.TODO())
	DB.Database("public").Collection("decoded_votes").Drop(context.TODO())
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
