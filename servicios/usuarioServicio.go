package servicios

import (
	"context"
	"log"

	modelo "github.com/LuisWaldman/fogon-servidor/modelo"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const database = "fogon" // Cambia esto al nombre de tu base de datos

type UsuarioServicio struct {
	db         *mongo.Client
	collection string
}

func NuevoUsuarioServicio(db *mongo.Client) *UsuarioServicio {
	return &UsuarioServicio{
		db:         db,
		collection: "usuarios", // Cambia esto al nombre de tu colección
	}
}

func (s *UsuarioServicio) CrearUsuario(user modelo.Usuario) error {
	col := s.db.Database(database).Collection(s.collection)
	inserta, err := col.InsertOne(context.TODO(), user)
	if err != nil {
		log.Println("Error creando usuario", "err", err)
		return err
	}
	log.Println("Usuario creado", inserta)
	return nil
}

func (s *UsuarioServicio) BuscarPorUsuario(usuario string) (*modelo.Usuario, error) {
	col := s.db.Database(database).Collection(s.collection)
	var user modelo.Usuario
	err := col.FindOne(context.TODO(), bson.M{"usuario": usuario}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &modelo.Usuario{Encontrado: false}, nil
		}
		return nil, err
	}
	user.Encontrado = true
	return &user, nil
}

func (s *UsuarioServicio) BorrarPorUsuario(usuario string) error {
	col := s.db.Database(database).Collection(s.collection)
	_, err := col.DeleteOne(context.TODO(), bson.M{"usuario": usuario})
	if err != nil {
		log.Println("Error borrando usuario", "err", err)
		return err
	}
	return nil
}

func (s *UsuarioServicio) ActualizarUsuario(user *modelo.Usuario) error {
	col := s.db.Database(database).Collection(s.collection)
	filter := bson.M{"usuario": user.Usuario}
	update := bson.M{"$set": user}
	_, err := col.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println("Error actualizando usuario", "err", err)
		return err
	}
	return nil
}
