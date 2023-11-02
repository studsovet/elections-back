package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func GetCandidates(userId string) ([]Candidate, error) {
	//TODO надо смотреть на userId и выдавать только нужных кандидатов
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

func (c *Candidate) SaveCandidate() (*Candidate, error) {
	coll := DB.Database("public").Collection("candidates")
	_, err := coll.InsertOne(context.TODO(), c)
	if err != nil {
		return &Candidate{}, err
	}
	return c, nil
}
