package service

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Org struct {
	Domain      string    `bson:"domain"`
	OrgName     string    `bson:"name"`
	OrgType     string    `bson:"type"`
	OrgSlug     string    `bson:"org_slug"`
	CreatedAt   time.Time `bson:"createdAt"`
	LogoURL     string    `bson:"logoUrl,omitempty"`
	Description string    `bson:"description,omitempty"`
}

func getOrgCollection() *mongo.Collection {
	return MongoDatabase.Collection("organizations")
}

func LookupOrg(domain string) (*Org, error) {
	filter := bson.M{"domain": domain}
	log.Printf("filter %v", filter)
	var org Org
	log.Printf("Org: %v", org)
	err := getOrgCollection().FindOne(context.TODO(), filter).Decode(&org)
	if err != nil {
		return nil, err
	}

	return &org, nil
}
