package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Test struct {
	Test string `json:"test" bindings:"required"`
}

func (t *Test) SaveTest() (*Test, error) {
	var err error
	coll := DB.Database("test").Collection("test")
	_, err = coll.InsertOne(context.TODO(), bson.D{
		{Key: "test", Value: t.Test},
	})
	if err != nil {
		return &Test{}, err
	}
	return t, nil
}

func GetTest() Test {
	coll := DB.Database("test").Collection("test")
	var t Test
	err := coll.FindOne(context.TODO(), bson.D{{}}).Decode(&t)
	if err == mongo.ErrNoDocuments {
		// return "Not found"
	}
	if err != nil {
		panic(err)
	}
	return t
}
