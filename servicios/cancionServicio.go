package servicios

import (
	"context"
	"log"

	modelo "github.com/LuisWaldman/fogon-servidor/modelo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CancionServicio struct {
	db         *mongo.Client
	collection string
}

func NuevoCancionServicio(db *mongo.Client) *CancionServicio {
	return &CancionServicio{
		db:         db,
		collection: "cancion",
	}
}

func (s *CancionServicio) CrearCancion(cancion modelo.Cancion) error {
	s.BorrarPorNombre(cancion.NombreArchivo) // Elimina la canción existente antes de crear una nueva
	col := s.db.Database(database).Collection(s.collection)
	inserta, err := col.InsertOne(context.TODO(), cancion)
	if err != nil {
		log.Println("Error creando Cancion", "err", err)
		return err
	}
	log.Println("Cancion creada", inserta)
	return nil
}

func (s *CancionServicio) BuscarPorNombre(nombreArchivo string) (*modelo.Cancion, error) {
	col := s.db.Database(database).Collection(s.collection)
	var cancion modelo.Cancion
	err := col.FindOne(context.TODO(), bson.M{"nombreArchivo": nombreArchivo}).Decode(&cancion)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &cancion, nil
}

func (s *CancionServicio) BorrarPorNombre(nombreArchivo string) error {
	col := s.db.Database(database).Collection(s.collection)
	_, err := col.DeleteOne(context.TODO(), bson.M{"nombreArchivo": nombreArchivo})
	if err != nil {
		log.Println("Error borrando cancion", "err", err)
		return err
	}
	return nil
}
