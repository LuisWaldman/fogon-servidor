package aplicacion

import (
	"testing"

	"github.com/LuisWaldman/fogon-servidor/aplicacion/logueadores"
	"github.com/stretchr/testify/assert"
)

func TestCreoSesion(t *testing.T) {
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

	if sesion.nombre != sesionID || sesion.latitud != latitud || sesion.longitud != longitud {
		t.Errorf("Los datos de la sesión no coinciden: got %v, want %s, %f, %f", sesion, sesionID, latitud, longitud)
	}

	assert.Equal(t, "ensesion", newSocket.UltimoEmitted().Event, "No dio ensesion")

	otromusico := NuevoMusico(newSocket, *loginRepo)
	app.AgregarMusico(otromusico)
	app.UnirseSesion(otromusico, sesionID)
	assert.Equal(t, "ensesion", newSocket.UltimoEmitted().Event, "No dio ensesion")
}

func TestSeUneASesionInexistente(t *testing.T) {
	app := NuevoAplicacion()
	loginRepo := logueadores.NewLogeadorRepository()
	claves := []string{"VALIDA"}
	loginRepo.Add("TEST", logueadores.NewTesterLogeador(claves))
	newSocket := &MockSocket{}
	musico := NuevoMusico(newSocket, *loginRepo)
	app.AgregarMusico(musico)

	app.UnirseSesion(musico, "sesion_inexistente")
	assert.Equal(t, "sesionFailed", newSocket.UltimoEmitted().Event, "No dio sesionFailed")
}

func TestEliminaSesionesSinUsuarios(t *testing.T) {
	app := NuevoAplicacion()
	loginRepo := logueadores.NewLogeadorRepository()
	claves := []string{"VALIDA"}
	loginRepo.Add("TEST", logueadores.NewTesterLogeador(claves))
	newSocket := &MockSocket{}
	musico := NuevoMusico(newSocket, *loginRepo)
	app.AgregarMusico(musico)

	app.CrearSesion(musico, "sesion", 0, 3.14)
	assert.Equal(t, 1, len(app.sesiones), "Hay sesiones")
	musico.SalirSesion()
	app.ActualizarSesiones()
	assert.Equal(t, 0, len(app.sesiones), "Hay sesiones")

}
