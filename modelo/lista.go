package modelo

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Lista struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Nombre         string        `bson:"nombre" json:"nombre"`
	Owner          string        `bson:"owner" json:"owner"`
	TotalCanciones int           `bson:"total_canciones" json:"total_canciones"`
}

func NuevaLista(nombre string, owner string) *Lista {
	return &Lista{
		ID:             bson.NewObjectID(),
		Nombre:         nombre,
		Owner:          owner,
		TotalCanciones: 0,
	}
}
