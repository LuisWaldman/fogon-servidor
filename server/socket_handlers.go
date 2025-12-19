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
			modo, ok1 := datas[0].(string)
			par_1, ok2 := datas[1].(string)
			par_2, ok3 := datas[2].(string)
			if ok1 && ok2 && ok3 {
				log.Println("LOGIN - Modo:", modo, "par_1:", par_1, "par_2:", par_2)
				newMusico.Login(modo, par_1, par_2)
			} else {
				log.Println("Error: Los datos de login no son del tipo string esperado")
			}
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
			sesion, ok1 := datas[0].(string)
			roldefault, ok2 := datas[1].(string)
			if ok1 && ok2 {
				log.Println("CREAR SESION - Sesion:", sesion)
				s.app.CrearSesion(newMusico, sesion, roldefault)
			} else {
				log.Println("Error: Los datos de crear sesion no son del tipo string esperado")
			}
		}
	})

	newSocket.On("salirsesion", func(datas ...any) {
		newMusico.SalirSesion()
		s.app.ActualizarSesiones()
		s.app.NotificarActualizarSesion()
	})

	newSocket.On("unirmesesion", func(datas ...any) {
		if len(datas) == 1 {
			sesion, ok := datas[0].(string)
			if ok {
				log.Println("unirmesesion - Sesion:", sesion)
				s.app.UnirseSesion(newMusico, sesion)
			} else {
				log.Println("Error: El dato de sesion no es del tipo string esperado")
			}
		}
	})

	newSocket.On("mensajeasesion", func(datas ...any) {
		if len(datas) == 1 {
			msj, ok := datas[0].(string)
			if ok {
				log.Println("mensajeasesion - Sesion:", msj)
				newMusico.MensajeSesion(msj)
			} else {
				log.Println("Error: El mensaje no es del tipo string esperado")
			}
		}
	})

	newSocket.On("iniciarReproduccion", func(datas ...any) {
		if len(datas) == 2 {
			compas, ok1 := datas[0].(float64)
			momento, ok2 := datas[1].(float64)
			if ok1 && ok2 {
				log.Println("iniciarReproduccion - Sesion:", compas, "Momento:", momento)
				newMusico.IniciarReproduccion(int(compas), momento)
			} else {
				log.Println("Error: Los datos de iniciar reproduccion no son del tipo float64 esperado")
			}
		}
	})

	newSocket.On("sincronizarReproduccion", func(datas ...any) {
		if len(datas) == 2 {
			compas, ok1 := datas[0].(float64)
			delayms, ok2 := datas[1].(float64)
			if ok1 && ok2 {
				log.Println("sincronizarReproduccion - Sesion:", compas, "Delay:", delayms)
				newMusico.SincronizarReproduccion(int(compas), delayms)
			} else {
				log.Println("Error: Los datos de sincronizar reproduccion no son del tipo float64 esperado")
			}
		}
	})

	newSocket.On("cambiarEstado", func(datas ...any) {
		log.Println("cambiarEstado - Sesion:")
		if len(datas) == 1 {
			estado, ok := datas[0].(string)
			if ok {
				log.Println("cambiarEstado - Estado:", estado)
				newMusico.CambiarEstado(estado)
			} else {
				log.Println("Error: El estado no es del tipo string esperado")
			}
		}
	})

	newSocket.On("actualizarCompas", func(datas ...any) {
		if len(datas) == 1 {
			compas, ok := datas[0].(float64)
			if ok {
				log.Println("actualizarCompas - Sesion:", compas)
				newMusico.ActualizarCompas(int(compas))
			} else {
				log.Println("Error: El compas no es del tipo float64 esperado")
			}
		}
	})

	newSocket.On("setrolausuario", func(datas ...any) {
		if len(datas) == 2 {
			usuarioIDInt, ok1 := datas[0].(float64) // Los números de JavaScript llegan como float64
			rol, ok2 := datas[1].(string)
			if ok1 && ok2 {
				usuarioID := int(usuarioIDInt)
				log.Println("setrolausuario - UsuarioID:", usuarioID, "Rol:", rol)
				newMusico.SetRolAUsuario(usuarioID, rol)
			} else {
				log.Println("Error: Los datos recibidos no son del tipo esperado (int, string)")
			}
		}
	})

	newSocket.On("disconnect", func(...any) {
		log.Println("Musico desconectado:", newMusico.ID)
		newMusico.SalirSesion()
		s.app.ActualizarSesiones()
		s.app.QuitarMusico(newMusico)
	})
}
