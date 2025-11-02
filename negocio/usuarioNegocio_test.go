package negocio

import (
	"testing"

	"github.com/LuisWaldman/fogon-servidor/datos"
	"github.com/LuisWaldman/fogon-servidor/modelo"
	"github.com/LuisWaldman/fogon-servidor/servicios"
	"github.com/stretchr/testify/assert"
)

func TestUsuarioNuevo_Crear(t *testing.T) {
	client, err := datos.ConnectDB()
	assert.Nil(t, err, "Error al conectar a la base de datos: %v", err)
	nombreUsuario := "NuevoUsuario"
	usuarioServicio := servicios.NuevoUsuarioServicio(client)
	cancionServicio := servicios.NuevoCancionServicio(client)
	listaServicio := servicios.NuevoListaServicio(client)
	itemServicio := servicios.NuevoItemIndiceCancionServicio(client)

	// CREO NEGOCIO USUARIO
	negocioUsuario := NuevoUsuarioNegocio(usuarioServicio, cancionServicio, listaServicio, itemServicio)
	negocioUsuario.BorrarPorUsuario(nombreUsuario)
	negocioUsuario.CrearUsuario(nombreUsuario)
	user, err := negocioUsuario.BuscarPorUsuario(nombreUsuario)
	assert.Nil(t, err, "Error al buscar usuario: %v", err)
	assert.NotNil(t, user, "Usuario no encontrado")
	assert.Equal(t, nombreUsuario, user.Usuario, "Nombre de usuario incorrecto")
	assert.Equal(t, "", user.Clave, "Clave de usuario incorrecta")
	listas, err := negocioUsuario.GetListasPorUsuario(user.Usuario)
	assert.Nil(t, err, "Error al obtener listas: %v", err)
	assert.Equal(t, 0, len(listas), "Se esperaban 0 listas para el usuario")
}

func TestUsuarioNuevo_CrearYBuscar(t *testing.T) {
	client, err := datos.ConnectDB()
	assert.Nil(t, err, "Error al conectar a la base de datos: %v", err)
	nombreUsuario := "NuevoUsuario"
	usuarioServicio := servicios.NuevoUsuarioServicio(client)
	cancionServicio := servicios.NuevoCancionServicio(client)
	listaServicio := servicios.NuevoListaServicio(client)
	itemServicio := servicios.NuevoItemIndiceCancionServicio(client)

	// CREO NEGOCIO USUARIO
	negocioUsuario := NuevoUsuarioNegocio(usuarioServicio, cancionServicio, listaServicio, itemServicio)
	negocioUsuario.BorrarPorUsuario(nombreUsuario)
	negocioUsuario.CrearUsuario(nombreUsuario)
	user, err := negocioUsuario.BuscarPorUsuario(nombreUsuario)
	assert.Nil(t, err, "Error al buscar usuario: %v", err)
	assert.NotNil(t, user, "Usuario no encontrado")
	assert.Equal(t, nombreUsuario, user.Usuario, "Nombre de usuario incorrecto")
	assert.Equal(t, "", user.Clave, "Clave de usuario incorrecta")
	listas, err := negocioUsuario.GetListasPorUsuario(user.Usuario)
	assert.Nil(t, err, "Error al obtener listas: %v", err)
	assert.Equal(t, 0, len(listas), "Se esperaban 0 listas para el usuario")
}

func TestAgregarALista(t *testing.T) {
	client, err := datos.ConnectDB()
	assert.Nil(t, err, "Error al conectar a la base de datos: %v", err)
	nombreUsuario := "NuevoUsuario"
	usuarioServicio := servicios.NuevoUsuarioServicio(client)
	cancionServicio := servicios.NuevoCancionServicio(client)
	listaServicio := servicios.NuevoListaServicio(client)
	itemServicio := servicios.NuevoItemIndiceCancionServicio(client)

	// CREO LISTA PARA USUARIO
	nombreLista := "ListaParaTest11"
	negocioUsuario := NuevoUsuarioNegocio(usuarioServicio, cancionServicio, listaServicio, itemServicio)
	negocioUsuario.BorrarPorUsuario(nombreUsuario)
	negocioUsuario.CrearUsuario(nombreUsuario)
	negocioUsuario.AgregarLista(nombreLista, nombreUsuario)

	user, err := negocioUsuario.BuscarPorUsuario(nombreUsuario)
	assert.Nil(t, err, "Error al buscar usuario: %v", err)
	assert.NotNil(t, user, "Usuario no encontrado")
	assert.Equal(t, nombreUsuario, user.Usuario, "Nombre de usuario incorrecto")
	assert.Equal(t, 1, len(user.Listas), "Se esperaban 1 lista para el usuario")
	assert.Equal(t, nombreLista, user.Listas[0], "Nombre de la lista incorrecto")
	// Limpiar
	negocioUsuario.BorrarPorUsuario(nombreUsuario)

}

func TestAgregarCancion(t *testing.T) {
	client, err := datos.ConnectDB()
	assert.Nil(t, err, "Error al conectar a la base de datos: %v", err)
	nombreUsuario := "NuevoUsuario"
	usuarioServicio := servicios.NuevoUsuarioServicio(client)
	cancionServicio := servicios.NuevoCancionServicio(client)
	listaServicio := servicios.NuevoListaServicio(client)
	itemServicio := servicios.NuevoItemIndiceCancionServicio(client)

	// CREO LISTA PARA USUARIO
	negocioUsuario := NuevoUsuarioNegocio(usuarioServicio, cancionServicio, listaServicio, itemServicio)
	negocioUsuario.BorrarPorUsuario(nombreUsuario)
	negocioUsuario.CrearUsuario(nombreUsuario)
	nuevaCancion := modelo.NuevaCancion("nombrecancion", nombreUsuario)
	negocioUsuario.AgregarCancion(nombreUsuario, nuevaCancion)

	user, err := negocioUsuario.BuscarPorUsuario(nombreUsuario)
	assert.Nil(t, err, "Error al buscar usuario: %v", err)
	assert.NotNil(t, user, "Usuario no encontrado")
	cancionesUsuario := negocioUsuario.GetCancionesPorUsuario(nombreUsuario)
	assert.Equal(t, nombreUsuario, user.Usuario, "Nombre de usuario incorrecto")
	assert.Equal(t, 1, len(cancionesUsuario), "Se esperaban 1 cancion para el usuario")
	// Limpiar
	negocioUsuario.BorrarPorUsuario(nombreUsuario)

}
