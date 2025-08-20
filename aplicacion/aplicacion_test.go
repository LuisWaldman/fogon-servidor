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
	app.CrearSesion(musico, sesionID)

	// Verificar que la sesión se haya creado correctamente
	sesion, exists := app.sesiones[sesionID]
	if !exists {
		t.Errorf("La sesión %s no fue creada", sesionID)
		return
	}

	if sesion.nombre != sesionID {
		t.Errorf("Los datos de la sesión no coinciden: got %v, want %s, %f, %f", sesion, sesionID, latitud, longitud)
	}
	otromusico := NuevoMusico(newSocket, *loginRepo)
	app.AgregarMusico(otromusico)
	app.UnirseSesion(otromusico, sesionID)
	assert.True(t, newSocket.TieneMensaje("ensesion"), "No dio ensesion")
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
	assert.True(t, newSocket.TieneMensaje("sesionFailed"), "No dio ensesion")
}

func TestEliminaSesionesSinUsuarios(t *testing.T) {
	app := NuevoAplicacion()
	loginRepo := logueadores.NewLogeadorRepository()
	claves := []string{"VALIDA"}
	loginRepo.Add("TEST", logueadores.NewTesterLogeador(claves))
	newSocket := &MockSocket{}
	musico := NuevoMusico(newSocket, *loginRepo)
	app.AgregarMusico(musico)

	app.CrearSesion(musico, "sesion")
	assert.Equal(t, 1, len(app.sesiones), "Hay sesiones")
	musico.SalirSesion()
	app.ActualizarSesiones()
	assert.Equal(t, 0, len(app.sesiones), "Hay sesiones")

}

func TestEliminaSesionesYOtraPideUsuarios(t *testing.T) {
	app := NuevoAplicacion()
	loginRepo := logueadores.NewLogeadorRepository()
	claves := []string{"VALIDA"}
	loginRepo.Add("TEST", logueadores.NewTesterLogeador(claves))
	newSocket := &MockSocket{}
	newSocket2 := &MockSocket{}
	musico := NuevoMusico(newSocket, *loginRepo)
	app.AgregarMusico(musico)
	musico2 := NuevoMusico(newSocket2, *loginRepo)
	app.AgregarMusico(musico2)

	app.CrearSesion(musico, "sesion")
	assert.Equal(t, 1, len(app.sesiones), "Hay sesiones")
	app.UnirseSesion(musico2, "sesion")

	musico2.SalirSesion()
	app.ActualizarSesiones()
	app.QuitarMusico(musico2)
	buscaMusico, _ := app.BuscarMusicoPorID(musico.ID)
	user := buscaMusico.Sesion.GetUsuariosView()
	assert.Equal(t, 1, len(user), "Hay sesiones")

}
