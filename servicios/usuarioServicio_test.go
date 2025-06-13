package servicios

import (
	"testing"

	db "github.com/LuisWaldman/fogon-servidor/db"
	modelo "github.com/LuisWaldman/fogon-servidor/modelo"
	"github.com/stretchr/testify/assert"
)

// Suponiendo que tienes una estructura Usuario y funciones CrearUsuario, BuscarPorUsuario, BorrarPorUsuario similares a las del ejemplo

func TestCrearUsuarioServicio(t *testing.T) {
	usuario := &modelo.Usuario{}
	usuario.Modologin = "USERPASS"
	usuario.Usuario = "servicio1"
	usuario.Clave = "clave1"
	client, err := db.ConnectDB()
	assert.Nil(t, err, "Error al crear base de datos: %v", err)
	servicio := NuevoUsuarioServicio(client)
	err = servicio.CrearUsuario(*usuario)
	assert.Nil(t, err, "Error al crear usuario: %v", err)

}

func TestObtenerUsuario(t *testing.T) {
	// Crea un mock de socket.Socket usando testify/mock o una estructura personalizada
	// Implementa los métodos necesarios de socket.Socket en mockSocket si es necesario

	client, err := db.ConnectDB()
	assert.Nil(t, err, "Error al crear base de datos: %v", err)
	servicio := NuevoUsuarioServicio(client)

	usuario, _ := servicio.BuscarPorUsuario("servicio1")
	assert.NotNil(t, usuario, "Usuario no encontrado")
	assert.True(t, usuario.Encontrado, "El usuario debería existir")
	assert.Equal(t, "USERPASS", usuario.Modologin, "El modo de login no coincide")
}
func TestObtenerUsuarioInexistente(t *testing.T) {
	// Crea un mock de socket.Socket usando testify/mock o una estructura personalizada
	// Implementa los métodos necesarios de socket.Socket en mockSocket si es necesario

	client, err := db.ConnectDB()
	assert.Nil(t, err, "Error al crear base de datos: %v", err)
	servicio := NuevoUsuarioServicio(client)

	usuario, _ := servicio.BuscarPorUsuario("INEXISTENTE")
	assert.NotNil(t, usuario, "Usuario no encontrado")
	assert.False(t, usuario.Encontrado, "El usuario no debería existir")
}
func TestCreaYBorra(t *testing.T) {

	client, errDB := db.ConnectDB()
	assert.Nil(t, errDB, "Error al crear base de datos: %v", errDB)
	servicio := NuevoUsuarioServicio(client)
	nombre := "test_user_" + RandString(8)
	usuario, _ := servicio.BuscarPorUsuario(nombre)
	if usuario.Encontrado {
		servicio.BorrarPorUsuario(nombre)
		usuarioborrado, _ := servicio.BuscarPorUsuario(nombre)
		assert.False(t, usuarioborrado.Encontrado, "El usuario no debería existir")

	}

	usuarioNuevo := &modelo.Usuario{}
	usuarioNuevo.Modologin = "USERPASS"
	usuarioNuevo.Usuario = nombre
	usuarioNuevo.Clave = "par_2"
	err := servicio.CrearUsuario(*usuarioNuevo)
	assert.Nil(t, err, "Error al crear usuario: %v", err)

	usuariocreado, _ := servicio.BuscarPorUsuario("pero")
	assert.NotNil(t, usuariocreado, "Usuario no encontrado")
	assert.True(t, usuariocreado.Encontrado, "El usuario debería existir")
	assert.Equal(t, "USERPASS", usuariocreado.Modologin, "El modo de login no coincide")

	servicio.BorrarPorUsuario(nombre)
	usuarioborrado, _ := servicio.BuscarPorUsuario(nombre)
	assert.False(t, usuarioborrado.Encontrado, "El usuario no debería existir")
}

func RandString(i int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, i)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}
