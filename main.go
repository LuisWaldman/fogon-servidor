package main

import (
	"log"
	"net/http"
	"strings"

	app "github.com/LuisWaldman/fogon-servidor/app"
	"github.com/LuisWaldman/fogon-servidor/app/logueadores"
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

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if strings.HasPrefix(c.Request.URL.RequestURI(), "/socket.io/") {
			c.Next()
			return
		}

		// Obtener el token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			c.Abort()
			return
		}

		// Extraer el token eliminando "Bearer "
		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := app.VerifyToken(token)
		if err != nil {
			log.Println("Error al verificar el token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}
		c.Set("userID", userID) // Almacenar el ID de usuario para su uso posterior

		// Validar el token (esto depende de tu lógica de autenticación)
		/*if !validateToken(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}
		*/
		// Si el token es válido, continuar con la solicitud
		c.Set("token", token) // Puedes almacenar el token para su uso posterior
		c.Next()
	}
}

var MyApp = app.NuevoAplicacion()

func main() {
	router := gin.Default()
	router.Use(corsMiddleware())
	router.Use(AuthMiddleware())

	AppConfig := Config.LoadConfiguration("config.json")

	client, err := db.ConnectDB()
	if err != nil {
		log.Fatalln("Error al conectar a la base de datos:", err)
		return
	}

	perfilServicio := servicios.NuevoPerfilServicio(client)
	constroladorServicio := controllers.NuevoPerfilController(perfilServicio, MyApp)

	usuarioServicio := servicios.NuevoUsuarioServicio(client)
	logRepo := logueadores.NewLogeadorRepository()
	logRepo.Add("USERPASS", logueadores.NewUserPassLogeador(usuarioServicio))

	log.Println("Iniciando servidor en puerto", AppConfig.Port)
	io := socket.NewServer(nil, nil)

	// Registrar el manejador de socket.io con el router de Gin
	// Se elimina http.Handle("/socket.io/", io.ServeHandler(nil))
	// y se añade la siguiente línea:
	router.Any("/socket.io/*any", gin.WrapH(io.ServeHandler(nil)))

	err = io.On("connection", func(clients ...any) {
		nuevaConexion(clients, *logRepo)
	})
	if err != nil {
		log.Fatalln("Error setting socket.io on connection", "err", err)
	}

	router.GET("/perfil", constroladorServicio.Get)
	router.POST("/perfil", constroladorServicio.Post)

	log.Fatalln(http.ListenAndServe(AppConfig.Port, router))
}
