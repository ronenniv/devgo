package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(".devcontainer/.env"); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	opts := options.Client()
	opts.SetConnectTimeout(1 * time.Second)
	opts.ApplyURI(uri)
	ctx := context.TODO()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("app").Collection("addresses")

	type address struct {
		ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		Street   string             `json:"street" bson:"street"`
		City     string             `json:"city" bson:"city"`
		State    string             `json:"state" bson:"state"`
		Zip_code string             `json:"zip_code" bson:"zip_code,omitempty"`
	}

	doc := address{
		Street:   "305 Knollcrest ct",
		City:     "Johns Creek",
		Zip_code: "30022",
		State:    "GA",
	}
	if err != nil {
		fmt.Println(err)
	}

	result, err := coll.InsertOne(ctx, doc)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.InsertedID)

}
