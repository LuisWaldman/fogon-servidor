package negocio

import (
	"testing"

	"github.com/LuisWaldman/fogon-servidor/datos"
	"github.com/LuisWaldman/fogon-servidor/servicios"
	"github.com/stretchr/testify/assert"
)

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
	canciones, err := negocioUsuario.GetCancionesPorUsuario(user.Usuario)
	assert.Nil(t, err, "Error al obtener canciones: %v", err)
	assert.Equal(t, 0, len(canciones), "Se esperaban 0 canciones para el usuario")
	//listas, err := negocioUsuario.GetListasPorUsuario(user.Usuario)
	//assert.Nil(t, err, "Error al obtener listas: %v", err)
	//assert.Equal(t, 0, len(listas), "Se esperaban 0 listas para el usuario")

	// Aquí iría la lógica para crear un nuevo usuario
	// GetCancionesPorUsuario GetListasPorUsuario
}
