package servicios

import (
	"context"
	"log"

	modelo "github.com/LuisWaldman/fogon-servidor/modelo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IndiceServicio struct {
	db         *mongo.Client
	collection string
}

func NuevoIndiceServicio(db *mongo.Client) *IndiceServicio {
	return &IndiceServicio{
		db:         db,
		collection: "indiceCancion",
	}
}

func (s *IndiceServicio) CrearIndice(indice *modelo.ItemIndiceCancion) error {
	// Elimina el índice existente si existe
	s.BorrarPorNombreYOwner(indice.FileName, indice.Owner)

	col := s.db.Database(database).Collection(s.collection)
	inserta, err := col.InsertOne(context.TODO(), indice)
	if err != nil {
		log.Println("Error creando índice", "err", err)
		return err
	}
	log.Println("Índice creado", inserta)
	return nil
}

func (s *IndiceServicio) BuscarPorNombre(nombreArchivo string) (*modelo.ItemIndiceCancion, error) {
	col := s.db.Database(database).Collection(s.collection)
	var indice modelo.ItemIndiceCancion
	err := col.FindOne(context.TODO(), bson.M{"origen.fileName": nombreArchivo}).Decode(&indice)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &indice, nil
}

func (s *IndiceServicio) BuscarPorNombreYOwner(nombreArchivo string, owner string) (*modelo.ItemIndiceCancion, error) {
	col := s.db.Database(database).Collection(s.collection)
	var indice modelo.ItemIndiceCancion
	filtro := bson.M{
		"origen.fileName": nombreArchivo,
		"owner":           owner,
	}
	err := col.FindOne(context.TODO(), filtro).Decode(&indice)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &indice, nil
}

func (s *IndiceServicio) BuscarPorOwner(owner string) ([]modelo.ItemIndiceCancion, error) {
	col := s.db.Database(database).Collection(s.collection)
	cursor, err := col.Find(context.TODO(), bson.M{"owner": owner})
	if err != nil {
		log.Println("Error buscando índices por owner", "err", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var indices []modelo.ItemIndiceCancion
	for cursor.Next(context.TODO()) {
		var indice modelo.ItemIndiceCancion
		if err := cursor.Decode(&indice); err != nil {
			log.Println("Error decodificando índice", "err", err)
			continue
		}
		indices = append(indices, indice)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Error en cursor", "err", err)
		return nil, err
	}

	return indices, nil
}

func (s *IndiceServicio) BorrarPorNombre(nombreArchivo string) error {
	col := s.db.Database(database).Collection(s.collection)
	_, err := col.DeleteMany(context.TODO(), bson.M{"origen.fileName": nombreArchivo})
	if err != nil {
		log.Println("Error borrando índice por nombre", "err", err)
		return err
	}
	return nil
}

func (s *IndiceServicio) BorrarPorNombreYOwner(nombreArchivo string, owner string) error {
	col := s.db.Database(database).Collection(s.collection)
	filtro := bson.M{
		"origen.fileName": nombreArchivo,
		"owner":           owner,
	}
	_, err := col.DeleteOne(context.TODO(), filtro)
	if err != nil {
		log.Println("Error borrando índice", "err", err)
		return err
	}
	return nil
}

func (s *IndiceServicio) BorrarPorOwner(owner string) error {
	col := s.db.Database(database).Collection(s.collection)
	_, err := col.DeleteMany(context.TODO(), bson.M{"owner": owner})
	if err != nil {
		log.Println("Error borrando índices por owner", "err", err)
		return err
	}
	return nil
}

func (s *IndiceServicio) ListarTodos() ([]modelo.ItemIndiceCancion, error) {
	col := s.db.Database(database).Collection(s.collection)
	cursor, err := col.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println("Error listando todos los índices", "err", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var indices []modelo.ItemIndiceCancion
	for cursor.Next(context.TODO()) {
		var indice modelo.ItemIndiceCancion
		if err := cursor.Decode(&indice); err != nil {
			log.Println("Error decodificando índice", "err", err)
			continue
		}
		indices = append(indices, indice)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Error en cursor", "err", err)
		return nil, err
	}

	return indices, nil
}
