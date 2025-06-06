package main

import (
	"log"
	"slices"
	"time"
)

const (
	ticksPerSecond          = 30
	cameraWidth             = int(1500 / 2)
	initialTimeWithoutSpeed = time.Second * 3
	maxSpeedMultiplier      = 5 // max speed will be maxSpeedMultiplier times the initial speed
)

type PlayerStatus string

const (
	PlayerStatusOK       PlayerStatus = "ok"
	PlayerStatusDead     PlayerStatus = "muerto"
	PlayerStatusFinished PlayerStatus = "cruzoMeta"
)

type PlayerInfo struct {
	playerNumber       int
	posX               int
	posY               int
	hasGravityInverted bool
	isWalking          bool
	isDead             bool
	hasFinished        bool
}

func (playerInfo PlayerInfo) ToMap() map[string]any {
	status := PlayerStatusOK

	if playerInfo.isDead {
		status = PlayerStatusDead
	} else if playerInfo.hasFinished {
		status = PlayerStatusFinished
	}

	return map[string]any{
		"numeroJugador":          playerInfo.playerNumber,
		"x":                      playerInfo.posX,
		"y":                      playerInfo.posY,
		"tieneGravedadInvertida": playerInfo.hasGravityInverted,
		"estaCaminando":          playerInfo.isWalking,
		"estado":                 status,
	}
}

func startGame(room *Room) {
	log.Println("Loading map")
	gameMap := loadMap(room.MapName)

	log.Println("Creating world")
	world := NewWorld(gameMap, room.Players)

	gameLoop(world, room)
}

func gameLoop(world *World, room *Room) {
	ticker := time.NewTicker(time.Second / ticksPerSecond)
	initialTimer := time.NewTimer(initialTimeWithoutSpeed)
	quit := make(chan struct{})

	amountOfPlayers := len(room.Players)

	playersThatFinished := make([]int, 0, amountOfPlayers)
	playersThatDied := make([]int, 0, amountOfPlayers)

	raceStarted := false // raceStarted represents whether the race has started, since there is an initial period where the game has started but the characters do not yet have speed

	raceFinishPosX := int(world.RaceFinish.Position.X)
	log.Println("Stating game loop")

	for {
		select {
		case <-initialTimer.C:
			// add speed to characters
			for _, player := range room.Players {
				player.Character.SetInitialSpeed()
			}

			raceStarted = true

			log.Println("race started")
		case <-ticker.C:
			playersPositions := []PlayerInfo{}

			if raceStarted {
				room.Mutex.Lock()

				for _, player := range room.Players {
					finished, died := world.Update(player.Character)

					if finished {
						playersThatFinished = append(playersThatFinished, player.ID)
					} else if died {
						playersThatDied = append(playersThatDied, player.ID)
					}
				}

				room.Mutex.Unlock()
			}

			if len(playersThatFinished)+len(playersThatDied) == amountOfPlayers {
				log.Println("race finished")

				slices.Reverse(playersThatDied)

				raceResult := append(playersThatFinished, playersThatDied...)

				sendCarreraTerminada(room, raceResult)

				ticker.Stop()
				return
			} else {
				for _, player := range room.Players {
					posX := player.Character.Object.Position.X
					posY := player.Character.Object.Position.Y
					hasGravityInverted := player.Character.HasGravityInverted
					isWalking := player.Character.IsWalking
					isDead := player.Character.IsDead
					hasFinished := player.Character.HasFinished

					point := Point{
						X: int(posX),
						Y: int(posY),
					}.FromServerToClient(mapHeight, characterWidth, characterHeight)

					playersPositions = append(playersPositions, PlayerInfo{
						playerNumber:       player.ID,
						posX:               point.X,
						posY:               point.Y,
						hasGravityInverted: hasGravityInverted,
						isWalking:          isWalking,
						isDead:             isDead,
						hasFinished:        hasFinished,
					})
				}

				maxPlayerPosX := getMaxPlayerPosX(playersPositions)

				setNewSpeeds(raceFinishPosX, room, maxPlayerPosX)

				cameraX := calculateCameraPosition(maxPlayerPosX)

				world.UpdateCameraLimitPosition(cameraX)

				playersPositionsProtocol := make([]any, 0, len(playersPositions))

				for _, playersPosition := range playersPositions {
					playersPositionsProtocol = append(playersPositionsProtocol, playersPosition.ToMap())
				}

				for _, player := range room.Players {
					err := player.SendTick(playersPositionsProtocol, cameraX)
					if err != nil {
						log.Println("failed to send tick", "err", err)
					}
				}
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func setNewSpeeds(raceFinishPosX int, room *Room, maxPlayerPosX int) {
	speedMultiplier := float64(maxPlayerPosX) / float64(raceFinishPosX)

	speedScale := 1 + (speedMultiplier * (maxSpeedMultiplier - 1))

	for _, player := range room.Players {
		player.Character.ScaleSpeed(speedScale)
	}
}

func getMaxPlayerPosX(playersPositions []PlayerInfo) int {
	maxPosX := 0

	for _, playersPosition := range playersPositions {
		if playersPosition.posX > maxPosX {
			maxPosX = playersPosition.posX
		}
	}

	return maxPosX
}

func calculateCameraPosition(maxPlayerPosX int) int {
	cameraX := maxPlayerPosX - cameraWidth
	if cameraX < 0 {
		cameraX = 0
	}

	return cameraX
}

func sendCarreraTerminada(room *Room, raceResult []int) {
	playersResults := make([]map[string]any, 0, len(raceResult))

	for _, playerID := range raceResult {
		playerIndex := slices.IndexFunc(room.Players, func(player *Musico) bool {
			return player.ID == playerID
		})

		playersResults = append(
			playersResults,
			room.Players[playerIndex].ToInformacionSalaInfo(),
		)
	}

	for _, player := range room.Players {
		err := player.SendCarreraTerminada(playersResults)
		if err != nil {
			log.Println("failed to send carreraTerminada", "err", err)
		}
	}
}
