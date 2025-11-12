package servicios

import (
	"testing"

	datos "github.com/LuisWaldman/fogon-servidor/datos"
	modelo "github.com/LuisWaldman/fogon-servidor/modelo"
	"github.com/stretchr/testify/assert"
)

func TestAgregarLista(t *testing.T) {
	nombreLista := "listaTest"
	ownerTest := "usuario_test"

	client, err := datos.ConnectDB()
	itemListaServicio := NuevoItemIndiceCancionServicio(client)
	listaservicio := NuevoListaServicio(client)
	assert.Nil(t, err, "Error al conectar a la base de datos: %v", err)

	lista, _ := listaservicio.BuscarPorNombreYOwner(nombreLista, ownerTest)
	if lista != nil {
		itemListaServicio.BorrarPorListaID(lista.ID.Hex())
		listaservicio.BorrarPorID(lista.ID.Hex())
	}
	listaservicio.CrearLista(nombreLista, ownerTest)
	lista, err = listaservicio.BuscarPorNombreYOwner(nombreLista, ownerTest)

	assert.Nil(t, err, "Error al buscar la lista: %v", err)
	assert.NotNil(t, lista, "Error al buscar la lista: %v", err)

	// Limpiar si existen las listas de prueba
	item := modelo.NewItemIndiceCancion("Canción de Prueba", "Banda de Prueba")
	item.ListaID = lista.ID
	itemListaServicio.AgregarCancion(item)

	item2 := modelo.NewItemIndiceCancion("Canción de Prueba", "Banda de Prueba")
	item2.ListaID = lista.ID
	itemListaServicio.AgregarCancion(item2)

	canciones := itemListaServicio.GetCancionesPorListaID(lista.ID)
	assert.NotNil(t, canciones, "Error al obtener canciones por lista ID: %v", err)
	assert.Len(t, canciones, 2, "Se esperaba 2 canciones en la lista")
	// Limpiar
	itemListaServicio.BorrarPorListaID(lista.ID.Hex())

}
