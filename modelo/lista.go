package modelo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Lista struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nombre         string             `bson:"nombre" json:"nombre"`
	Owner          string             `bson:"owner" json:"owner"`
	TotalCanciones int                `bson:"total_canciones" json:"total_canciones"`
}

func NuevaLista(nombre string, owner string) *Lista {
	return &Lista{
		ID:             primitive.NewObjectID(),
		Nombre:         nombre,
		Owner:          owner,
		TotalCanciones: 0,
	}
}
