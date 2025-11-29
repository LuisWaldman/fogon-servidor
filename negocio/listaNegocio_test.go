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
	nombreUsuario := "bot1"
	nombreLista := "tomados"
	cancionServicio := servicios.NuevoCancionServicio(client)
	listaServicio := servicios.NuevoListaServicio(client)
	itemServicio := servicios.NuevoItemIndiceCancionServicio(client)
	negocioLista := NuevoListaNegocio(cancionServicio, listaServicio, itemServicio)
	err = negocioLista.NuevaListaForzarCreacion(nombreLista, nombreUsuario)
	assert.Nil(t, err, "Error al crear la lista: %v", err)
	lista, _ := negocioLista.GetLista(nombreLista, nombreUsuario)

	assert.Equal(t, 0, lista.TotalCanciones, "Se esperaban 0 canciones en la lista nueva")
	item := modelo.NewItemIndiceCancion("Canción de Prueba", "Banda de Prueba")
	negocioLista.AgregarCancionALista(nombreLista, nombreUsuario, item)
	lista, _ = negocioLista.GetLista(nombreLista, nombreUsuario)
	assert.Equal(t, 1, lista.TotalCanciones, "Se esperaban 1 canción en la lista después de agregarla")

	item2 := modelo.NewItemIndiceCancion("Canción de Prueba 2", "Banda de Prueba 2")
	negocioLista.AgregarCancionALista(nombreLista, nombreUsuario, item2)
	lista, _ = negocioLista.GetLista(nombreLista, nombreUsuario)
	assert.Equal(t, 2, lista.TotalCanciones, "Se esperaban 2 canciones en la lista después de agregarla")
	// Limpiar
	err = negocioLista.BorrarLista(nombreLista, nombreUsuario)
	assert.Nil(t, err, "Error al borrar la lista: %v", err)
	// Aquí iría la lógica para crear un nuevo usuario
	// GetCancionesPorUsuario GetListasPorUsuario
}
