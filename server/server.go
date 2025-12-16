package server

import (
	"context"
	"log"
	"net/http"

	aplicacion "github.com/LuisWaldman/fogon-servidor/aplicacion"
	"github.com/LuisWaldman/fogon-servidor/aplicacion/logueadores"
	"github.com/LuisWaldman/fogon-servidor/config"
	"github.com/LuisWaldman/fogon-servidor/controllers"
	"github.com/LuisWaldman/fogon-servidor/datos"
	"github.com/LuisWaldman/fogon-servidor/interfaces"
	"github.com/LuisWaldman/fogon-servidor/negocio"
	"github.com/LuisWaldman/fogon-servidor/servicios"

	"github.com/gin-gonic/gin"
	"github.com/zishang520/socket.io/v2/socket"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Server contiene todo el estado de la aplicación y sus dependencias
type Server struct {
	config    config.Config
	app       *aplicacion.Aplicacion
	dbClient  *mongo.Client
	router    *gin.Engine
	socketIO  *socket.Server
	loginRepo *logueadores.LogeadorRepository

	// Servicios (usando interfaces)
	perfilServicio     interfaces.IPerfilServicio
	listaServicio      interfaces.IListaServicio
	cancionServicio    interfaces.ICancionServicio
	itemIndiceServicio interfaces.IItemIndiceCancionServicio
	usuarioServicio    interfaces.IUsuarioServicio

	// Servicios concretos (para dependency injection)
	perfilServicioConcreto     *servicios.PerfilServicio
	listaServicioConcreto      *servicios.ListaServicio
	cancionServicioConcreto    *servicios.CancionServicio
	itemIndiceServicioConcreto *servicios.ItemIndiceCancionServicio
	usuarioServicioConcreto    *servicios.UsuarioServicio

	// Negocio
	usuarioNegocio *negocio.UsuarioNegocio

	// Controladores
	controllers *Controllers
}

// Controllers agrupa todos los controladores
type Controllers struct {
	Perfil         *controllers.PerfilController
	RTC            *controllers.RTCController
	AnswerRTC      *controllers.AnswerRTCController
	UpdateRTC      *controllers.UpdateRTCController
	Sesiones       *controllers.SesionesController
	UsuariosSesion *controllers.UsuariosSesion
	CancionSesion  *controllers.CancionSesionController
	Cancion        *controllers.CancionController
	Lista          *controllers.ListaController
	ItemIndice     *controllers.ItemCancionesListasController
	ListaSesion    *controllers.ListaSesionController
	NumeroCancion  *controllers.NumeroCancionSesionController
	Reproductor    *controllers.ReproductorSesionController
}

// NewServer crea una nueva instancia del servidor con todas las dependencias inicializadas
func NewServer(ctx context.Context) (*Server, error) {
	// Cargar configuración
	cfg := config.LoadConfiguration()

	// Conectar a la base de datos
	dbClient, err := datos.ConnectDB()
	if err != nil {
		return nil, err
	}

	// Crear instancia del servidor
	server := &Server{
		config:   cfg,
		dbClient: dbClient,
		app:      aplicacion.NuevoAplicacion(),
	}

	// Inicializar dependencias
	if err := server.initServices(); err != nil {
		return nil, err
	}

	if err := server.initBusiness(); err != nil {
		return nil, err
	}

	if err := server.initControllers(); err != nil {
		return nil, err
	}

	if err := server.initRouter(); err != nil {
		return nil, err
	}

	if err := server.initSocketIO(ctx); err != nil {
		return nil, err
	}

	return server, nil
}

// initServices inicializa todos los servicios
func (s *Server) initServices() error {
	// Crear servicios concretos
	s.perfilServicioConcreto = servicios.NuevoPerfilServicio(s.dbClient)
	s.listaServicioConcreto = servicios.NuevoListaServicio(s.dbClient)
	s.cancionServicioConcreto = servicios.NuevoCancionServicio(s.dbClient)
	s.itemIndiceServicioConcreto = servicios.NuevoItemIndiceCancionServicio(s.dbClient)
	s.usuarioServicioConcreto = servicios.NuevoUsuarioServicio(s.dbClient)

	// Asignar a interfaces
	s.perfilServicio = s.perfilServicioConcreto
	s.listaServicio = s.listaServicioConcreto
	s.cancionServicio = s.cancionServicioConcreto
	s.itemIndiceServicio = s.itemIndiceServicioConcreto
	s.usuarioServicio = s.usuarioServicioConcreto

	return nil
}

// initBusiness inicializa la lógica de negocio
func (s *Server) initBusiness() error {
	s.usuarioNegocio = negocio.NuevoUsuarioNegocio(
		s.usuarioServicioConcreto,
		s.cancionServicioConcreto,
		s.listaServicioConcreto,
		s.itemIndiceServicioConcreto,
	)

	return nil
}

// initControllers inicializa todos los controladores
func (s *Server) initControllers() error {
	s.controllers = &Controllers{
		Perfil:         controllers.NuevoPerfilController(s.perfilServicioConcreto, s.app),
		RTC:            controllers.NuevoRTCController(s.app),
		AnswerRTC:      controllers.NuevoAnswerRTCController(s.app),
		UpdateRTC:      controllers.NuevoUpdateRTCController(s.app),
		Sesiones:       controllers.NuevoSesionesController(s.app),
		UsuariosSesion: controllers.NuevoUsuariosSesion(s.app),
		CancionSesion:  controllers.NuevoCancionSesionController(s.app),
		Cancion:        controllers.NuevoCancionController(s.cancionServicioConcreto, s.usuarioNegocio, s.app),
		Lista:          controllers.NuevoListaController(s.usuarioNegocio, s.app),
		ItemIndice:     controllers.NuevoItemCancionesListasController(s.usuarioNegocio, s.app),
		ListaSesion:    controllers.NuevoListaSesionController(s.app),
		NumeroCancion:  controllers.NuevoNumeroCancionSesionController(s.app),
		Reproductor:    controllers.NuevoReproductorSesionController(s.app),
	}

	// Inicializar repositorio de login
	s.loginRepo = logueadores.NewLogeadorRepository()
	s.loginRepo.Add("USERPASS", logueadores.NewUserPassLogeador(s.usuarioServicioConcreto))

	return nil
}

// initRouter inicializa el router con middleware y rutas
func (s *Server) initRouter() error {
	// Configurar modo según configuración
	if s.config.LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	s.router = gin.Default()
	s.router.Use(s.corsMiddleware())
	s.router.Use(s.authMiddleware())

	s.setupRoutes()

	return nil
}

// setupRoutes configura todas las rutas del servidor
func (s *Server) setupRoutes() {
	// Rutas de perfil
	s.router.GET("/perfil", s.controllers.Perfil.Get)
	s.router.POST("/perfil", s.controllers.Perfil.Post)

	// Rutas WebRTC
	s.router.POST("/answerrtc", s.controllers.AnswerRTC.Post)
	s.router.POST("/webrtc", s.controllers.RTC.Post)
	s.router.POST("/updatertc", s.controllers.UpdateRTC.Post)
	s.router.GET("/webrtc", s.controllers.RTC.Get)

	// Rutas de sesiones
	s.router.GET("/sesiones", s.controllers.Sesiones.Get)
	s.router.GET("/usersesion", s.controllers.UsuariosSesion.Get)
	s.router.GET("/cancionsesion", s.controllers.CancionSesion.Get)
	s.router.POST("/cancionsesion", s.controllers.CancionSesion.Post)

	// Rutas de listas
	s.router.GET("/lista", s.controllers.Lista.Get)
	s.router.POST("/lista", s.controllers.Lista.Post)
	s.router.PUT("/lista", s.controllers.Lista.Put)
	s.router.DELETE("/lista", s.controllers.Lista.Delete)

	// Rutas de canciones
	s.router.GET("/cancion", s.controllers.Cancion.Get)
	s.router.POST("/cancion", s.controllers.Cancion.Post)
	s.router.DELETE("/cancion", s.controllers.Cancion.Delete)

	// Rutas de items de canciones
	s.router.GET("/itemcancionlista", s.controllers.ItemIndice.GetCancionesLista)
	s.router.POST("/itemcancionlista", s.controllers.ItemIndice.PostCancionesLista)
	s.router.GET("/itemcancionusuario", s.controllers.ItemIndice.GetCancionesPorUsuario)

	// Rutas de sesiones de lista
	s.router.GET("/listasesion", s.controllers.ListaSesion.Get)
	s.router.POST("/listasesion", s.controllers.ListaSesion.Post)
	s.router.POST("/listasesionitem", s.controllers.ListaSesion.PostItem)

	// Rutas de reproductor
	s.router.GET("/numerocancion", s.controllers.NumeroCancion.Get)
	s.router.POST("/numerocancion", s.controllers.NumeroCancion.Post)
	s.router.PUT("/numerocancion", s.controllers.NumeroCancion.Put)
	s.router.POST("/tocar", s.controllers.Reproductor.PostTocar)
	s.router.POST("/tocarnro", s.controllers.Reproductor.PostTocarNro)
	s.router.PUT("/tocarnro", s.controllers.Reproductor.PutTocarNro)
}

// initSocketIO inicializa Socket.IO
func (s *Server) initSocketIO(ctx context.Context) error {
	s.socketIO = socket.NewServer(nil, nil)

	// Registrar el manejador de socket.io con el router de Gin
	s.router.Any("/socket.io/*any", gin.WrapH(s.socketIO.ServeHandler(nil)))

	err := s.socketIO.On("connection", func(clients ...any) {
		s.nuevaConexion(clients)
	})
	if err != nil {
		return err
	}

	return nil
}

// Start inicia el servidor
func (s *Server) Start() error {
	log.Printf("Iniciando servidor en puerto %s", s.config.Port)
	log.Printf("Nivel de log configurado: %s", s.config.LogLevel)

	return http.ListenAndServe(s.config.Port, s.router)
}

// Shutdown cierra gracefully el servidor
func (s *Server) Shutdown(ctx context.Context) error {
	if s.dbClient != nil {
		return s.dbClient.Disconnect(ctx)
	}
	return nil
}

// GetApp retorna la instancia de la aplicación (para compatibilidad temporal)
func (s *Server) GetApp() *aplicacion.Aplicacion {
	return s.app
}

// GetConfig retorna la configuración
func (s *Server) GetConfig() config.Config {
	return s.config
}

// SetPerfilServicio permite inyectar un mock del servicio de perfil para testing
func (s *Server) SetPerfilServicio(servicio interfaces.IPerfilServicio) {
	s.perfilServicio = servicio
}

// SetUsuarioServicio permite inyectar un mock del servicio de usuario para testing
func (s *Server) SetUsuarioServicio(servicio interfaces.IUsuarioServicio) {
	s.usuarioServicio = servicio
}

// SetCancionServicio permite inyectar un mock del servicio de canción para testing
func (s *Server) SetCancionServicio(servicio interfaces.ICancionServicio) {
	s.cancionServicio = servicio
}

// SetListaServicio permite inyectar un mock del servicio de lista para testing
func (s *Server) SetListaServicio(servicio interfaces.IListaServicio) {
	s.listaServicio = servicio
}

// SetItemIndiceServicio permite inyectar un mock del servicio de item índice para testing
func (s *Server) SetItemIndiceServicio(servicio interfaces.IItemIndiceCancionServicio) {
	s.itemIndiceServicio = servicio
}
