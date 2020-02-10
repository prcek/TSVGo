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

// PaymentsReport ...
type PaymentsReport struct {
	Payments                                                 []Payment
	CashAm, CardAm, CashAmIn, CashAmOut, CardAmIn, CardAmOut int64
}

// Add - sum
func (r *PaymentsReport) Add(a *PaymentsReport) {
	r.CardAm += a.CardAm
	r.CardAmIn += a.CardAmIn
	r.CardAmOut += a.CardAmOut
	r.CashAm += a.CashAm
	r.CashAmIn += a.CashAmIn
	r.CashAmOut += a.CashAmOut
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
func GetPaymentsForOrder(orderID primitive.ObjectID) PaymentsReport {
	cursor, err := paymentsColl.Find(context.Background(), bson.M{"order_id": orderID})
	if err != nil {
		log.Fatal(err)
	}

	var r PaymentsReport
	if err = cursor.All(context.TODO(), &r.Payments); err != nil {
		log.Fatal(err)
	}

	for _, p := range r.Payments {
		if p.Type == "cash" {
			if p.Status == "paid" {
				r.CashAm += p.Amount
				r.CashAmIn += p.Amount
				if p.Amount < 0 {
					log.Fatal("negative payment", p)
				}
			} else if p.Status == "refunded" {
				r.CashAm -= p.Amount
				r.CashAmOut -= p.Amount
				if p.Amount < 0 {
					log.Fatal("negative refund payment", p)
				}
			} else {
				log.Fatal("unexpected payment status - ", p.Status)
			}
		} else if p.Type == "card" {
			if p.Status == "paid" {
				r.CardAm += p.Amount
				r.CardAmIn += p.Amount
				if p.Amount < 0 {
					log.Fatal("negative payment", p)
				}
			} else if p.Status == "refunded" {
				r.CardAm -= p.Amount
				r.CardAmOut -= p.Amount
				if p.Amount < 0 {
					log.Fatal("negative refund payment", p)
				}
			} else if p.Status == "declined" {

			} else {
				log.Fatal("unexpected payment status - ", p.Status)
			}
		} else {
			log.Fatal("unexpected payment type - ", p.Type)
		}
	}

	return r
}
