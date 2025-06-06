package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/zishang520/socket.io/v2/socket"
)

type Config struct {
	Port        string `json:"Port"`
	MONGODB_URI string `json:"MONGODB_URI"`
}

// AppConfig es la configuración global de la aplicación.
var AppConfig Config

func LoadConfiguration(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		return config, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}

func main() {
	var err error
	AppConfig, err = LoadConfiguration("config.json")
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
	_, err = ConnectDB()
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
