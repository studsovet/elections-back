package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	ID         string `bson:"id" json:"id" bindings:"required"`
	IsObserver bool   `bson:"is_observer" json:"is_observer" bindings:"required"`
	IsAdmin    bool   `bson:"is_admin" json:"is_admin" bindings:"required"`
}

func (u *User) SaveUser() (*User, error) {
	var err error

	_, err = GetUserByID(u.ID)

	if err != nil {
		_, err = DB.Database("protected").Collection("auth_data").InsertOne(context.TODO(), u)
		if err != nil {
			return &User{}, err
		}
		return u, nil
	} else {
		DB.Database("protected").Collection("auth_data").FindOneAndUpdate(context.TODO(), bson.D{{"id", u.ID}}, bson.D{{"is_observer", u.IsObserver}, {"is_admin", u.IsAdmin}})
		return u, nil
	}
}

func GetUserByID(id string) (User, error) {
	u := User{}
	err := DB.Database("protected").Collection("auth_data").FindOne(context.TODO(), bson.D{{"id", id}}).Decode(&u)
	return u, err
}
