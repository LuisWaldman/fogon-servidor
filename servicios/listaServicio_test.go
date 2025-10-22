package servicios

import (
	"testing"

	datos "github.com/LuisWaldman/fogon-servidor/datos"
	modelo "github.com/LuisWaldman/fogon-servidor/modelo"
	"github.com/stretchr/testify/assert"
)

func TestCrearLista(t *testing.T) {
	nombreTest := "Mi Lista de Prueba " + RandString(4)
	ownerTest := "usuario_test_" + RandString(4)
	lista := modelo.NuevaLista(nombreTest, ownerTest)

	client, err := datos.ConnectDB()
	assert.Nil(t, err, "Error al crear base de datos: %v", err)
	servicio := NuevoListaServicio(client)

	// Limpiar si existe
	existente, _ := servicio.BuscarPorNombreYOwner(nombreTest, ownerTest)
	if existente != nil {
		servicio.BorrarPorID(existente.ID.Hex())
	}

	err = servicio.CrearLista(lista)
	assert.Nil(t, err, "Error al crear lista: %v", err)

	// Verificar que se puede encontrar
	encontrada, _ := servicio.BuscarPorNombreYOwner(nombreTest, ownerTest)
	if encontrada != nil {
		assert.Equal(t, nombreTest, encontrada.Nombre, "El nombre no se guardó correctamente")
		assert.Equal(t, ownerTest, encontrada.Owner, "El owner no se guardó correctamente")
		// Limpiar
		servicio.BorrarPorID(encontrada.ID.Hex())
	}
}

func TestTraerListaNoExistente(t *testing.T) {
	client, err := datos.ConnectDB()
	assert.Nil(t, err, "Error al crear base de datos: %v", err)
	servicio := NuevoListaServicio(client)

	lista, _ := servicio.BuscarPorNombreYOwner("LISTA_INEXISTENTE", "USUARIO_INEXISTENTE")
	assert.Nil(t, lista, "Lista no encontrada")
}

func TestCrearYBorrarLista(t *testing.T) {
	nombreLista := "test_lista_" + RandString(8)
	ownerLista := "test_owner_" + RandString(8)
	lista := modelo.NuevaLista(nombreLista, ownerLista)

	client, err := datos.ConnectDB()
	assert.Nil(t, err, "Error al crear base de datos: %v", err)
	servicio := NuevoListaServicio(client)

	// Limpiar completamente todas las listas de prueba con este nombre/owner
	err = servicio.BorrarPorNombreYOwner(nombreLista, ownerLista)
	assert.Nil(t, err, "Error al limpiar listas preexistentes: %v", err)

	// Verificar los datos de la lista antes de crear
	assert.Equal(t, nombreLista, lista.Nombre, "El nombre en la lista original no coincide")
	assert.Equal(t, ownerLista, lista.Owner, "El owner en la lista original no coincide")

	// Crear la lista
	err = servicio.CrearLista(lista)
	assert.Nil(t, err, "Error al crear lista: %v", err)

	// Verificar que se creó correctamente
	lista_b, err := servicio.BuscarPorNombreYOwner(nombreLista, ownerLista)
	assert.Nil(t, err, "Error al buscar lista: %v", err)
	assert.NotNil(t, lista_b, "Lista no encontrada")
	if lista_b != nil {
		assert.Equal(t, nombreLista, lista_b.Nombre, "El nombre de la lista no coincide")
		assert.Equal(t, ownerLista, lista_b.Owner, "El owner de la lista no coincide")

		// Debug: Imprimir el ID de la lista creada
		t.Logf("ID de la lista creada: %s", lista_b.ID.Hex())
	}

	// Borrar la lista usando el nuevo método directo
	err = servicio.BorrarPorNombreYOwner(nombreLista, ownerLista)
	assert.Nil(t, err, "Error al borrar lista: %v", err)

	// Verificar que se borró correctamente
	lista_c, _ := servicio.BuscarPorNombreYOwner(nombreLista, ownerLista)
	assert.Nil(t, lista_c, "Lista debería haber sido borrada pero aún existe")
}

func TestRenombrarLista(t *testing.T) {
	// Crear nombres únicos para evitar conflictos
	nombreOriginal := "ListaARenombrar_" + RandString(8)
	nombreNuevo := "ListaRenombrada_" + RandString(8)
	ownerTest := "usuario_test_" + RandString(8)

	client, err := datos.ConnectDB()
	assert.Nil(t, err, "Error al conectar a la base de datos: %v", err)
	servicio := NuevoListaServicio(client)

	// Limpiar si existen las listas de prueba
	servicio.BorrarPorNombreYOwner(nombreOriginal, ownerTest)
	servicio.BorrarPorNombreYOwner(nombreNuevo, ownerTest)

	// Crear la lista con nombre original
	lista := modelo.NuevaLista(nombreOriginal, ownerTest)
	err = servicio.CrearLista(lista)
	assert.Nil(t, err, "Error al crear lista original: %v", err)

	// Verificar que se creó correctamente
	listaOriginal, err := servicio.BuscarPorNombreYOwner(nombreOriginal, ownerTest)
	assert.Nil(t, err, "Error al buscar lista original: %v", err)
	assert.NotNil(t, listaOriginal, "Lista original no se creó correctamente")

	// Renombrar la lista
	err = servicio.RenombrarLista(listaOriginal.ID.Hex(), nombreNuevo)
	assert.Nil(t, err, "Error al renombrar lista: %v", err)

	// Verificar que existe la lista con el nuevo nombre
	listaRenombrada, err := servicio.BuscarPorNombreYOwner(nombreNuevo, ownerTest)
	assert.Nil(t, err, "Error al buscar lista renombrada: %v", err)
	assert.NotNil(t, listaRenombrada, "Lista renombrada no existe")

	// Verificar que no existe la lista con el nombre original
	listaAntigua, _ := servicio.BuscarPorNombreYOwner(nombreOriginal, ownerTest)
	assert.Nil(t, listaAntigua, "Lista original no debería existir después del renombrado")

	// Verificar que los datos se mantuvieron, excepto el nombre
	if listaRenombrada != nil {
		assert.Equal(t, nombreNuevo, listaRenombrada.Nombre, "El nuevo nombre no se actualizó correctamente")
		assert.Equal(t, ownerTest, listaRenombrada.Owner, "El owner no se mantuvo correctamente")
		assert.Equal(t, listaOriginal.ID, listaRenombrada.ID, "El ID no debería cambiar al renombrar")
	}

	// Limpiar
	if listaRenombrada != nil {
		servicio.BorrarPorID(listaRenombrada.ID.Hex())
	}
}
