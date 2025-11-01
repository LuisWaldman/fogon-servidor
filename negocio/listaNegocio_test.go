package negocio

import (
	"testing"

	"github.com/LuisWaldman/fogon-servidor/datos"
	"github.com/LuisWaldman/fogon-servidor/modelo"
	"github.com/LuisWaldman/fogon-servidor/servicios"
	"github.com/stretchr/testify/assert"
)

func TestLista(t *testing.T) {
	client, err := datos.ConnectDB()
	assert.Nil(t, err, "Error al conectar a la base de datos: %v", err)
	nombreUsuario := "NuevoUsuario"
	nombreLista := "NuevaLista"
	cancionServicio := servicios.NuevoCancionServicio(client)
	listaServicio := servicios.NuevoListaServicio(client)
	itemServicio := servicios.NuevoItemIndiceCancionServicio(client)
	negocioLista := NuevoListaNegocio(cancionServicio, listaServicio, itemServicio)
	negocioLista.NuevaListaForzarCreacion(nombreLista, nombreUsuario)
	lista, _ := negocioLista.GetLista(nombreLista, nombreUsuario)

	assert.Equal(t, 0, lista.TotalCanciones, "Se esperaban 0 canciones en la lista nueva")
	item := modelo.NewItemIndiceCancion("Canción de Prueba", "Banda de Prueba")
	negocioLista.AgregarCancionALista(nombreLista, nombreUsuario, item)
	lista, _ = negocioLista.GetLista(nombreLista, nombreUsuario)
	assert.Equal(t, 1, lista.TotalCanciones, "Se esperaban 1 canción en la lista después de agregarla")

	// Aquí iría la lógica para crear un nuevo usuario
	// GetCancionesPorUsuario GetListasPorUsuario
}
