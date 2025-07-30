package aplicacion

import (
	"log"
	"sync"
	"time"
	// Adjust the import path as necessary
)

type Sesion struct {
	nombre   string
	cancion  string
	latitud  float64
	longitud float64
	musicos  map[int]*Musico
	estado   string
	inicio   float64
	compas   int
	Mutex    *sync.Mutex
}

func NuevaSesion(nombre string) *Sesion {
	return &Sesion{
		nombre:  nombre,
		musicos: make(map[int]*Musico),
		Mutex:   &sync.Mutex{},
	}
}

func (sesion *Sesion) MensajeSesion(msj string) {
	sesion.Mutex.Lock()
	for _, musicos := range sesion.musicos {
		musicos.emit("mensajesesion", msj)
	}
	sesion.Mutex.Unlock()
}

func (sesion *Sesion) IniciarReproduccion(compas int, delay float64) {
	sesion.compas = compas
	sesion.estado = "reproduciendo"

	now := time.Now()
	_, min, sec := now.Clock()
	nsec := now.Nanosecond()
	elapsedMicrosSinceHourStart := (((min*60 + sec) * 1000000) + (nsec / 1000)) / 1000
	time := float64(elapsedMicrosSinceHourStart)

	sesion.inicio = time + delay
	log.Print("Hora para tomar: ", time, " - Inicio: ", sesion.inicio, " - Compas: ", compas, " - Delay: ", delay)
	sesion.Mutex.Lock()
	log.Print("Hora toma: ", time, " - Inicio: ", sesion.inicio, " - Compas: ", compas, " - Delay: ", delay)
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

func (sesion *Sesion) ActualizarCancion(nmCancion string) {
	sesion.cancion = nmCancion
	sesion.Mutex.Lock()
	for _, musico := range sesion.musicos {
		musico.Socket.Emit("cancionActualizada", sesion.cancion)
	}
	sesion.Mutex.Unlock()
}

type UsuarioSesionView struct {
	Usuario      string `bson:"usuario"`
	NombrePerfil string `bson:"nombre_perfil"`
	RolSesion    string `bson:"rol_sesion"`
}

func (sesion *Sesion) GetUsuariosView() []UsuarioSesionView {
	usuarios := make([]UsuarioSesionView, 0, len(sesion.musicos))
	sesion.Mutex.Lock()
	for _, musico := range sesion.musicos {
		usuarios = append(usuarios, UsuarioSesionView{
			Usuario:      musico.Usuario,
			NombrePerfil: musico.NombrePerfil,
			RolSesion:    musico.rolSesion,
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
	if sesion.cancion != "" {
		musico.Socket.Emit("cancionActualizada", sesion.cancion)
	}
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
