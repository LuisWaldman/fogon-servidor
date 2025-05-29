package main

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/elliotchance/pie/v2"
)

type Room struct {
	ID             string
	Mutex          *sync.Mutex
	Players        []*Musico
	GameStarted    atomic.Bool
	MapName        string
	listaBandas    []string
	listaCanciones []string
	nro_cancion    int
	nro_compas     int
	director       string
	NextPlayerID   int
}

func NewRoom(mapName string) *Room {
	return &Room{
		ID:           "nueva",
		Mutex:        &sync.Mutex{},
		MapName:      mapName,
		NextPlayerID: 1,
	}
}

func (room *Room) AddPlayer(newPlayer *Musico) {
	room.Mutex.Lock()

	newPlayer.ID = room.NextPlayerID
	room.NextPlayerID++

	room.Players = append(room.Players, newPlayer)
	room.sendInformacionSala()

	newPlayer.Room = room

	room.Mutex.Unlock()
}

func (room *Room) RemovePlayer(oldPlayer *Musico) int {
	room.Mutex.Lock()

	// if the game didn't start yet, remove the player
	if !room.GameStarted.Load() {
		room.Players = pie.Filter(room.Players, func(player *Musico) bool {
			return player != oldPlayer
		})

		room.sendInformacionSala()
	}

	oldPlayer.Socket = nil
	oldPlayer.Room = nil

	playerAmount := len(room.Players)
	room.Mutex.Unlock()

	return playerAmount
}

func (room *Room) Replicar(iduser int, user string, data interface{}) bool {

	for _, player := range room.Players {
		if player.ID != iduser {
			err := player.SendReplica(user, data)
			if err != nil {
				log.Println("failed to send replica", "err", err)
			}
		}
	}
	return true
}

func (room *Room) ComunicarCambioLista() bool {

	for _, player := range room.Players {
		err := player.SendLista(room.listaBandas, room.listaCanciones)
		if err != nil {
			log.Println("failed to send replica", "err", err)
		}
	}
	return true
}

func (room *Room) ComunicarCambioCompas() bool {
	for _, player := range room.Players {
		err := player.SendCambioCompas(room.nro_compas)
		if err != nil {
			log.Println("failed to send cambioCompas", "err", err)
		}
	}
	return true
}

func (room *Room) IniciarCompas() bool {
	for _, player := range room.Players {
		err := player.IniciarCompas(room.nro_compas)
		if err != nil {
			log.Println("failed to send IniciarCompas", "err", err)
		}
	}
	return true
}

func (room *Room) ComunicarCambioCancion() bool {
	for _, player := range room.Players {
		err := player.SendCambioCancion(room.nro_cancion)
		if err != nil {
			log.Println("failed to send cambioCancion", "err", err)
		}
	}
	return true
}

func (room *Room) StartGame() bool {
	if room.GameStarted.Load() {
		return false
	}

	room.Mutex.Lock()
	room.GameStarted.Store(true)

	for _, player := range room.Players {
		err := player.SendInicioJuego()
		if err != nil {
			log.Println("failed to send inicioJuego", "err", err)
		}
	}

	go func() {
		startGame(room)

		delete(rooms, room.ID)
	}()

	room.Mutex.Unlock()

	return true
}

func (room *Room) sendInformacionSala() {
	playersInfo := make([]map[string]any, 0, len(room.Players))

	for _, player := range room.Players {
		playersInfo = append(playersInfo, player.ToInformacionSalaInfo())
	}

	for _, player := range room.Players {
		err := player.SendInformacionSala(room.ID, room.MapName, playersInfo)
		if err != nil {
			log.Println("failed to send informacionSala", "err", err)
		}
	}
}
