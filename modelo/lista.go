package modelo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Lista struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nombre         string             `bson:"nombre" json:"nombre"`
	Owner          string             `bson:"owner" json:"owner"`
	TotalCanciones int                `bson:"total_canciones" json:"total_canciones"`
}

// UnmarshalBSON implementa bson.Unmarshaler para manejar ObjectID correctamente
func (l *Lista) UnmarshalBSON(data []byte) error {
	aux := &struct {
		ID     interface{} `bson:"_id,omitempty"`
		Nombre string      `bson:"nombre"`
		Owner  string      `bson:"owner"`
	}{}

	if err := bson.Unmarshal(data, aux); err != nil {
		return err
	}

	// Asignar los campos
	l.Nombre = aux.Nombre
	l.Owner = aux.Owner

	// Manejar diferentes formatos de ObjectID
	switch id := aux.ID.(type) {
	case primitive.ObjectID:
		l.ID = id
	case string:
		if objID, err := primitive.ObjectIDFromHex(id); err == nil {
			l.ID = objID
		}
	case map[string]interface{}:
		if oidStr, ok := id["$oid"].(string); ok {
			if objID, err := primitive.ObjectIDFromHex(oidStr); err == nil {
				l.ID = objID
			}
		}
	}

	return nil
}

func NuevaLista(nombre string, owner string) *Lista {
	return &Lista{
		Nombre:         nombre,
		Owner:          owner,
		TotalCanciones: 0,
	}
}
