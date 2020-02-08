package main

import (
	"store"

	"go.mongodb.org/mongo-driver/bson/primitive"
	//"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	// "regexp"
	// "strconv"
	// "strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func processEventReport(eventID string) {
	ID, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		log.Fatal(err)
	}

	orders := store.GetAllOrdersForEvent(ID)
	for _, order := range orders {
		pays, paysAmount := store.GetPays(order.PublicKey)
		if order.Status == "paid" {
			if paysAmount != int(order.Price) {
				fmt.Println("has diff pays amount!", order, pays)
			}
		} else if order.Status == "cancelled" {
			if paysAmount != 0 {
				fmt.Println("has pays amount!", order, pays)
			}
		} else {
			fmt.Println("wrong status", order)
		}
		/*
			if len(pays) == 0 {

			} else {
				if amount != int(order.Price) {
					fmt.Println(order, amount, pays)
				}
			}
		*/
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	mongoURI, mongoURIp := os.LookupEnv("MONGODB_URI")
	if !mongoURIp {
		log.Fatal("env MONGODB_URI missing")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal("NewClient", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Ping", err)
	}

	// var sr bson.M
	/*
		var sr2 struct {
			ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
			Amount int64              `bson:"amounta,omitempty"`
		}

		sr2.Amount = 1

		rr, _ := client.Database("testdb").Collection("pbla").InsertOne(context.Background(), sr2)
		fmt.Println(rr)
	*/
	// testEv(client)
	store.InitCollections(client.Database("tsv"))
	store.ReadPaysFromCSV("pays_csv/all.csv")
	processEventReport("5d964fea2bfbc5000ff2a19a")

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

	fmt.Println("bye")
}
