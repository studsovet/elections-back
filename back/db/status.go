package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Status struct {
	Code        int       `bson:"id" json:"id" bindings:"required"`
	Description string    `bson:"description" json:"description" bindings:"required"`
	Time        time.Time `bson:"time" json:"time" bindings:"required"`
}

func SetStatus(code int, desc string) error {
	s := Status{
		Code:        code,
		Description: desc,
		Time:        time.Now(),
	}

	_, err := DB.Database("public").Collection("status").InsertOne(context.TODO(), s)

	return err
}

func GetLastStatus() Status {
	var s Status
	sort := bson.D{{"time", -1}}
	opts := options.FindOne().SetSort(sort)

	err := DB.Database("public").Collection("status").FindOne(context.TODO(), bson.D{}, opts).Decode(&s)

	if err != nil {
		return Status{
			Code: -1,
		}
	}

	return s
}
