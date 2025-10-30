package servicios

import (
	"testing"

	datos "github.com/LuisWaldman/fogon-servidor/datos"
	modelo "github.com/LuisWaldman/fogon-servidor/modelo"
	"github.com/stretchr/testify/assert"
)

func TestAgregarLista(t *testing.T) {
	nombreNuevo := "ListaRenombrada_" + RandString(8)
	ownerTest := "usuario_test"

	client, err := datos.ConnectDB()
	assert.Nil(t, err, "Error al conectar a la base de datos: %v", err)
	servicio := NuevoListaServicio(client)

	// Limpiar si existen las listas de prueba
	servicio.BorrarPorNombreYOwner(nombreNuevo, ownerTest)
	listaNueva := modelo.NuevaLista(nombreNuevo, ownerTest)
	cancionServicio := NuevoListaCancionServicio(client)
	cancionServicio.AgregarCancion()

}
