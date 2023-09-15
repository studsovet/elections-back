package db

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Key struct {
	Data   string `bson:"data" json:"data" bindings:"required"`
	Type   string `bson:"type" json:"type" bindings:"required"`
	PartID int    `bson:"partid" json:"partid" bindings:"required"`
}

func DropKeys() {
	DB.Database("protected").Collection("keys_data").Drop(context.TODO())
}

func (u *Key) SaveKey() (*Key, error) {
	var err error

	n, err := DB.Database("protected").Collection("keys_data").CountDocuments(context.TODO(), bson.D{{"type", u.Type}, {"partid", u.PartID}})

	if err != nil {
		return &Key{}, err
	}

	println(n)

	if n == 0 {
		_, err = DB.Database("protected").Collection("keys_data").InsertOne(context.TODO(), u)
		if err != nil {
			return &Key{}, err
		}
		return u, nil
	} else {
		_, err = DB.Database("protected").Collection("keys_data").UpdateOne(context.TODO(), bson.D{{"type", u.Type}, {"partid", u.PartID}}, bson.D{{"$set", bson.D{{"data", u.Data}}}})
		if err != nil {
			return &Key{}, err
		}
		return u, nil
	}
}

func GetPublicKey() (Key, error) {
	k := Key{}
	err := DB.Database("protected").Collection("keys_data").FindOne(context.TODO(), bson.D{{"type", "public"}}).Decode(&k)
	return k, err
}

func GetParsedPublicKey() (*rsa.PublicKey, error) {
	key, err := GetPublicKey()

	if err != nil {
		return &rsa.PublicKey{}, err
	}

	bKey, err := hex.DecodeString(key.Data)

	if err != nil {
		return &rsa.PublicKey{}, err
	}

	tkey, err := x509.ParsePKIXPublicKey(bKey)

	if err != nil {
		return &rsa.PublicKey{}, err
	}

	publicKey, ok := tkey.(*rsa.PublicKey)

	if ok {
		return publicKey, err
	} else {
		return &rsa.PublicKey{}, errors.New("can't parse public key")
	}
}

func GetPrivateKey() (Key, error) {
	k := Key{}
	err := DB.Database("protected").Collection("keys_data").FindOne(context.TODO(), bson.D{{"type", "private"}, {"partid", 0}}).Decode(&k)
	return k, err
}

func GetParsedPrivateKey() (*rsa.PrivateKey, error) {
	key, err := GetPrivateKey()

	if err != nil {
		return &rsa.PrivateKey{}, err
	}

	bKey, err := hex.DecodeString(key.Data)

	if err != nil {
		return &rsa.PrivateKey{}, err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(bKey)
	return privateKey, err
}

func PrivateKeyRecovery() error {
	coll := DB.Database("protected").Collection("keys_data")
	cursor, err := coll.Find(context.TODO(), bson.D{{"type", "private"}})
	if err != nil {
		return err
	}

	var votes []Key
	err = cursor.All(context.TODO(), &votes)
	if err != nil {
		return err
	}

	var privateKeyHex string = ""

	for _, v := range votes {
		privateKeyHex += v.Data
	}

	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return err
	}

	key, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		return err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return err
	}
	publicKeyHex := hex.EncodeToString(publicKeyBytes)

	publicKey, err := GetPublicKey()
	if err != nil {
		return err
	}

	if publicKey.Data != publicKeyHex {
		return errors.New("private key is wrong")
	}

	_, err = DB.Database("protected").Collection("keys_data").DeleteMany(context.TODO(), bson.D{{"type", "private"}})
	if err != nil {
		return err
	}

	(&Key{
		Data:   privateKeyHex,
		Type:   "private",
		PartID: 0,
	}).SaveKey()

	return nil
}

func (u *Key) IsTokenExist() (bool, error) {
	opts := options.Count().SetHint("_id_")
	count, err := DB.Database("protected").Collection("keys_data").CountDocuments(context.TODO(), bson.D{{"type", u.Type}, {"partid", u.PartID}}, opts)

	return int(count) > 0, err
}
