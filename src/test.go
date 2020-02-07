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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func testEv(cli *mongo.Client) {
	store.InitCollections(cli.Database("tsv"))
	ID, err := primitive.ObjectIDFromHex("5d964fea2bfbc5000ff2a19a")
	if err != nil {
		log.Fatal(err)
	}
	ev := store.GetEvent(ID)
	fmt.Println("EV:", ev)
	order := store.GetOrderForEvent(ID)
	fmt.Println("OR:", order)
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
	var result []string

	result, err = client.Database("testdb").ListCollectionNames(ctx, bson.D{})
	fmt.Println(result)

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
	testEv(client)

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

	fmt.Println("bye")
}
