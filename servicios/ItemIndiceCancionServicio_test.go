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
	lista, err := listaservicio.BuscarPorNombreYOwner(nombreLista, ownerTest)
	if lista == nil {
		listaservicio.CrearLista(nombreLista, ownerTest)
		lista, err = listaservicio.BuscarPorNombreYOwner(nombreLista, ownerTest)
	}
	assert.Nil(t, err, "Error al buscar la lista: %v", err)
	assert.NotNil(t, lista, "Error al buscar la lista: %v", err)

	// Limpiar si existen las listas de prueba
	item := modelo.NewItemIndiceCancion("Canción de Prueba", "Banda de Prueba")
	item.ListaID = lista.ID
	itemListaServicio.AgregarCancion(item)

}
