package servicios

import (
	"context"
	"log"

	modelo "github.com/LuisWaldman/fogon-servidor/modelo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CancionServicio struct {
	db             *mongo.Client
	collection     string
	indiceServicio *IndiceServicio
}

func NuevoCancionServicio(db *mongo.Client) *CancionServicio {
	return &CancionServicio{
		db:             db,
		collection:     "cancion",
		indiceServicio: NuevoIndiceServicio(db),
	}
}

func (s *CancionServicio) CrearCancion(cancion modelo.Cancion) error {
	s.BorrarPorNombreYOwner(cancion.NombreArchivo, cancion.Owner) // Elimina la canción existente antes de crear una nueva
	col := s.db.Database(database).Collection(s.collection)
	inserta, err := col.InsertOne(context.TODO(), cancion)
	if err != nil {
		log.Println("Error creando Cancion", "err", err)
		return err
	}

	indice := modelo.BuildFromCancion(&cancion)
	err = s.indiceServicio.CrearIndice(indice)
	if err != nil {
		log.Println("Error creando índice para canción", "err", err)
		// No retornamos el error para no fallar la creación de la canción
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

func (s *CancionServicio) BuscarPorNombreYOwner(nombreArchivo string, owner string) (*modelo.Cancion, error) {
	col := s.db.Database(database).Collection(s.collection)
	var cancion modelo.Cancion
	filtro := bson.M{
		"nombreArchivo": nombreArchivo,
		"owner":         owner,
	}
	err := col.FindOne(context.TODO(), filtro).Decode(&cancion)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &cancion, nil
}

func (s *CancionServicio) BuscarPorOwner(owner string) ([]modelo.Cancion, error) {
	col := s.db.Database(database).Collection(s.collection)
	cursor, err := col.Find(context.TODO(), bson.M{"owner": owner})
	if err != nil {
		log.Println("Error buscando canciones por owner", "err", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var canciones []modelo.Cancion
	for cursor.Next(context.TODO()) {
		var cancion modelo.Cancion
		if err := cursor.Decode(&cancion); err != nil {
			log.Println("Error decodificando canción", "err", err)
			continue
		}
		canciones = append(canciones, cancion)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Error en cursor", "err", err)
		return nil, err
	}

	return canciones, nil
}

func (s *CancionServicio) BorrarPorNombre(nombreArchivo string) error {
	col := s.db.Database(database).Collection(s.collection)
	_, err := col.DeleteMany(context.TODO(), bson.M{"nombreArchivo": nombreArchivo})
	if err != nil {
		log.Println("Error borrando cancion", "err", err)
		return err
	}

	// También borrar del índice
	s.indiceServicio.BorrarPorNombre(nombreArchivo)

	return nil
}

func (s *CancionServicio) BorrarPorNombreYOwner(nombreArchivo string, owner string) error {
	col := s.db.Database(database).Collection(s.collection)
	filtro := bson.M{
		"nombreArchivo": nombreArchivo,
		"owner":         owner,
	}
	_, err := col.DeleteOne(context.TODO(), filtro)
	if err != nil {
		log.Println("Error borrando cancion por nombre y owner", "err", err)
		return err
	}

	// También borrar del índice
	s.indiceServicio.BorrarPorNombreYOwner(nombreArchivo, owner)

	return nil
}
