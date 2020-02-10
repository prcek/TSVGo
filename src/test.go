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
	reportOrdersNo := 0
	reportPays := 0
	reportCash := 0
	reportPays2 := 0
	reportPays2In := 0
	reportPays2Out := 0
	var cumReport store.PaymentsReport
	for _, order := range orders {
		pays, paysAmount, paysAmountIn, paysAmountOut := store.GetPays(order.PublicKey)
		reportPays2 += paysAmount
		reportPays2In += paysAmountIn
		reportPays2Out += paysAmountOut
		paymentsReport := store.GetPaymentsForOrder(order.ID)
		if paysAmount != int(paymentsReport.CardAm) {
			fmt.Println("order ", order, " has pays!= paymentsCard", pays, paymentsReport.Payments)
		}

		anyPayment := (paysAmount != 0) || (paymentsReport.CashAm != 0) || (paymentsReport.CardAm != 0)
		if order.Status == "paid" {
			reportOrdersNo++
			cumReport.Add(&paymentsReport)
			reportPays += int(paymentsReport.CardAm)
			reportCash += int(paymentsReport.CashAm)
			if paysAmount+int(paymentsReport.CashAm) != int(order.Price) {
				fmt.Println("has diff pays amount!", order, pays)
				fmt.Println("...", paymentsReport.Payments, ".", paymentsReport.CashAm, ".", paymentsReport.CardAm)
			}
		} else if (order.Status == "cancelled") || (order.Status == "declined") {
			if anyPayment {
				fmt.Println("has pays amount!", order, pays)
			}
			reportPays += int(paymentsReport.CardAm)
			reportCash += int(paymentsReport.CashAm)
			cumReport.Add(&paymentsReport)
			//fmt.Println("cancelled order", paymentsReport)
		} else if (order.Status == "pending") || (order.Status == "reservation") {
			if anyPayment {
				fmt.Println("has paid amount!", order, pays, paymentsReport.Payments)
			}
		} else if order.Status == "refunded" {
			reportPays += int(paymentsReport.CardAm)
			reportCash += int(paymentsReport.CashAm)
			if (paymentsReport.CashAm + paymentsReport.CardAm) != 0 {
				fmt.Println("has paid amount!", order, pays, paymentsReport.Payments)
			}
			// fmt.Println("refunded order", paymentsReport)
			cumReport.Add(&paymentsReport)
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
	fmt.Printf("Orders %d, cash %d, card %d pays %d (in:%d out:%d)\n", reportOrdersNo, reportCash, reportPays, reportPays2, reportPays2In, reportPays2Out)
	fmt.Printf("Vybrano hotovost %d vratky %d, celkem %d\n", cumReport.CashAmIn, cumReport.CashAmOut, cumReport.CashAm)
	fmt.Printf("Vybrano karta %d vratky %d, celkem %d\n", cumReport.CardAmIn, cumReport.CardAmOut, cumReport.CardAm)
	fmt.Printf("Sumar: %d\n", cumReport.CardAm+cumReport.CashAm)
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
	fmt.Println("ples 1")
	processEventReport("5d964fea2bfbc5000ff2a19a")
	fmt.Println("ples 2")
	processEventReport("5de80f84521551000f25aa85")
	fmt.Println("kazma")
	processEventReport("5de8152f521551000f25ab42") //kazma
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

	fmt.Println("bye")
}
