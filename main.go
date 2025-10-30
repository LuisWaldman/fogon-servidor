package main

import (
	"log"
	"net/http"
	"strings"

	aplicacion "github.com/LuisWaldman/fogon-servidor/aplicacion"
	"github.com/LuisWaldman/fogon-servidor/aplicacion/logueadores"
	Config "github.com/LuisWaldman/fogon-servidor/config"
	"github.com/LuisWaldman/fogon-servidor/controllers"
	"github.com/LuisWaldman/fogon-servidor/datos"
	"github.com/LuisWaldman/fogon-servidor/servicios"

	"github.com/gin-gonic/gin"
	"github.com/zishang520/socket.io/v2/socket"
)

var AppConfig = Config.LoadConfiguration()

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Origin", AppConfig.Site)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
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

		if c.Request.Method == "GET" && c.Request.URL.Path == "/ntp" {
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
		userID, err := aplicacion.VerifyToken(token)
		if err != nil {
			log.Println("Error al verificar el token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}
		c.Set("userID", userID) // Almacenar el ID de usuario para su uso posterior
		c.Set("token", token)   // Puedes almacenar el token para su uso posterior

		c.Next()
	}
}

var MyApp = aplicacion.NuevoAplicacion()

func main() {
	router := gin.Default()
	router.Use(corsMiddleware())
	router.Use(AuthMiddleware())
	gin.SetMode(gin.ReleaseMode)

	client, err := datos.ConnectDB()
	if err != nil {
		log.Fatalln("Error al conectar a la base de datos:", err)
		return
	}
	log.Printf("Nivel de log configurado: %s", AppConfig.LogLevel)

	perfilServicio := servicios.NuevoPerfilServicio(client)
	listaServicio := servicios.NuevoListaServicio(client)
	//listaCancionServicio := servicios.NuevoListaCancionServicio(client)
	cancionServicio := servicios.NuevoCancionServicio(client)
	indiceServicio := servicios.NuevoIndiceServicio(client)
	usuarioServicio := servicios.NuevoUsuarioServicio(client)
	constroladorPerfil := controllers.NuevoPerfilController(perfilServicio, MyApp)
	constroladorRTC := controllers.NuevoRTCController(MyApp)
	constroladorAnswerRTC := controllers.NuevoAnswerRTCController(MyApp)
	constroladorUpdateRTC := controllers.NuevoUpdateRTCController(MyApp)
	constroladorSesiones := controllers.NuevoSesionesController(MyApp)
	constroladorUsuarioSesiones := controllers.NuevoUsuariosSesion(MyApp)
	constroladorCancionSesion := controllers.NuevoCancionSesionController(MyApp)
	constroladorCancion := controllers.NuevoCancionController(cancionServicio, MyApp)
	constroladorIndice := controllers.NuevoIndiceController(indiceServicio, MyApp)
	controladorLista := controllers.NuevoListaController(listaServicio, MyApp)
	//controladorListaCancion := controllers.NuevoListaCancionController(listaCancionServicio, listaServicio, indiceServicio, MyApp)

	loginRepo := logueadores.NewLogeadorRepository()
	loginRepo.Add("USERPASS", logueadores.NewUserPassLogeador(usuarioServicio))

	log.Println("Iniciando servidor en puerto", AppConfig.Port)
	io := socket.NewServer(nil, nil)

	// Registrar el manejador de socket.io con el router de Gin
	// Se elimina http.Handle("/socket.io/", io.ServeHandler(nil))
	// y se añade la siguiente línea:
	router.Any("/socket.io/*any", gin.WrapH(io.ServeHandler(nil)))

	err = io.On("connection", func(clients ...any) {
		nuevaConexion(clients, *loginRepo)
	})
	if err != nil {
		log.Fatalln("Error setting socket.io on connection", "err", err)
	}

	router.GET("/perfil", constroladorPerfil.Get)
	router.POST("/perfil", constroladorPerfil.Post)
	router.POST("/answerrtc", constroladorAnswerRTC.Post)
	router.POST("/webrtc", constroladorRTC.Post)
	router.POST("/updatertc", constroladorUpdateRTC.Post)
	router.GET("/webrtc", constroladorRTC.Get)

	router.GET("/sesiones", constroladorSesiones.Get)
	router.GET("/usersesion", constroladorUsuarioSesiones.Get)
	router.GET("/cancionsesion", constroladorCancionSesion.Get)
	router.POST("/cancionsesion", constroladorCancionSesion.Post)

	// Rutas para índices
	router.GET("/indice", constroladorIndice.GetByName)
	router.DELETE("/indice", constroladorIndice.Delete)
	router.GET("/indice/owner", constroladorIndice.GetByOwner)
	router.GET("/indice/search", constroladorIndice.GetByNameAndOwner)
	router.GET("/indices", constroladorIndice.GetAll)

	// Rutas para listas
	router.GET("/lista", controladorLista.GetByOwner)
	router.POST("/lista", controladorLista.Post)
	router.PUT("/lista", controladorLista.Put)
	router.DELETE("/lista", controladorLista.Delete)

	router.GET("/cancion", constroladorCancion.Get)

	router.GET("/cancion/owner", constroladorCancion.GetByOwner)
	router.POST("/cancion", constroladorCancion.Post)
	router.DELETE("/cancion", constroladorCancion.Delete)

	log.Fatalln(http.ListenAndServe(AppConfig.Port, router))
}
