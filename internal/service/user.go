package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func getUserCollection() *mongo.Collection {
	return MongoDatabase.Collection("users")
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	EmailHash string             `bson:"emailHash"`
	Password  string             `bson:"password"` // bcrypt hash
	Slug      string             `bson:"org"`
	OrgType   string             `bson:"orgType"`
	CreatedAt time.Time          `bson:"createdAt"`
}

func HashPassword(preHashed string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(preHashed), bcrypt.DefaultCost)
	return string(hash), err
}

func SaveUser(user User) (primitive.ObjectID, error) {
	coll := getUserCollection()
	res, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return oid, nil
	}
	return primitive.NilObjectID, fmt.Errorf("inserted id not ObjectID")
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

func IsValidEmailHash(emailHash string) bool {
	// Example validation: Ensure the hash is a fixed-length alphanumeric string
	const hashLength = 64 // Adjust based on the actual hash length
	if len(emailHash) != hashLength {
		return false
	}
	for _, char := range emailHash {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return false
		}
	}
	return true
}
