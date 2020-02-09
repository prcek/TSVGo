package store

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var paymentsColl *mongo.Collection

// Payment - mongo doc
type Payment struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Status         string             `json:"status,omitempty" bson:"status,omitempty"`
	Type           string             `json:"type,omitempty" bson:"type,omitempty"`
	Amount         int64              `json:"amount,omitempty" bson:"amount,omitempty"`
	OrderPublicKey string             `json:"order_public_key,omitempty" bson:"order_public_key,omitempty"`
}

// GetOrderForEvent reads one Event
/*
func GetOrderForEvent(eventID primitive.ObjectID) Order {
	var doc Order
	err := orderColl.FindOne(context.Background(), bson.M{"event_id": eventID, "type": "t"}).Decode(&doc)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}
*/

// GetPaymentsForOrder returns all order for event
func GetPaymentsForOrder(orderID primitive.ObjectID) ([]Payment, int64, int64) {
	cursor, err := paymentsColl.Find(context.Background(), bson.M{"order_id": orderID})
	if err != nil {
		log.Fatal(err)
	}
	var results []Payment
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	var cardAm int64 = 0
	var cashAm int64 = 0
	for _, p := range results {
		if p.Type == "cash" {
			if p.Status == "paid" {
				cashAm += p.Amount
			} else if p.Status == "refunded" {
				cashAm -= p.Amount
			} else {
				log.Fatal("unexpected payment status - ", p.Status)
			}
		} else if p.Type == "card" {
			if p.Status == "paid" {
				cardAm += p.Amount
			} else if p.Status == "refunded" {
				cardAm -= p.Amount
			} else if p.Status == "declined" {

			} else {
				log.Fatal("unexpected payment status - ", p.Status)
			}
		} else {
			log.Fatal("unexpected payment type - ", p.Type)
		}
	}
	return results, cashAm, cardAm
}
