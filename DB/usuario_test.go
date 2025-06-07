package DB

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrearUsuario(t *testing.T) {
	// Crea un mock de socket.Socket usando testify/mock o una estructura personalizada
	// Implementa los métodos necesarios de socket.Socket en mockSocket si es necesario

	usuario := &usuarioDB{}
	usuario.Modologin = "USERPASS"
	usuario.Usuario = "pero"
	usuario.Clave = "par_2"
	err := crear_usuario(*usuario)
	assert.Nil(t, err, "Error al crear usuario: %v", err)
}
func TestCrearUsuario2(t *testing.T) {
	// Crea un mock de socket.Socket usando testify/mock o una estructura personalizada
	// Implementa los métodos necesarios de socket.Socket en mockSocket si es necesario

	usuario := &usuarioDB{}
	usuario.Modologin = "USERPASS"
	usuario.Usuario = "OTROUSER"
	usuario.Clave = "par_2"
	err := crear_usuario(*usuario)
	assert.Nil(t, err, "Error al crear usuario: %v", err)
}
func TestObtenerUsuario(t *testing.T) {
	// Crea un mock de socket.Socket usando testify/mock o una estructura personalizada
	// Implementa los métodos necesarios de socket.Socket en mockSocket si es necesario

	usuario, _ := buscarxusuario("pero")
	assert.NotNil(t, usuario, "Usuario no encontrado")
	assert.True(t, usuario.Encontrado, "El usuario debería existir")
	assert.Equal(t, "USERPASS", usuario.Modologin, "El modo de login no coincide")
}
func TestObtenerUsuarioInexistente(t *testing.T) {
	// Crea un mock de socket.Socket usando testify/mock o una estructura personalizada
	// Implementa los métodos necesarios de socket.Socket en mockSocket si es necesario

	usuario, _ := buscarxusuario("INEXISTENTE")
	assert.NotNil(t, usuario, "Usuario no encontrado")
	assert.False(t, usuario.Encontrado, "El usuario no debería existir")
}
func TestCreaYBorra(t *testing.T) {

	nombre := "test_user_" + RandString(8)
	usuario, _ := buscarxusuario(nombre)
	if usuario.Encontrado {
		borrarxusuario(nombre)
		usuarioborrado, _ := buscarxusuario(nombre)
		assert.False(t, usuarioborrado.Encontrado, "El usuario no debería existir")

	}

	usuarioNuevo := &usuarioDB{}
	usuarioNuevo.Modologin = "USERPASS"
	usuarioNuevo.Usuario = nombre
	usuarioNuevo.Clave = "par_2"
	err := crear_usuario(*usuarioNuevo)
	assert.Nil(t, err, "Error al crear usuario: %v", err)

	usuariocreado, _ := buscarxusuario("pero")
	assert.NotNil(t, usuariocreado, "Usuario no encontrado")
	assert.True(t, usuariocreado.Encontrado, "El usuario debería existir")
	assert.Equal(t, "USERPASS", usuariocreado.Modologin, "El modo de login no coincide")

	borrarxusuario(nombre)
	usuarioborrado, _ := buscarxusuario(nombre)
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
