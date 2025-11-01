package servicios

import (
	"context"
	"log"

	modelo "github.com/LuisWaldman/fogon-servidor/modelo"

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
