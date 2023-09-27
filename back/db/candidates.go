package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type Candidate struct {
	Name     string `bson:"name" json:"name" bindings:"required"`
	Surname  string `bson:"surname" json:"surname" bindings:"required"`
	Program  string `bson:"program" json:"program" bindings:"required"`
	PhotoUrl string `bson:"photourl" json:"photourl" bindings:"required"`
}

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
