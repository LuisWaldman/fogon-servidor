package servicios

import (
	"context"
	"log"

	modelo "github.com/LuisWaldman/fogon-servidor/modelo"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ListaCancionServicio struct {
	db         *mongo.Client
	collection string
}

func NuevoListaCancionServicio(db *mongo.Client) *ListaCancionServicio {
	return &ListaCancionServicio{
		db:         db,
		collection: "listaCanciones",
	}
}

func (s *ListaCancionServicio) AgregarCancion(listaCancion *modelo.ListaCancion) error {
	col := s.db.Database(database).Collection(s.collection)
	inserta, err := col.InsertOne(context.TODO(), listaCancion)
	if err != nil {
		log.Println("Error agregando canción a lista", "err", err)
		return err
	}
	log.Println("Canción agregada a lista", inserta)
	return nil
}
