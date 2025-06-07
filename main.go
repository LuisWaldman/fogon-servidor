package main

import (
	"log"
	"net/http"

	ConfigP "fogon-servidor/configP"

	"github.com/zishang520/socket.io/v2/socket"
)

func main() {
	AppConfig, err := ConfigP.LoadConfiguration("config.json")
	if err != nil {
		log.Fatalln("Failed to load configuration:", err)
	}

	log.Println("Iniciando servidor en puerto", AppConfig.Port)
	io := socket.NewServer(nil, nil)
	http.Handle("/socket.io/", io.ServeHandler(nil))

	err = io.On("connection", func(clients ...any) {
		nuevaConexion(clients)
	})
	if err != nil {
		log.Fatalln("Error setting socket.io on connection", "err", err)
	}

	log.Fatalln(http.ListenAndServe(AppConfig.Port, nil))
}
