package servicios

import (
	"context"
	"log"

	modelo "github.com/LuisWaldman/fogon-servidor/modelo"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PerfilServicio struct {
	db         *mongo.Client
	collection string
}

func NuevoPerfilServicio(db *mongo.Client) *PerfilServicio {
	return &PerfilServicio{
		db:         db,
		collection: "perfil", // Cambia esto al nombre de tu colección
	}
}

func (s *PerfilServicio) CrearPerfil(user modelo.Perfil) error {
	s.BorrarPorUsuario(user.Usuario) // Elimina el perfil existente antes de crear uno nuevo
	col := s.db.Database(database).Collection(s.collection)
	inserta, err := col.InsertOne(context.TODO(), user)
	if err != nil {
		log.Println("Error creando Perfil", "err", err)
		return err
	}
	log.Println("Perfil creado", inserta)
	return nil
}

func (s *PerfilServicio) BuscarPorUsuario(usuario string) (*modelo.Perfil, error) {
	col := s.db.Database(database).Collection(s.collection)
	var user modelo.Perfil
	err := col.FindOne(context.TODO(), bson.M{"usuario": usuario}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *PerfilServicio) BorrarPorUsuario(usuario string) error {
	col := s.db.Database(database).Collection(s.collection)
	_, err := col.DeleteOne(context.TODO(), bson.M{"usuario": usuario})
	if err != nil {
		log.Println("Error borrando usuario", "Perfil", err)
		return err
	}
	return nil
}
