package servicios

import (
	"context"
	"log"

	modelo "github.com/LuisWaldman/fogon-servidor/modelo"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

	// Filtro para buscar documentos con el mismo fileName y origenUrl
	filter := bson.M{
		"fileName":  item.FileName,
		"origenUrl": item.OrigenUrl,
	}

	// Opciones para hacer upsert (update si existe, insert si no existe)
	opts := options.Replace().SetUpsert(true)

	resultado, err := col.ReplaceOne(context.TODO(), filter, item, opts)
	if err != nil {
		log.Println("Error agregando/actualizando canción en lista", "err", err)
		return err
	}

	if resultado.UpsertedID != nil {
		log.Println("Nueva canción insertada en lista", "id", resultado.UpsertedID)
	} else {
		log.Println("Canción actualizada en lista", "modificados", resultado.ModifiedCount)
	}

	return nil
}

func (s *ItemIndiceCancionServicio) GetCancionesPorListaID(listaID bson.ObjectID) []*modelo.ItemIndiceCancion {
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
	objID, err := bson.ObjectIDFromHex(id)
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
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = col.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		return err
	}
	return nil
}
