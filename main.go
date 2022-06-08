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
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func insertMongo(colHandler *mongo.Collection) {
	// sample for insertOne to MongoDB
	type address struct {
		ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		Street   string             `json:"street" bson:"street"`
		City     string             `json:"city" bson:"city"`
		State    string             `json:"state" bson:"state"`
		Zip_code string             `json:"zip_code" bson:"zip_code,omitempty"`
	}
	ctx := context.TODO()

	doc := address{
		Street:   "305 Knollcrest ct",
		City:     "Johns Creek",
		Zip_code: "30022",
		State:    "GA",
	}

	result, err := colHandler.InsertOne(ctx, doc)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.InsertedID)
}

func initMongoDB() *mongo.Collection {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	opts := options.Client()
	opts.SetConnectTimeout(2000 * time.Millisecond)
	opts.ApplyURI(uri)
	ctx := context.TODO()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	// defer func() {
	// 	if err := client.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	return client.Database("app").Collection("addresses")
}

func initRedis() {
	// redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping(ctx)
	log.Printf("redisClient status %v\n", status)

	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
}

func main() {
	// environment variables load
	if err := godotenv.Load(".devcontainer/.env"); err != nil {
		log.Println("No .env file found")
	}

	// mongodb
	colHandler := initMongoDB()
	insertMongo(colHandler)

}
