package modelo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListaCancion struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ListaID           primitive.ObjectID `bson:"listaId" json:"listaId"`
	ItemIndiceCancion ItemIndiceCancion  `bson:"itemIndiceCancion" json:"itemIndiceCancion"`
	Orden             int                `bson:"orden" json:"orden"`
	FechaAgregada     time.Time          `bson:"fechaAgregada" json:"fechaAgregada"`
}

func NuevaListaCancion(listaID primitive.ObjectID, item *ItemIndiceCancion, orden int) *ListaCancion {
	return &ListaCancion{
		ListaID:           listaID,
		ItemIndiceCancion: *item,
		Orden:             orden,
		FechaAgregada:     time.Now(),
	}
}
