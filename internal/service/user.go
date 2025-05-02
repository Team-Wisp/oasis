package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func getUserCollection() *mongo.Collection {
	return MongoDatabase.Collection("users")
}

type User struct {
	EmailHash string    `bson:"emailHash"`
	Password  string    `bson:"password"` // bcrypt hash
	OrgSlug   string    `bson:"org"`
	OrgType   string    `bson:"orgType"`
	CreatedAt time.Time `bson:"createdAt"`
}

func HashPassword(preHashed string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(preHashed), bcrypt.DefaultCost)
	return string(hash), err
}

func SaveUser(user User) error {
	coll := getUserCollection()
	_, err := coll.InsertOne(context.TODO(), user)
	return err
}

func DoesUserExist(emailHash string) (bool, error) {
	filter := bson.M{"emailHash": emailHash}
	count, err := getUserCollection().CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func GetUserByEmailHash(emailHash string) (*User, error) {
	filter := bson.M{"emailHash": emailHash}
	var user User
	err := getUserCollection().FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CheckPassword(hashed string, inputHash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(inputHash)) == nil
}
