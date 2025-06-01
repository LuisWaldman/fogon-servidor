package main

import (
	"log"

	"github.com/zishang520/socket.io/v2/socket"
)

var rooms = map[string]*Room{}

func LoginUser(datas ...any) {
}

func nuevaConexion(clients []any) {
	newSocket := clients[0].(*socket.Socket)
	newClientID := newSocket.Id()

	log.Println("Nuevo Musico ID: ", newClientID)
	newMusico := NuevoMusico(newSocket)
	log.Println("Nuevo Musico: ", newMusico)
	err := newSocket.On("login", func(datas ...any) {
		log.Println("evento recivido: login", newMusico.Name, "with data:", datas)
		if len(datas) == 3 {
			modo := datas[0].(string)
			par_1 := datas[1].(string)
			par_2 := datas[2].(string)
			log.Println("Modo:", modo, "par_1:", par_1, "par_2:", par_2)
			newMusico.login(modo, par_1, par_2)
		}

	})
	if err != nil {
		log.Println("fallo registrando el mensaje holamundo", "err", err)
		newSocket.Disconnect(true)
		return
	}
	/*
		err := newClient.On("changeGravity", func(datas ...any) {
			log.Println("changeGravity event received")

			if newMusico.Room == nil {
				log.Println("changeGravity event received when the player is not in a room, ignoring message")

				return
			}

			if !newMusico.Room.GameStarted.Load() {
				log.Println("changeGravity event received when the game has not yet started, ignoring message")

				return
			}

			// TODO: mutex for each player, not only for the list
			newMusico.Room.Mutex.Lock()
			newMusico.Character.InvertGravity()
			newMusico.Room.Mutex.Unlock()
		})
		if err != nil {
			log.Println("failed to register on changeGravity message", "err", err)
			newClient.Disconnect(true)

			return
		}

		err = newClient.On("iniciarJuego", func(datas ...any) {
			log.Println("iniciarJuego event received")

			if newMusico.Room == nil {
				log.Println("iniciarJuego event received when the player is not in a room, ignoring message")

				return
			}
			return
		})
		if err != nil {
			log.Println("failed to register on iniciarJuego message", "err", err)
			newClient.Disconnect(true)

			return
		}

		err = newClient.On("unirme_sesion", func(datas ...any) {
			if len(datas) == 2 {
				sesion := datas[0].(string)
				usuario := datas[1].(string)
				newMusico.Name = usuario

				log.Println("unirSala event received with:", sesion, usuario)

				room, roomExists := rooms[sesion]
				if !roomExists {

					log.Println("crearSala event received with:", sesion, usuario)
					removeFromRoom(newMusico)
					newRoom := NewRoom(sesion)
					newRoom.director = usuario
					newRoom.ID = sesion
					newMusico.Name = usuario
					newRoom.AddPlayer(newMusico)
					rooms[newRoom.ID] = newRoom
					log.Println("Room ", newRoom.ID, " created, waiting for players")

					newMusico.emit("director", usuario)
					return
				}

				removeFromRoom(newMusico)
				newMusico.Name = usuario
				room.AddPlayer(newMusico)

				newMusico.SendLista(room.listaBandas, room.listaCanciones)

			} else {
				log.Println("unirSala event received without correct params")
			}
		})
		if err != nil {
			log.Println("failed to register on unirSala message", "err", err)
			newClient.Disconnect(true)

			return
		}

		err = newClient.On("replicar", func(datas ...any) {

			log.Println(newMusico.Name, "replicar event received with:", datas)
			if len(datas) > 2 {
				sesion := datas[0].(string)
				usuario := datas[1].(string)
				log.Println("replicar event received with:", sesion, usuario, datas[2])
				newMusico.Room.Replicar(newMusico.ID, usuario, datas[2])

			}
		})
		if err != nil {
			log.Println("failed to register on crearSala message", "err", err)
			newClient.Disconnect(true)

			return
		}

		err = newClient.On("get_director", func(datas ...any) {

			log.Println(newMusico.Name, "get_director:", datas)
			if newMusico.Room.director == "" {
				newMusico.Room.director = newMusico.Name
			}
			newMusico.emit("director", newMusico.Room.director)
		})
		if err != nil {
			log.Println("failed to register on get_director message", "err", err)
			newClient.Disconnect(true)

			return
		}

		err = newClient.On("set_cancion", func(datas ...any) {
			log.Println(newMusico.Name, "set_cancion:", datas)

			if len(datas) > 0 {
				valor := datas[0]
				log.Println("Updated cancion:", valor)

				switch v := valor.(type) {
				case float64:
					newMusico.Room.nro_cancion = int(v)
				case string:
					if nroCancion, err := strconv.Atoi(v); err == nil {
						newMusico.Room.nro_cancion = nroCancion
					} else {
						log.Println("Failed to convert valor to int:", err)
					}
				default:
					log.Println("Unexpected type for valor:", v)
				}

				newMusico.Room.ComunicarCambioCancion()
				log.Println("ComunicarCambioCancion called")
			}
		})
		if err != nil {
			log.Println("failed to register on set_cancion message", "err", err)
			newClient.Disconnect(true)

			return
		}
		err = newClient.On("setstart_compas", func(datas ...any) {
			log.Println(newMusico.Name, "setstart_compas:", datas)

			if len(datas) > 0 {

				valor := datas[0]
				switch v := valor.(type) {
				case float64:
					newMusico.Room.nro_compas = int(v)
				case string:
					if nroCancion, err := strconv.Atoi(v); err == nil {
						newMusico.Room.nro_compas = nroCancion
					} else {
						log.Println("Failed to convert valor to int:", err)
					}
				default:
					log.Println("Unexpected type for valor:", v)
				}

				log.Println("Updated compas:", newMusico.Room.nro_compas)
				newMusico.Room.IniciarCompas()
				log.Println("IniciarCompas called")
			}
		})
		if err != nil {
			log.Println("failed to register on start_compas message", "err", err)
			newClient.Disconnect(true)

			return
		}
		err = newClient.On("set_compas", func(datas ...any) {
			log.Println(newMusico.Name, "set_compas:", datas)

			if len(datas) > 0 {

				valor := datas[0]
				switch v := valor.(type) {
				case float64:
					newMusico.Room.nro_compas = int(v)
				case string:
					if nroCancion, err := strconv.Atoi(v); err == nil {
						newMusico.Room.nro_compas = nroCancion
					} else {
						log.Println("Failed to convert valor to int:", err)
					}
				default:
					log.Println("Unexpected type for valor:", v)
				}

				log.Println("Updated compas:", newMusico.Room.nro_compas)
				newMusico.Room.ComunicarCambioCompas()
				log.Println("ComunicarCambioCompas called")
			}
		})
		if err != nil {
			log.Println("failed to register on set_compas message", "err", err)
			newClient.Disconnect(true)

			return
		}

		err = newClient.On("set_lista", func(datas ...any) {

			log.Println(newMusico.Name, "set_lista:", datas)

			log.Println(newMusico.Name, "set_lista event received with:", datas)

			listaBandas := make([]string, len(datas[0].([]interface{})))
			for i, v := range datas[0].([]interface{}) {
				listaBandas[i] = v.(string)
			}
			newMusico.Room.listaBandas = listaBandas
			log.Println("Updated listaBandas:", listaBandas)

			listaCanciones := make([]string, len(datas[1].([]interface{})))
			for i, v := range datas[1].([]interface{}) {
				listaCanciones[i] = v.(string)
			}
			newMusico.Room.listaCanciones = listaCanciones
			log.Println("Updated listaCanciones:", listaCanciones)

			newMusico.Room.ComunicarCambioLista()
			log.Println("ComunicarCambioLista called")
		})
		if err != nil {
			log.Println("failed to register on get_director message", "err", err)
			newClient.Disconnect(true)

			return
		}

		err = newClient.On("disconnect", func(...any) {
			log.Println("client disconnected", newClient.Id())

			newMusico.Socket = nil

			removeFromRoom(newMusico)
		})
		if err != nil {
			log.Println("failed to register on disconnect message", "err", err)
			newClient.Disconnect(true)

			return
		}*/
}

func removeFromRoom(player *Musico) {
	if player.Room != nil {
		roomID := player.Room.ID

		playersAmount := player.Room.RemovePlayer(player)
		if playersAmount == 0 {
			delete(rooms, roomID)
		}
	}
}
