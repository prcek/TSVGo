package main

import (
	//"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"time"

	// "regexp"
	// "strconv"
	// "strings"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Println("hello")
	clientOptions := options.Client().ApplyURI("mongodb+srv://xxxt:xxx@cluster0-tlc14.mongodb.net/test")

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
