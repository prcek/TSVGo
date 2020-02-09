package store

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Event - mongo doc
type Event struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
}

var eventColl *mongo.Collection

// InitCollections setups Collections
func InitCollections(db *mongo.Database) {
	eventColl = db.Collection("events")
	orderColl = db.Collection("orders")
	paymentsColl = db.Collection("payments")
}

// GetEvent reads one Event
func GetEvent(ID primitive.ObjectID) Event {
	var ev Event
	err := eventColl.FindOne(context.Background(), bson.M{"_id": ID}).Decode(&ev)
	if err != nil {
		log.Fatal(err)
	}
	return ev
}
