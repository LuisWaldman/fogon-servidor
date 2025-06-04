package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {
	// Crea un mock de socket.Socket usando testify/mock o una estructura personalizada
	// Implementa los métodos necesarios de socket.Socket en mockSocket si es necesario

	newSocket := &MockSocket{}
	newMusico := NuevoMusico(newSocket)
	newMusico.login("USERPASS", "pero", "par_2")
	assert.Equal(t, "loginSuccess", newSocket.UltimoEmitted().Event, "No dio loginSuccess")
}
