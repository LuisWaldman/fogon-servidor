package aplicacion

import (
	"log"
	"strings"
	"sync"

	"github.com/LuisWaldman/fogon-servidor/modelo"
	// Adjust the import path as necessary
)

type Sesion struct {
	nombre  string
	musicos map[int]*Musico
	estado  string
	inicio  float64
	compas  int
	Mutex   *sync.Mutex
	cancion modelo.Cancion
}

func mismaRed(ipsA, ipsB []string) bool {
	for _, ipA := range ipsA {
		for _, ipB := range ipsB {
			if strings.HasPrefix(ipA, "192.168.") && strings.HasPrefix(ipB, "192.168.") {
				if strings.Split(ipA, ".")[2] == strings.Split(ipB, ".")[2] {
					return true
				}
			}
		}
	}
	return false
}

func (sesion *Sesion) NuevoSDP(musico *Musico) {

	sesion.Mutex.Lock()
	for _, musicoSess := range sesion.musicos {
		if (musico.ID != musicoSess.ID) && mismaRed(musico.IPs, musicoSess.IPs) {

			musico.Socket.Emit("sincronizarRTC", musicoSess.ID)
			break
		}
	}
	sesion.Mutex.Unlock()
}

func NuevaSesion(nombre string) *Sesion {
	return &Sesion{
		nombre:  nombre,
		musicos: make(map[int]*Musico),
		Mutex:   &sync.Mutex{},
	}
}

func (sesion *Sesion) GetCancion() modelo.Cancion {
	return sesion.cancion
}
func (sesion *Sesion) SetCancion(cancion modelo.Cancion) {
	sesion.cancion = cancion
	sesion.Mutex.Lock()
	for _, musico := range sesion.musicos {
		log.Print("Actualizando canción para el músico:", musico.ID)
		musico.Socket.Emit("cancionActualizada")
	}
	sesion.Mutex.Unlock()
}

func (sesion *Sesion) MensajeSesion(msj string) {
	sesion.Mutex.Lock()
	for _, musicos := range sesion.musicos {
		musicos.emit("mensajesesion", msj)
	}
	sesion.Mutex.Unlock()
}

func (sesion *Sesion) ActualizarUsuarios() {
	sesion.Mutex.Lock()
	println("Actualizando usuarios en la sesión:", sesion.nombre)
	for _, musicos := range sesion.musicos {
		musicos.emit("actualizarusuarios")
	}
	sesion.Mutex.Unlock()
}

func (sesion *Sesion) SincronizarReproduccion(compas int, time float64) {
	sesion.compas = compas
	sesion.inicio = time
	sesion.Mutex.Lock()
	for _, musico := range sesion.musicos {
		musico.emit("sincronizar", compas, sesion.inicio)
	}
	sesion.Mutex.Unlock()
}

func (sesion *Sesion) IniciarReproduccion(compas int, time float64) {
	sesion.compas = compas
	sesion.estado = "reproduciendo"
	sesion.inicio = time
	sesion.Mutex.Lock()
	for _, musico := range sesion.musicos {
		musico.emit("cancionIniciada", compas, sesion.inicio)
	}
	sesion.Mutex.Unlock()
}

func (sesion *Sesion) DetenerReproduccion() {
	sesion.Mutex.Lock()
	for _, musico := range sesion.musicos {
		musico.emit("cancionDetenida")
		sesion.estado = "pausada"
	}
	sesion.Mutex.Unlock()

}

func (sesion *Sesion) ActualizarCompas(compas int) {
	sesion.compas = compas
	sesion.Mutex.Lock()
	for _, musico := range sesion.musicos {
		musico.emit("compasActualizado", compas)
	}
	sesion.Mutex.Unlock()
}

type UsuarioSesionView struct {
	ID        int    `bson:"id"`
	Usuario   string `bson:"usuario"`
	Perfil    *modelo.Perfil
	RolSesion string `bson:"rolSesion"`
}

func (sesion *Sesion) GetUsuariosView() []UsuarioSesionView {
	usuarios := make([]UsuarioSesionView, 0, len(sesion.musicos))
	sesion.Mutex.Lock()
	for _, musico := range sesion.musicos {
		usuarios = append(usuarios, UsuarioSesionView{
			ID:        musico.ID,
			Perfil:    musico.Perfil,
			Usuario:   musico.Usuario,
			RolSesion: musico.rolSesion,
		})
	}
	sesion.Mutex.Unlock()
	return usuarios
}

func (sesion *Sesion) AgregarMusico(musico *Musico) {
	if musico == nil {
		return
	}
	if sesion.musicos == nil {
		sesion.musicos = make(map[int]*Musico)
		musico.SetRolSesion("director")

	}
	sesion.musicos[musico.ID] = musico
	if sesion.estado == "reproduciendo" {
		musico.Socket.Emit("cancionIniciada", sesion.compas, sesion.inicio)
	}
}

func (sesion *Sesion) SalirSesion(musico *Musico) {
	if musico == nil {
		return
	}
	sesion.Mutex.Lock()
	delete(sesion.musicos, musico.ID)
	if len(sesion.musicos) > 0 {
		for _, m := range sesion.musicos {
			if m.rolSesion == "director" {
				sesion.Mutex.Unlock()
				return // Al menos un director sigue en la sesión
			}
		}
		// Si no hay directores, el primero se convierte en director
		for _, m := range sesion.musicos {
			m.SetRolSesion("director")
			sesion.Mutex.Unlock()
			return
		}
	}
	sesion.Mutex.Unlock()
}
