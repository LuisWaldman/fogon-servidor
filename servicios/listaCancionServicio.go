package servicios

import (
	"context"
	"log"

	modelo "github.com/LuisWaldman/fogon-servidor/modelo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ListaCancionServicio struct {
	db         *mongo.Client
	collection string
}

func NuevoListaCancionServicio(db *mongo.Client) *ListaCancionServicio {
	return &ListaCancionServicio{
		db:         db,
		collection: "listaCanciones",
	}
}

func (s *ListaCancionServicio) AgregarCancion(listaCancion *modelo.ListaCancion) error {
	col := s.db.Database(database).Collection(s.collection)

	// Si no se especifica el orden, asignar el siguiente número
	if listaCancion.Orden == 0 {
		ultimoOrden, err := s.ObtenerUltimoOrden(listaCancion.ListaID)
		if err != nil {
			return err
		}
		listaCancion.Orden = ultimoOrden + 1
	}

	inserta, err := col.InsertOne(context.TODO(), listaCancion)
	if err != nil {
		log.Println("Error agregando canción a lista", "err", err)
		return err
	}
	log.Println("Canción agregada a lista", inserta)
	return nil
}

func (s *ListaCancionServicio) ObtenerUltimoOrden(listaID primitive.ObjectID) (int, error) {
	col := s.db.Database(database).Collection(s.collection)

	// Buscar el documento con mayor orden
	opts := options.FindOne().SetSort(bson.D{{Key: "orden", Value: -1}})
	var ultimaCancion modelo.ListaCancion

	err := col.FindOne(context.TODO(), bson.M{"listaId": listaID}, opts).Decode(&ultimaCancion)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil // No hay canciones, empezar desde 0
		}
		return 0, err
	}

	return ultimaCancion.Orden, nil
}

func (s *ListaCancionServicio) ObtenerCancionesPorLista(listaID primitive.ObjectID) ([]modelo.ListaCancion, error) {
	col := s.db.Database(database).Collection(s.collection)

	// Ordenar por el campo 'orden'
	opts := options.Find().SetSort(bson.D{{Key: "orden", Value: 1}})
	cursor, err := col.Find(context.TODO(), bson.M{"listaId": listaID}, opts)
	if err != nil {
		log.Println("Error obteniendo canciones de lista", "err", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var canciones []modelo.ListaCancion
	for cursor.Next(context.TODO()) {
		var cancion modelo.ListaCancion
		if err := cursor.Decode(&cancion); err != nil {
			log.Println("Error decodificando canción de lista", "err", err)
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

func (s *ListaCancionServicio) BuscarPorID(id string) (*modelo.ListaCancion, error) {
	col := s.db.Database(database).Collection(s.collection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var listaCancion modelo.ListaCancion
	err = col.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&listaCancion)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &listaCancion, nil
}

func (s *ListaCancionServicio) EliminarCancion(id string) error {
	// Primero obtener la canción para conocer su lista y orden
	cancion, err := s.BuscarPorID(id)
	if err != nil || cancion == nil {
		return err
	}

	col := s.db.Database(database).Collection(s.collection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = col.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		log.Println("Error eliminando canción de lista", "err", err)
		return err
	}

	// Reordenar las canciones posteriores
	err = s.ReordenarTrasEliminacion(cancion.ListaID, cancion.Orden)
	if err != nil {
		log.Println("Error reordenando tras eliminación", "err", err)
		// No retornamos error para no fallar la eliminación
	}

	return nil
}

func (s *ListaCancionServicio) ReordenarTrasEliminacion(listaID primitive.ObjectID, ordenEliminado int) error {
	col := s.db.Database(database).Collection(s.collection)

	// Decrementar en 1 el orden de todas las canciones que tengan orden mayor al eliminado
	filtro := bson.M{
		"listaId": listaID,
		"orden":   bson.M{"$gt": ordenEliminado},
	}

	actualizacion := bson.M{"$inc": bson.M{"orden": -1}}

	_, err := col.UpdateMany(context.TODO(), filtro, actualizacion)
	if err != nil {
		log.Println("Error reordenando canciones", "err", err)
		return err
	}

	return nil
}

func (s *ListaCancionServicio) CambiarOrden(id string, nuevoOrden int) error {
	cancion, err := s.BuscarPorID(id)
	if err != nil || cancion == nil {
		return err
	}

	ordenAnterior := cancion.Orden
	if ordenAnterior == nuevoOrden {
		return nil // No hay cambio
	}

	col := s.db.Database(database).Collection(s.collection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Ajustar los órdenes de las otras canciones
	if nuevoOrden > ordenAnterior {
		// Mover hacia abajo: decrementar orden de canciones entre ordenAnterior+1 y nuevoOrden
		filtro := bson.M{
			"listaId": cancion.ListaID,
			"orden":   bson.M{"$gt": ordenAnterior, "$lte": nuevoOrden},
		}
		actualizacion := bson.M{"$inc": bson.M{"orden": -1}}
		_, err = col.UpdateMany(context.TODO(), filtro, actualizacion)
		if err != nil {
			return err
		}
	} else {
		// Mover hacia arriba: incrementar orden de canciones entre nuevoOrden y ordenAnterior-1
		filtro := bson.M{
			"listaId": cancion.ListaID,
			"orden":   bson.M{"$gte": nuevoOrden, "$lt": ordenAnterior},
		}
		actualizacion := bson.M{"$inc": bson.M{"orden": 1}}
		_, err = col.UpdateMany(context.TODO(), filtro, actualizacion)
		if err != nil {
			return err
		}
	}

	// Actualizar el orden de la canción específica
	_, err = col.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": bson.M{"orden": nuevoOrden}})
	if err != nil {
		log.Println("Error actualizando orden de canción", "err", err)
		return err
	}

	return nil
}

func (s *ListaCancionServicio) ActualizarNotas(id string, notas string) error {
	col := s.db.Database(database).Collection(s.collection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = col.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": bson.M{"notas": notas}})
	if err != nil {
		log.Println("Error actualizando notas de canción", "err", err)
		return err
	}

	return nil
}

func (s *ListaCancionServicio) BorrarPorLista(listaID primitive.ObjectID) error {
	col := s.db.Database(database).Collection(s.collection)
	_, err := col.DeleteMany(context.TODO(), bson.M{"listaId": listaID})
	if err != nil {
		log.Println("Error borrando canciones de lista", "err", err)
		return err
	}
	return nil
}
