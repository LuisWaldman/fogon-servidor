package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/LuisWaldman/fogon-servidor/aplicacion"
	"github.com/LuisWaldman/fogon-servidor/aplicacion/logueadores"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestControllerSesion(t *testing.T) {
	app := aplicacion.NuevoAplicacion()
	loginRepo := logueadores.NewLogeadorRepository()
	claves := []string{"VALIDA"}
	loginRepo.Add("TEST", logueadores.NewTesterLogeador(claves))
	newSocket := &aplicacion.MockSocket{}
	musico := aplicacion.NuevoMusico(newSocket, *loginRepo)
	app.AgregarMusico(musico)

	// Crear una sesión
	sesionID := "sesion_1"
	latitud := 12.34
	longitud := 56.78
	app.CrearSesion(musico, sesionID, latitud, longitud)
	micontroller := NuevoSesionesController(app) // Initialize the controller with the application

	// Crear un contexto de prueba
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Create a mock or appropriate context for 'c'

	c.Set("userID", musico.ID) // Set a mock user ID in the context
	micontroller.Get(c)

	assert.Equal(t, 200, w.Code, "Expected status code 200")

}
