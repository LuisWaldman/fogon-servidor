package servicios

import (
	"context"
	"log"

	modelo "github.com/LuisWaldman/fogon-servidor/modelo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ItemIndiceCancionServicio struct {
	db         *mongo.Client
	collection string
}

func NuevoItemIndiceCancionServicio(db *mongo.Client) *ItemIndiceCancionServicio {
	return &ItemIndiceCancionServicio{
		db:         db,
		collection: "listaCanciones",
	}
}

func (s *ItemIndiceCancionServicio) AgregarCancion(item *modelo.ItemIndiceCancion) error {
	col := s.db.Database(database).Collection(s.collection)
	inserta, err := col.InsertOne(context.TODO(), item)
	if err != nil {
		log.Println("Error agregando canción a lista", "err", err)
		return err
	}
	log.Println("Canción agregada a lista", inserta)
	return nil
}

func (s *ItemIndiceCancionServicio) GetCancionesPorListaID(listaID primitive.ObjectID) []*modelo.ItemIndiceCancion {
	col := s.db.Database(database).Collection(s.collection)
	filter := bson.M{"listaId": listaID}
	cursor, err := col.Find(context.TODO(), filter)
	if err != nil {
		log.Println("Error obteniendo canciones por lista ID", "err", err)
		return nil
	}
	defer cursor.Close(context.TODO())

	var canciones []*modelo.ItemIndiceCancion
	for cursor.Next(context.TODO()) {
		var cancion modelo.ItemIndiceCancion
		if err := cursor.Decode(&cancion); err != nil {
			log.Println("Error decodificando canción", "err", err)
			continue
		}
		canciones = append(canciones, &cancion)
	}
	return canciones
}

func (s *ItemIndiceCancionServicio) BorrarPorListaID(id string) error {
	col := s.db.Database(database).Collection(s.collection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = col.DeleteMany(context.TODO(), bson.M{"listaId": objID})
	if err != nil {
		return err
	}
	return nil
}

func (s *ItemIndiceCancionServicio) BorrarPorID(id string) error {
	col := s.db.Database(database).Collection(s.collection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = col.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		return err
	}
	return nil
}
