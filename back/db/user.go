package db

import (
	"context"
	"html"
	"strings"

	token "elections-back/utils"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         string `bson:"id" json:"id" bindings:"required"`
	Username   string `bson:"username" json:"username" bindings:"required"`
	Password   string `bson:"password" json:"password" bindings:"required"`
	IsObserver bool   `bson:"is_observer" json:"is_observer" bindings:"required"`
	IsAdmin    bool   `bson:"is_admin" json:"is_admin" bindings:"required"`
}

func (u *User) SaveUser() (*User, error) {
	var err error
	err = u.BeforeSave()
	if err != nil {
		return &User{}, err
	}
	_, err = DB.Database("protected").Collection("auth_data").InsertOne(context.TODO(), u)
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave() error {
	u.ID = uuid.NewString()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	return nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string) (string, error) {
	var err error

	u := User{}
	err = DB.Database("protected").Collection("auth_data").FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&u)
	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(u.ID)

	if err != nil {
		return "", err
	}

	return token, nil
}

func CountObservers() (int, error) {
	opts := options.Count().SetHint("_id_")
	count, err := DB.Database("protected").Collection("auth_data").CountDocuments(context.TODO(), bson.D{{"is_observer", true}}, opts)

	return int(count), err
}

func GetUserByID(id string) (User, error) {
	u := User{}
	err := DB.Database("protected").Collection("auth_data").FindOne(context.TODO(), bson.D{{"id", id}}).Decode(&u)
	return u, err
}
