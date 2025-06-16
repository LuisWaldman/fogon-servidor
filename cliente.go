package main

import (
	"log"

	aplicacion "github.com/LuisWaldman/fogon-servidor/aplicacion"
	"github.com/LuisWaldman/fogon-servidor/aplicacion/logueadores"
	"github.com/zishang520/socket.io/v2/socket"
)

func LoginUser(datas ...any) {
}

func nuevaConexion(clients []any, logRepo logueadores.LogeadorRepository) {
	newSocket := clients[0].(*socket.Socket)
	newMusico := aplicacion.NuevoMusico(newSocket, logRepo)
	MyApp.AgregarMusico(newMusico)
	log.Println("Nuevo Musico: ", newMusico)
	newSocket.On("login", func(datas ...any) {
		if len(datas) == 3 {
			modo := datas[0].(string)
			par_1 := datas[1].(string)
			par_2 := datas[2].(string)
			log.Println("LOGIN - Modo:", modo, "par_1:", par_1, "par_2:", par_2)
			newMusico.Login(modo, par_1, par_2)
		}
	})
	newSocket.On("crearsesion", func(datas ...any) {
		if len(datas) == 3 {
			sesion := datas[0].(string)
			latitud := datas[1].(float64)
			longitud := datas[2].(float64)
			log.Println("CREAR SESION - Sesion:", sesion, "Latitud:", latitud, "Longitud:", longitud)
			MyApp.CrearSesion(newMusico, sesion, latitud, longitud)
		}
	})
	newSocket.On("disconnect", func(...any) {
		MyApp.QuitarMusico(newMusico)
	})
}

func removeFromRoom(player *aplicacion.Musico) {
	/*
		if player.Room != nil {
			roomID := player.Room.ID

			playersAmount := player.Room.RemovePlayer(player)
			if playersAmount == 0 {
				delete(rooms, roomID)
			}
		}*/
}
