package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Org struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Domain      string             `bson:"domain"`
	OrgName     string             `bson:"name"`
	OrgType     string             `bson:"type"`
	OrgSlug     string             `bson:"org_slug"`
	CreatedAt   time.Time          `bson:"createdAt"`
	LogoURL     string             `bson:"logoUrl,omitempty"`
	Description string             `bson:"description,omitempty"`
}

func getOrgCollection() *mongo.Collection {
	return MongoDatabase.Collection("organizations")
}

func LookupOrg(domain string) (*Org, error) {
	filter := bson.M{"domain": domain}
	var org Org
	err := getOrgCollection().FindOne(context.TODO(), filter).Decode(&org)
	if err != nil {
		return nil, err
	}

	return &org, nil
}
func LookupOrgBySlug(slug string) (*Org, error) {
	filter := bson.M{"org_slug": slug}
	var org Org
	err := getOrgCollection().FindOne(context.TODO(), filter).Decode(&org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}
