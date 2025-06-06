package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client

func ConnectDB() {
	var err error
	client, err = mongo.Connect(options.Client().ApplyURI(MONGODB_URI))
	if err != nil {
		log.Fatalln("Error conectando con MongoDB", "err", err)
	}
	// Ping the primary
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatalln("Error haciendo ping a MongoDB", "err", err)
	}
	log.Println("Conectado a MongoDB!")
}

func FindUserByName(nombre string) (bson.M, error) {
	coll := client.Database("fogon").Collection("usuarios")
	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{"nombre", nombre}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("no se encontro el usuario %s: %w", nombre, err)
	}
	if err != nil {
		return nil, fmt.Errorf("error buscando usuario %s: %w", nombre, err)
	}
	return result, nil
}
