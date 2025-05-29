package main

import (
	"github.com/zishang520/socket.io/v2/socket"
)

type Musico struct {
	ID        int
	Name      string
	Socket    *socket.Socket
	Room      *Room
	Character *Character
}

func NuevoMusico(socket *socket.Socket) *Musico {
	return &Musico{
		Socket: socket,
	}
}

func (player *Musico) SendTick(playersPositions []any, cameraX int) error {
	return player.emit("tick", playersPositions, cameraX)
}

func (player *Musico) SendInicioJuego() error {
	return player.emit("inicioJuego")
}

func (player *Musico) SendReplica(nombre_usuario string, datos interface{}) error {
	return player.emit("replica", nombre_usuario, datos)
}

func (player *Musico) SendLista(bandas []string, temas []string) error {
	return player.emit("lista", bandas, temas)
}
func (player *Musico) SendCambioCompas(compas int) error {
	return player.emit("compas", compas)
}

func (player *Musico) IniciarCompas(compas int) error {
	return player.emit("start_compas", compas)
}

func (player *Musico) SendCambioCancion(cancion int) error {
	return player.emit("cancion", cancion)
}

func (player *Musico) SendInformacionSala(roomUUID string, mapName string, players []map[string]any) error {
	return player.emit("informacionSala", roomUUID, mapName, players)
}

func (player *Musico) ToInformacionSalaInfo() map[string]any {
	return map[string]any{
		"numeroJugador": player.ID,
		"nombre":        player.Name,
	}
}

func (player *Musico) SendCarreraTerminada(raceResult []map[string]any) error {
	return player.emit("carreraTerminada", raceResult)
}

func (player *Musico) emit(ev string, args ...any) error {
	if player.Socket == nil {
		return nil
	}

	return player.Socket.Emit(ev, args...)
}
