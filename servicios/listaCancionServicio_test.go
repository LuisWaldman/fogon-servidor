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
	assert.Nil(t, err, "Error al conectar a la base de datos: %v", err)
	servicio := NuevoListaServicio(client)
	lista, err := servicio.BuscarPorNombreYOwner(nombreLista, ownerTest)

	// Limpiar si existen las listas de prueba
	item := modelo.NewItemIndiceCancion("Canción de Prueba", "Banda de Prueba")
	listaCancion := modelo.NuevaListaCancion(lista.ID, item, 0)
	cancionServicio := NuevoListaCancionServicio(client)
	cancionServicio.AgregarCancion(listaCancion)

}
