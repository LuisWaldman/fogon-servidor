package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {
	// Crea un mock de socket.Socket usando testify/mock o una estructura personalizada
	newSocket := &MockSocket{}
	newMusico := NuevoMusico(newSocket)
	newMusico.ID = 123 // Asigna un ID al usuario para la prueba

	// Llama al método login
	newMusico.Login("USERPASS", "par_1", "VALIDA")

	// Verifica que el evento emitido sea "loginSuccess"
	assert.Equal(t, "loginSuccess", newSocket.UltimoEmitted().Event, "No dio loginSuccess")

	// Obtén el token emitido
	tokenData, ok := newSocket.UltimoEmitted().Args[0].(map[string]string)
	assert.True(t, ok, "El payload del evento no es un mapa de strings")
	tokenString, exists := tokenData["token"]
	assert.True(t, exists, "El token no fue emitido")

	// Usa VerifyToken para validar el token
	userID, err := VerifyToken(tokenString)
	assert.NoError(t, err, "Error al verificar el token")
	assert.Equal(t, newMusico.ID, userID, "El ID del usuario no coincide con el token")
}

func TestLoginFailed(t *testing.T) {
	// Crea un mock de socket.Socket usando testify/mock o una estructura personalizada
	newSocket := &MockSocket{}
	newMusico := NuevoMusico(newSocket)
	newMusico.ID = 123 // Asigna un ID al usuario para la prueba

	// Llama al método login
	newMusico.Login("USERPASS", "par_1", "OTRACONTRASEÑA")

	// Verifica que el evento emitido sea "loginSuccess"
	assert.Equal(t, "loginFailed", newSocket.UltimoEmitted().Event, "No dio loginFailed")

}
