package servicios

import (
	"context"
	"log"

	modelo "github.com/LuisWaldman/fogon-servidor/modelo"

	"go.mongodb.org/mongo-driver/bson"
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

func (s *ListaServicio) CrearLista(Nombre string, Owner string) error {
	col := s.db.Database(database).Collection(s.collection)
	lista := modelo.NuevaLista(Nombre, Owner)
	inserta, err := col.InsertOne(context.TODO(), lista)
	if err != nil {
		log.Println("Error creando lista", "err", err)
		return err
	}
	log.Println("Lista creada", inserta)
	return nil
}

func (s *ListaServicio) BuscarPorNombreYOwner(nombre string, owner string) (*modelo.Lista, error) {
	col := s.db.Database(database).Collection(s.collection)
	filter := bson.M{"owner": owner, "nombre": nombre}
	cursor, err := col.Find(context.TODO(), filter)
	if err != nil {
		log.Println("Error obteniendo listas", "err", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var lista *modelo.Lista
	for cursor.Next(context.TODO()) {
		var l modelo.Lista
		if err := cursor.Decode(&l); err != nil {
			log.Println("Error decodificando lista", "err", err)
			continue
		}
		lista = &l
	}
	if err := cursor.Err(); err != nil {
		log.Println("Error iterando cursor", "err", err)
		return nil, err
	}
	return lista, nil
}
func (s *ListaServicio) ActualizarLista(lista *modelo.Lista) error {
	col := s.db.Database(database).Collection(s.collection)
	filter := bson.M{"_id": lista.ID}
	update := bson.M{"$set": bson.M{
		"nombre":          lista.Nombre,
		"owner":           lista.Owner,
		"total_canciones": lista.TotalCanciones,
	}}
	_, err := col.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println("Error actualizando lista", "err", err)
		return err
	}
	log.Println("Lista actualizada", lista.ID)
	return nil
}
func (s *ListaServicio) BorrarPorID(id string) error {
	col := s.db.Database(database).Collection(s.collection)
	filter := bson.M{"_id": id}
	_, err := col.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Println("Error borrando lista", "err", err)
		return err
	}
	log.Println("Lista borrada", id)
	return nil
}
