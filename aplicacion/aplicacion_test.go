package aplicacion

import (
	"testing"

	"github.com/LuisWaldman/fogon-servidor/aplicacion/logueadores"
	"github.com/stretchr/testify/assert"
)

func CreoSesion(t *testing.T) {
	app := NuevoAplicacion()
	loginRepo := logueadores.NewLogeadorRepository()
	claves := []string{"VALIDA"}
	loginRepo.Add("TEST", logueadores.NewTesterLogeador(claves))
	newSocket := &MockSocket{}
	musico := NuevoMusico(newSocket, *loginRepo)
	app.AgregarMusico(musico)

	// Crear una sesión
	sesionID := "sesion_1"
	latitud := 12.34
	longitud := 56.78
	app.CrearSesion(musico, sesionID, latitud, longitud)

	// Verificar que la sesión se haya creado correctamente
	sesion, exists := app.sesiones[sesionID]
	if !exists {
		t.Errorf("La sesión %s no fue creada", sesionID)
		return
	}

	if sesion.sesion != sesionID || sesion.latitud != latitud || sesion.longitud != longitud {
		t.Errorf("Los datos de la sesión no coinciden: got %v, want %s, %f, %f", sesion, sesionID, latitud, longitud)
	}

	assert.Equal(t, "ensesion", newSocket.UltimoEmitted().Event, "No dio ensesion")
}
