package store

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var orderColl *mongo.Collection

// Order - mongo doc
type Order struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Status    string             `json:"status,omitempty" bson:"status,omitempty"`
	Type      string             `json:"type,omitempty" bson:"type,omitempty"`
	Price     int64              `json:"price,omitempty" bson:"price,omitempty"`
	PublicKey string             `json:"public_key,omitempty" bson:"public_key,omitempty"`
}

// GetOrderForEvent reads one Event
func GetOrderForEvent(ID primitive.ObjectID) Order {
	var doc Order
	err := orderColl.FindOne(context.Background(), bson.M{"event_id": ID, "type": "t"}).Decode(&doc)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}
