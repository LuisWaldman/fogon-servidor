package aplicacion

import (
	"testing"

	"github.com/LuisWaldman/fogon-servidor/aplicacion/logueadores"
	"github.com/stretchr/testify/assert"
)

func TestSesionMensaje(t *testing.T) {
	// Crea un mock de socket.Socket usando testify/mock o una estructura personalizada
	loginRepo := logueadores.NewLogeadorRepository()
	claves := []string{"VALIDA"}
	loginRepo.Add("TEST", logueadores.NewTesterLogeador(claves))
	newSocket := &MockSocket{}
	newSocket2 := &MockSocket{}
	newMusico := NuevoMusico(newSocket, *loginRepo)
	newMusico2 := NuevoMusico(newSocket2, *loginRepo)

	newMusico.ID = 123  // Asigna un ID al usuario para la prueba
	newMusico2.ID = 423 // Asigna un ID al usuario para la prueba

	sesion := NuevaSesion("TestSession")
	newMusico.UnirseSesion(sesion)
	newMusico2.UnirseSesion(sesion)
	// Obtén el token emitido

	newMusico.MensajeSesion("Hola a todos")
	assert.Equal(t, "mensajesesion", newSocket.UltimoEmitted().Event, "No dio ensesion")
	mensaje, _ := newSocket.UltimoEmitted().Args[0].(string)
	assert.Equal(t, "Hola a todos", mensaje, "Mensaje no coincide")

	assert.Equal(t, "mensajesesion", newSocket2.UltimoEmitted().Event, "No dio ensesion")
	mensaje, _ = newSocket2.UltimoEmitted().Args[0].(string)
	assert.Equal(t, mensaje, "Hola a todos", "Mensaje no coincide")

}

func TestQuedaDirector(t *testing.T) {
	// Crea un mock de socket.Socket usando testify/mock o una estructura personalizada
	loginRepo := logueadores.NewLogeadorRepository()
	claves := []string{"VALIDA"}
	loginRepo.Add("TEST", logueadores.NewTesterLogeador(claves))
	newSocket := &MockSocket{}
	newSocket2 := &MockSocket{}
	newMusico := NuevoMusico(newSocket, *loginRepo)
	newMusico2 := NuevoMusico(newSocket2, *loginRepo)

	newMusico.ID = 123  // Asigna un ID al usuario para la prueba
	newMusico2.ID = 423 // Asigna un ID al usuario para la prueba

	sesion := NuevaSesion("TestSession")
	newMusico.UnirseSesion(sesion)
	newMusico2.UnirseSesion(sesion)
	newMusico.SalirSesion()

	assert.Equal(t, "actualizarusuarios", newSocket2.UltimoEmitted().Event, "No dio cambiarrolsesion")

}
