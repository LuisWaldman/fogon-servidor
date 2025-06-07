package main

import (
	"log"
	"net/http"

	DB "fogon-servidor/DB"
	ConfigP "fogon-servidor/configP"

	"github.com/zishang520/socket.io/v2/socket"
)

func main() {
	AppConfig, err := ConfigP.LoadConfiguration("config.json")
	if err != nil {
		log.Fatalln("Failed to load configuration:", err)
	}

	// Puedes establecer valores predeterminados si es necesario
	if AppConfig.Port == "" {
		AppConfig.Port = ":8080" // O algún valor predeterminado
	}
	if AppConfig.MONGODB_URI == "" {
		// Manejar el caso de URI vacía, quizás con un valor predeterminado o un error fatal
		log.Fatalln("MONGODB_URI cannot be empty in configuration")
	}

	// ConnectDB ahora usará AppConfig.MONGODB_URI internamente
	_, err = DB.ConnectDB()
	if err != nil {
		log.Fatalln("Failed to connect to MongoDB:", err)
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
