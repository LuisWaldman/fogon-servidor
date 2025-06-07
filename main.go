package main

import (
	"log"
	"net/http"

	ConfigP "fogon-servidor/configP"

	"github.com/zishang520/socket.io/v2/socket"
)

func songsHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the request for songs here
	// For example, you can return a JSON response with a list of songs
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"songs": ["Song1", "Song2", "Song3"]}`))
}

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
	// Handle the new REST endpoint for songs
	//http.Handle("/socket.io/", io.ServeHandler(nil))
	http.HandleFunc("/api/songs", songsHandler)

	log.Fatalln(http.ListenAndServe(AppConfig.Port, nil))
}
