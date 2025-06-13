package servicios

import (
	"testing"

	db "github.com/LuisWaldman/fogon-servidor/db"
	modelo "github.com/LuisWaldman/fogon-servidor/modelo"
	"github.com/stretchr/testify/assert"
)

// Suponiendo que tienes una estructura Usuario y funciones CrearUsuario, BuscarPorUsuario, BorrarPorUsuario similares a las del ejemplo

func TestCrearPerfil(t *testing.T) {
	perfil := &modelo.Perfil{}
	perfil.Usuario = "servicio1"
	perfil.Imagen = "asd"
	perfil.Nombre = "sdf"
	perfil.Instrumento = "dfgdg"
	perfil.Descripcion = "fdgfdgdf"

	client, err := db.ConnectDB()
	assert.Nil(t, err, "Error al crear base de datos: %v", err)
	servicio := NuevoPerfilServicio(client)
	err = servicio.CrearPerfil(*perfil)
	assert.Nil(t, err, "Error al crear perfil: %v", err)

}

func TestTraerPerfilNoExistente(t *testing.T) {

	client, err := db.ConnectDB()
	assert.Nil(t, err, "Error al crear base de datos: %v", err)
	servicio := NuevoPerfilServicio(client)

	perfil, _ := servicio.BuscarPorUsuario("INEXISTENTE")
	assert.Nil(t, perfil, "Usuario no encontrado")
}

func TestCrearYBorrarPerfil(t *testing.T) {
	perfil := &modelo.Perfil{}

	usuario := "test_user_" + RandString(8)
	perfil.Usuario = usuario
	perfil.Imagen = "asd"
	perfil.Nombre = "sdf"
	perfil.Instrumento = "dfgdg"
	perfil.Descripcion = "fdgfdgdf"

	client, err := db.ConnectDB()
	assert.Nil(t, err, "Error al crear base de datos: %v", err)
	servicio := NuevoPerfilServicio(client)

	elperfil, _ := servicio.BuscarPorUsuario(usuario)
	if elperfil != nil {
		err = servicio.BorrarPorUsuario(usuario)
		assert.Nil(t, err, "Error al borrar perfil existente: %v", err)
	}

	err = servicio.CrearPerfil(*perfil)
	assert.Nil(t, err, "Error al crear perfil: %v", err)
	perfil_b, _ := servicio.BuscarPorUsuario(usuario)
	assert.NotNil(t, perfil_b, "Perfil no encontrado")
	err = servicio.BorrarPorUsuario(usuario)
	assert.Nil(t, err, "Error al borrar perfil: %v", err)
	perfil_b, _ = servicio.BuscarPorUsuario(usuario)
	assert.Nil(t, perfil_b, "Perfil debería haber sido borrado pero aún existe")

}
