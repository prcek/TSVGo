package main

import (
	//"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"

	// "regexp"
	// "strconv"
	// "strings"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	r, _, err := net.LookupSRV("mongodb", "tcp", "ts-cluster-ci22a.mongodb.net")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("bye", r)
}
