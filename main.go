package main

import (
	"log"
	"net/http"

	"github.com/zishang520/socket.io/v2/socket"
)

const Port = ":8080"
const MONGODB_URI = "mongodb+srv://luis:luis@cluster0.n2rothk.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

func main() {

	log.Println("Iniciando servidor en puerto", Port)
	ConnectDB() // Connect to the database using the function from DataBase.go

	io := socket.NewServer(nil, nil)
	http.Handle("/socket.io/", io.ServeHandler(nil))

	err := io.On("connection", func(clients ...any) {
		nuevaConexion(clients)
	})
	if err != nil {
		log.Fatalln("Error setting socket.io on connection", "err", err)
	}

	log.Fatalln(http.ListenAndServe(Port, nil))
}
