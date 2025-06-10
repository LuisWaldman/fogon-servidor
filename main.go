package main

import (
	"log"
	"net/http"

	Controllers "github.com/LuisWaldman/fogon-servidor/controllers"

	ConfigP "github.com/LuisWaldman/fogon-servidor/configP"

	"github.com/gin-gonic/gin"
	"github.com/zishang520/socket.io/v2/socket"
)

func main() {

	router := gin.Default()
	AppConfig, err := ConfigP.LoadConfiguration("config.json")
	if err != nil {
		log.Fatalln("Failed to load configuration:", err)
	}

	log.Println("Iniciando servidor en puerto", AppConfig.Port)
	io := socket.NewServer(nil, nil)

	// Registrar el manejador de socket.io con el router de Gin
	// Se elimina http.Handle("/socket.io/", io.ServeHandler(nil))
	// y se añade la siguiente línea:
	router.Any("/socket.io/*any", gin.WrapH(io.ServeHandler(nil)))

	err = io.On("connection", func(clients ...any) {
		nuevaConexion(clients)
	})
	if err != nil {
		log.Fatalln("Error setting socket.io on connection", "err", err)
	}
	// Handle the new REST endpoint for songs
	//http.Handle("/socket.io/", io.ServeHandler(nil)) // Esta línea ya no es necesaria y puede ser eliminada

	controller := Controllers.NuevoCancionesController()
	router.GET("/api/songs", controller.GetSongs)

	log.Fatalln(http.ListenAndServe(AppConfig.Port, router))
}
