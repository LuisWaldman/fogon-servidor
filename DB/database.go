package db

import (
	"context"
	"log"

	config "github.com/LuisWaldman/fogon-servidor/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDB() (*mongo.Client, error) {
	var err error
	AppConfig := config.LoadConfiguration("./../config.json")
	var client *mongo.Client
	client, err = mongo.Connect(options.Client().ApplyURI(AppConfig.MONGODB_URI))
	if err != nil {
		log.Println("Error conectando con MongoDB", "err", err)
		return nil, err
	}
	// Ping the primary
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Println("Error haciendo ping a MongoDB", "err", err)
		return nil, err
	}
	log.Println("Conectado a MongoDB!", AppConfig.MONGODB_URI)
	return client, err
}
