package logueadores

import (
	"testing"

	"github.com/LuisWaldman/fogon-servidor/db"
	"github.com/LuisWaldman/fogon-servidor/servicios"
	"github.com/stretchr/testify/assert"
)

func TestNewUserPassLogeador(t *testing.T) {

	client, err := db.ConnectDB()
	assert.Nil(t, err, "Error al crear base de datos: %v", err)
	servicio := servicios.NuevoUsuarioServicio(client)
	logeador := NewUserPassLogeador(servicio)
	if logeador == nil {
		t.Fatal("Expected non-nil UserPassLogeador")
	}
	servicio.BorrarPorUsuario("BORRADO1")
	ret := logeador.Login("BORRADO1", "clave1")
	assert.True(t, ret, "La primera vez debería ser true, ya que el usuario no existe y se crea")
	ret = logeador.Login("BORRADO1", "clave2")
	assert.False(t, ret, "La segunda vez debería ser false, ya que el usuario ya existe con una clave diferente")

}
