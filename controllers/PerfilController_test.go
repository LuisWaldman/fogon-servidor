package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http/httptest"
	"testing"

	"github.com/LuisWaldman/fogon-servidor/aplicacion"
	"github.com/LuisWaldman/fogon-servidor/aplicacion/logueadores"
	"github.com/LuisWaldman/fogon-servidor/datos"
	"github.com/LuisWaldman/fogon-servidor/modelo"
	"github.com/LuisWaldman/fogon-servidor/servicios"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func MockJsonPost(c *gin.Context, content interface{}) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")
	jsonBytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
}

func textoalazar() string {
	// Genera un texto aleatorio de 10 caracteres
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 10)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)

}

func TestPerfilControllerSesion(t *testing.T) {
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

	client, err := datos.ConnectDB()
	assert.Nil(t, err, "Error al crear base de datos: %v", err)
	servicio := servicios.NuevoPerfilServicio(client)
	micontroller := NuevoPerfilController(servicio, app) // Initialize the controller with the application

	// Crear un contexto de prueba
	w := httptest.NewRecorder()
	context_POST, _ := gin.CreateTestContext(w)
	context_POST.Set("userID", musico.ID) // Set a mock user ID in the context

	perfil := modelo.Perfil{}
	perfil.Nombre = textoalazar()
	perfil.Descripcion = textoalazar()
	perfil.Imagen = textoalazar()
	context_POST.Request = httptest.NewRequest("POST", "/perfil", nil)
	context_POST.Request.Header.Set("Content-Type", "application/json")
	MockJsonPost(context_POST, perfil) // Mock the JSON body for the POST request
	micontroller.Post(context_POST)
	assert.Equal(t, 201, w.Code, "Expected status code 201")

	wGET := httptest.NewRecorder()
	context_GET, _ := gin.CreateTestContext(wGET)
	context_GET.Set("userID", musico.ID) // Set a mock user ID in the context
	micontroller.Get(context_GET)
	assert.Equal(t, 200, wGET.Code, "Expected status code 200")
	perfilResponse := modelo.Perfil{}
	err = json.Unmarshal(wGET.Body.Bytes(), &perfilResponse)
	assert.Nil(t, err, "Error al deserializar la respuesta: %v", err)
	assert.Equal(t, perfil.Nombre, perfilResponse.Nombre, "Expected profile name to match")
	assert.Equal(t, perfil.Descripcion, perfilResponse.Descripcion, "Expected profile description to match")
	assert.Equal(t, perfil.Imagen, perfilResponse.Imagen, "Expected profile image to match")
	assert.Equal(t, musico.Usuario, perfilResponse.Usuario, "Expected profile user to match")

}
