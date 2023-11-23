package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func GetPublicKey(election_id string) (PublicKey, error) {
	key := PublicKey{}
	filter := bson.D{{Key: "id", Value: election_id}}
	err := DB.Database("public").Collection("public_keys").FindOne(context.TODO(), filter).Decode(&key)
	return key, err
}

func (k *PublicKey) Save() (*PublicKey, error) {
	_, err := DB.Database("public").Collection("public_keys").DeleteMany(context.TODO(), bson.D{{Key: "id", Value: k.ID}})

	if err != nil {
		return nil, err
	}

	_, err = DB.Database("public").Collection("public_keys").InsertOne(context.TODO(), k)
	return k, err
}

func GetPrivateKey(election_id string) (PrivateKey, error) {
	key := PrivateKey{}
	filter := bson.D{{Key: "id", Value: election_id}}
	err := DB.Database("public").Collection("private_keys").FindOne(context.TODO(), filter).Decode(&key)
	return key, err
}

func (k *PrivateKey) Save() (*PrivateKey, error) {
	_, err := DB.Database("public").Collection("private_keys").DeleteMany(context.TODO(), bson.D{{Key: "id", Value: k.ID}})

	if err != nil {
		return nil, err
	}

	_, err = DB.Database("public").Collection("private_keys").InsertOne(context.TODO(), k)
	return k, err
}
