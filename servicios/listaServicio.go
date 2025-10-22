package servicios

import (
	"context"
	"log"

	modelo "github.com/LuisWaldman/fogon-servidor/modelo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ListaServicio struct {
	db         *mongo.Client
	collection string
}

func NuevoListaServicio(db *mongo.Client) *ListaServicio {
	return &ListaServicio{
		db:         db,
		collection: "listas",
	}
}

func (s *ListaServicio) CrearLista(lista *modelo.Lista) error {
	col := s.db.Database(database).Collection(s.collection)
	inserta, err := col.InsertOne(context.TODO(), lista)
	if err != nil {
		log.Println("Error creando lista", "err", err)
		return err
	}
	log.Println("Lista creada", inserta)
	return nil
}

func (s *ListaServicio) BuscarPorID(id string) (*modelo.Lista, error) {
	col := s.db.Database(database).Collection(s.collection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var lista modelo.Lista
	err = col.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&lista)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &lista, nil
}

func (s *ListaServicio) BuscarPorNombreYOwner(nombre string, owner string) (*modelo.Lista, error) {
	col := s.db.Database(database).Collection(s.collection)
	var lista modelo.Lista
	filtro := bson.M{
		"nombre": nombre,
		"owner":  owner,
	}
	err := col.FindOne(context.TODO(), filtro).Decode(&lista)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &lista, nil
}

func (s *ListaServicio) BuscarPorOwner(owner string) ([]modelo.Lista, error) {
	col := s.db.Database(database).Collection(s.collection)
	cursor, err := col.Find(context.TODO(), bson.M{"owner": owner})
	if err != nil {
		log.Println("Error buscando listas por owner", "err", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var listas []modelo.Lista
	for cursor.Next(context.TODO()) {
		var lista modelo.Lista
		if err := cursor.Decode(&lista); err != nil {
			log.Println("Error decodificando lista", "err", err)
			continue
		}
		listas = append(listas, lista)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Error en cursor", "err", err)
		return nil, err
	}

	return listas, nil
}

func (s *ListaServicio) ActualizarLista(id string, actualizacion bson.M) error {
	col := s.db.Database(database).Collection(s.collection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filtro := bson.M{"_id": objID}
	_, err = col.UpdateOne(context.TODO(), filtro, bson.M{"$set": actualizacion})
	if err != nil {
		log.Println("Error actualizando lista", "err", err)
		return err
	}
	return nil
}

func (s *ListaServicio) RenombrarLista(id string, nuevoNombre string) error {
	actualizacion := bson.M{"nombre": nuevoNombre}
	return s.ActualizarLista(id, actualizacion)
}

func (s *ListaServicio) BorrarPorID(id string) error {
	col := s.db.Database(database).Collection(s.collection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := col.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		log.Println("Error borrando lista", "err", err)
		return err
	}

	if result.DeletedCount == 0 {
		log.Printf("No se borró ninguna lista con ID: %s", id)
	} else {
		log.Printf("Lista borrada exitosamente. ID: %s, DeletedCount: %d", id, result.DeletedCount)
	}

	return nil
}

func (s *ListaServicio) BorrarPorOwner(owner string) error {
	col := s.db.Database(database).Collection(s.collection)
	result, err := col.DeleteMany(context.TODO(), bson.M{"owner": owner})
	if err != nil {
		log.Println("Error borrando listas por owner", "err", err)
		return err
	}
	log.Printf("Borradas %d listas del owner: %s", result.DeletedCount, owner)
	return nil
}

func (s *ListaServicio) BorrarPorNombreYOwner(nombre string, owner string) error {
	col := s.db.Database(database).Collection(s.collection)
	filtro := bson.M{
		"nombre": nombre,
		"owner":  owner,
	}
	result, err := col.DeleteMany(context.TODO(), filtro)
	if err != nil {
		log.Println("Error borrando listas por nombre y owner", "err", err)
		return err
	}
	log.Printf("Borradas %d listas con nombre: %s y owner: %s", result.DeletedCount, nombre, owner)
	return nil
}
