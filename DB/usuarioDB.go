package DB

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

const database = "fogon"      // Cambia esto al nombre de tu base de datos
const collection = "usuarios" // Cambia esto al nombre de tu colección

type usuarioDB struct {
	Encontrado bool   `bson:-`
	Modologin  string `bson:"modologin"`
	Usuario    string `bson:"usuario"`
	Clave      string `bson:"clave"`
}

func crear_usuario(user usuarioDB) error {
	client, err := ConnectDB()
	if err != nil {
		return err
	}
	collection := client.Database(database).Collection(collection)
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Println("Error creando usuario", "err", err)
		return err
	}
	return nil
}

func buscarxusuario(usuario string) (*usuarioDB, error) {
	client, err := ConnectDB()
	if err != nil {
		return nil, err
	}
	collection := client.Database(database).Collection(collection)
	var user usuarioDB
	err = collection.FindOne(context.TODO(), bson.M{"usuario": usuario}).Decode(&user)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return &usuarioDB{Encontrado: false}, nil
		}
		return nil, err
	}
	user.Encontrado = true
	return &user, nil
}

func borrarxusuario(usuario string) error {
	client, err := ConnectDB()
	if err != nil {
		return err
	}
	collection := client.Database(database).Collection(collection)
	_, err = collection.DeleteOne(context.TODO(), bson.M{"usuario": usuario})
	if err != nil {
		log.Println("Error borrando usuario", "err", err)
		return err
	}
	return nil
}
