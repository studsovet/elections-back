package db

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

func GetAllCandidates() ([]Candidate, error) {
	coll := DB.Database("public").Collection("candidates")
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return []Candidate{}, err
	}
	var candidates []Candidate
	err = cursor.All(context.TODO(), &candidates)
	if err != nil {
		return []Candidate{}, err
	}
	return candidates, nil
}

func (c *Candidate) Save() (*Candidate, error) {
	coll := DB.Database("public").Collection("candidates")
	_, err := coll.InsertOne(context.TODO(), c)
	if err != nil {
		return &Candidate{}, err
	}
	return c, nil
}

func ApproveCandidate(id string, approved bool) error {
	coll := DB.Database("public").Collection("candidates")
	filter := bson.D{{Key: "id", Value: id}}
	update := bson.D{{Key: "$set",
		Value: bson.D{{Key: "approved", Value: approved},
			{Key: "waitingForApprove", Value: false}},
	}}
	res, err := coll.UpdateOne(context.TODO(), filter, update)
	if res.MatchedCount == 0 {
		return errors.New("no candidate with id `" + id + "` found")
	}
	return err
}
