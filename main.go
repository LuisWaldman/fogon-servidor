package main

import (
	"log"
	"net/http"

	Config "github.com/LuisWaldman/fogon-servidor/config"
	"github.com/LuisWaldman/fogon-servidor/controllers"
	"github.com/LuisWaldman/fogon-servidor/db"
	"github.com/LuisWaldman/fogon-servidor/servicios"

	"github.com/gin-gonic/gin"
	"github.com/zishang520/socket.io/v2/socket"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		//c.Header("Access-Control-Allow-Origin", "https://www.fogon.ar/")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Manejar preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func main() {
	router := gin.Default()
	router.Use(corsMiddleware())
	AppConfig := Config.LoadConfiguration("config.json")

	log.Println("Iniciando servidor en puerto", AppConfig.Port)
	io := socket.NewServer(nil, nil)

	// Registrar el manejador de socket.io con el router de Gin
	// Se elimina http.Handle("/socket.io/", io.ServeHandler(nil))
	// y se añade la siguiente línea:
	router.Any("/socket.io/*any", gin.WrapH(io.ServeHandler(nil)))

	err := io.On("connection", func(clients ...any) {
		nuevaConexion(clients)
	})
	if err != nil {
		log.Fatalln("Error setting socket.io on connection", "err", err)
	}

	client, err := db.ConnectDB()
	if err != nil {
		log.Fatalln("Error al conectar a la base de datos:", err)
		return
	}

	perfilServicio := servicios.NuevoPerfilServicio(client)
	constroladorServicio := controllers.NuevoPerfilController(perfilServicio)

	router.GET("/perfil", constroladorServicio.Get)
	router.POST("/perfil", constroladorServicio.Post)

	log.Fatalln(http.ListenAndServe(AppConfig.Port, router))
}
