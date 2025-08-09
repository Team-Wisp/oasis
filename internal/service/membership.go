package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Membership struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UserID        primitive.ObjectID `bson:"userId"`
	OrgID         primitive.ObjectID `bson:"orgId"`
	DisplayHandle string             `bson:"displayHandle"`
	Role          string             `bson:"role"` // "member"
	CreatedAt     time.Time          `bson:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt"`
}

func getMembershipCollection() *mongo.Collection {
	return MongoDatabase.Collection("memberships")
}

// UpsertMembership returns an existing or newly-created membership ID.
// Idempotent on (userId, orgId). Generates a handle once.
func UpsertMembership(ctx context.Context, userId, orgId primitive.ObjectID) (primitive.ObjectID, string, error) {
	coll := getMembershipCollection()

	// try find existing
	var m Membership
	err := coll.FindOne(ctx, bson.M{"userId": userId, "orgId": orgId}).Decode(&m)
	if err == nil {
		return m.ID, m.DisplayHandle, nil
	}
	if err != mongo.ErrNoDocuments {
		return primitive.NilObjectID, "", err
	}

	// create new
	handle := generateHandle()
	now := time.Now()
	doc := Membership{
		UserID:        userId,
		OrgID:         orgId,
		DisplayHandle: handle,
		Role:          "member",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	res, err := coll.InsertOne(ctx, doc)
	if err != nil {
		// race: unique index may reject, try read again
		if writeErr, ok := err.(mongo.WriteException); ok && len(writeErr.WriteErrors) > 0 && writeErr.WriteErrors[0].Code == 11000 {
			if err := coll.FindOne(ctx, bson.M{"userId": userId, "orgId": orgId}).Decode(&m); err == nil {
				return m.ID, m.DisplayHandle, nil
			}
		}
		return primitive.NilObjectID, "", err
	}
	id, _ := res.InsertedID.(primitive.ObjectID)
	return id, handle, nil
}

// EnsureMembershipIndexes can be called once at boot.
func EnsureMembershipIndexes(ctx context.Context) error {
	c := getMembershipCollection()
	// unique compound (userId, orgId)
	models := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "orgId", Value: 1}},
			Options: indexOpts(true),
		},
		{Keys: bson.D{{Key: "orgId", Value: 1}}},
		{Keys: bson.D{{Key: "userId", Value: 1}}},
	}
	_, err := c.Indexes().CreateMany(ctx, models)
	return err
}

// ---- tiny handle generator ----
var adjectives = []string{"calm", "brisk", "curious", "steady", "bright", "quiet", "bold", "wry", "kind", "spry", "eager", "mellow"}
var animals = []string{"otter", "sparrow", "lynx", "walrus", "fox", "owl", "badger", "heron", "orca", "koala", "yak", "hare"}

func generateHandle() string {
	rand.Seed(time.Now().UnixNano())
	a := adjectives[rand.Intn(len(adjectives))]
	b := animals[rand.Intn(len(animals))]
	n := rand.Intn(90) + 10 // 10..99
	return fmt.Sprintf("%s-%s-%d", a, b, n)
}

// helper to set unique options
func indexOpts(unique bool) *options.IndexOptions {
	o := options.Index()
	o.SetUnique(unique)
	o.SetBackground(true)
	return o
}
