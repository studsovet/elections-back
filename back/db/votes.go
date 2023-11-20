package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func (v *EncryptedVote) Save(election_id string) (*EncryptedVote, error) {
	var err error

	coll := DB.Database("public").Collection("encrypted_votes_election_" + election_id)
	_, err = coll.InsertOne(context.TODO(), v)
	if err != nil {
		return &EncryptedVote{}, err
	}
	return v, nil
}

func (v *DecryptedVote) Save(election_id string) (*DecryptedVote, error) {
	var err error

	coll := DB.Database("public").Collection("decrypted_votes_election_" + election_id)
	_, err = coll.InsertOne(context.TODO(), v)
	if err != nil {
		return &DecryptedVote{}, err
	}
	return v, nil
}

func DropEncryptedVotes(election_id string) error {
	return DB.Database("public").Collection("decrypted_votes_election_" + election_id).Drop(context.TODO())
}

func IsVoted(election_id string, voter_id string) (bool, error) {
	coll := DB.Database("public").Collection("alreadyVoted")

	filter := bson.D{{Key: "voterId", Value: voter_id}, {Key: "electionId", Value: election_id}}
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func SetVoted(election_id string, voter_id string) {
	coll := DB.Database("public").Collection("alreadyVoted")

	elem := bson.D{{Key: "voterId", Value: voter_id}, {Key: "electionId", Value: election_id}}
	coll.InsertOne(context.TODO(), elem)
}

func GetEncryptedVotes(election_id string) ([]EncryptedVote, error) {
	coll := DB.Database("public").Collection("encrypted_votes_election_" + election_id)
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return []EncryptedVote{}, err
	}
	var votes []EncryptedVote
	err = cursor.All(context.TODO(), &votes)
	if err != nil {
		return []EncryptedVote{}, err
	}
	if votes == nil {
		return []EncryptedVote{}, nil
	}
	return votes, nil
}
