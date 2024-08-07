package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"

	"github.com/tantanok221/cloudflare-image-backend/internal/helper"
)

var (
	MongoString string = helper.GetEnv("MongoString")
)

func Init() (*mongo.Collection, func()) {
	//dsn := fmt.Sprintf("%v://%v:%v@%v:%v/postgres", DBNAME, USER, PASSWORD, HOST, PORT)
	//dsn := fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v", USER, PASSWORD, HOST, PORT, DBNAME)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(MongoString).SetServerAPIOptions(serverAPI)
	db, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}
	closeDB := func() {
		if err := db.Disconnect(context.TODO()); err != nil {
			log.Printf("Error disconnecting from mongodb: %v", err)
		}
	}

	// Send a ping to confirm a successful connection
	if err := db.Database("cloudflare-image").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	client := db.Database("cloudflare-image").Collection("douren-image")
	print("You had connected to mongodb")
	return client, closeDB

}
