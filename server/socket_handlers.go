package server

import (
	"log"
	"time"

	aplicacion "github.com/LuisWaldman/fogon-servidor/aplicacion"
	"github.com/zishang520/socket.io/v2/socket"
)

// nuevaConexion maneja nuevas conexiones WebSocket
func (s *Server) nuevaConexion(clients []any) {
	newSocket := clients[0].(*socket.Socket)
	newMusico := aplicacion.NuevoMusico(newSocket, *s.loginRepo)
	s.app.AgregarMusico(newMusico)
	log.Println("Nuevo Musico: ", newMusico)

	s.setupSocketHandlers(newSocket, newMusico)
} // setupSocketHandlers configura todos los manejadores de eventos de socket
func (s *Server) setupSocketHandlers(newSocket *socket.Socket, newMusico *aplicacion.Musico) {
	newSocket.On("login", func(datas ...any) {
		if len(datas) == 3 {
			modo := datas[0].(string)
			par_1 := datas[1].(string)
			par_2 := datas[2].(string)
			log.Println("LOGIN - Modo:", modo, "par_1:", par_1, "par_2:", par_2)
			newMusico.Login(modo, par_1, par_2)
		}
	})

	newSocket.On("gettime", func(datas ...any) {
		now := time.Now()
		_, min, sec := now.Clock()
		nsec := now.Nanosecond()
		elapsedMicrosSinceHourStart := ((min*60 + sec) * 1000000) + (nsec / 1000)
		timeMs := int(elapsedMicrosSinceHourStart / 1000) // Convert to milliseconds
		newMusico.Socket.Emit("time", timeMs)
	})

	newSocket.On("crearsesion", func(datas ...any) {
		if len(datas) == 2 {
			sesion := datas[0].(string)
			roldefault := datas[1].(string)

			log.Println("CREAR SESION - Sesion:", sesion)
			s.app.CrearSesion(newMusico, sesion, roldefault)
		}
	})

	newSocket.On("salirsesion", func(datas ...any) {
		newMusico.SalirSesion()
		s.app.ActualizarSesiones()
		s.app.NotificarActualizarSesion()
	})

	newSocket.On("unirmesesion", func(datas ...any) {
		if len(datas) == 1 {
			sesion := datas[0].(string)
			log.Println("unirmesesion - Sesion:", sesion)
			s.app.UnirseSesion(newMusico, sesion)
		}
	})

	newSocket.On("mensajeasesion", func(datas ...any) {
		if len(datas) == 1 {
			msj := datas[0].(string)
			log.Println("mensajeasesion - Sesion:", msj)
			newMusico.MensajeSesion(msj)
		}
	})

	newSocket.On("iniciarReproduccion", func(datas ...any) {
		if len(datas) == 2 {
			compas := datas[0].(float64)
			momento := datas[1].(float64)
			log.Println("iniciarReproduccion - Sesion:", compas, "Momento:", momento)
			newMusico.IniciarReproduccion(int(compas), momento)
		}
	})

	newSocket.On("sincronizarReproduccion", func(datas ...any) {
		if len(datas) == 2 {
			compas := datas[0].(float64)
			delayms := datas[1].(float64)
			log.Println("sincronizarReproduccion - Sesion:", compas, "Delay:", delayms)
			newMusico.SincronizarReproduccion(int(compas), delayms)
		}
	})

	newSocket.On("cambiarEstado", func(datas ...any) {
		log.Println("cambiarEstado - Sesion:")
		if len(datas) == 1 {
			estado := datas[0].(string)
			log.Println("iniciarReproduccion - Sesion:", estado)
			newMusico.CambiarEstado(estado)
		}
	})

	newSocket.On("actualizarCompas", func(datas ...any) {
		if len(datas) == 1 {
			compas := datas[0].(float64)
			log.Println("actualizarCompas - Sesion:", compas)
			newMusico.ActualizarCompas(int(compas))
		}
	})

	newSocket.On("disconnect", func(...any) {
		log.Println("Musico desconectado:", newMusico.ID)
		newMusico.SalirSesion()
		s.app.ActualizarSesiones()
		s.app.QuitarMusico(newMusico)
	})
}
