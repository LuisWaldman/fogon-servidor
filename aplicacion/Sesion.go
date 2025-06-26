package aplicacion

import (
	"log"
	"time"

	servicios "github.com/LuisWaldman/fogon-servidor/servicios" // Adjust the import path as necessary
)

type Sesion struct {
	nombre   string
	cancion  string
	latitud  float64
	longitud float64
	musicos  map[int]*Musico
	estado   string
	inicio   time.Time
	compas   int
}

func (sesion *Sesion) MensajeSesion(msj string) {
	for _, musicos := range sesion.musicos {
		musicos.emit("mensajesesion", msj)
	}
}

func (sesion *Sesion) IniciarReproduccion(compas int, delay float64) {
	NTPServicio := servicios.NuevoNTPServicio()
	sesion.compas = compas
	hora, _ := NTPServicio.Get()
	sesion.estado = "reproduciendo"

	sesion.inicio = hora.Add(time.Duration(delay*1000) * time.Millisecond)
	log.Print("Hora: ", hora, " - Inicio: ", sesion.inicio, " - Compas: ", compas, " - Delay: ", delay)
	for _, musico := range sesion.musicos {
		musico.emit("cancionIniciada", compas, sesion.inicio.Format("2006-01-02 15:04:05.000"))
	}
}

func (sesion *Sesion) DetenerReproduccion() {
	for _, musico := range sesion.musicos {
		musico.emit("cancionDetenida")
		sesion.estado = "pausada"
	}

}

func (sesion *Sesion) ActualizarCompas(compas int) {
	sesion.compas = compas
	for _, musico := range sesion.musicos {
		musico.emit("compasActualizado", compas)
	}
}

func (sesion *Sesion) ActualizarCancion(nmCancion string) {
	sesion.cancion = nmCancion
	for _, musico := range sesion.musicos {
		musico.Socket.Emit("cancionActualizada", sesion.cancion)
	}
}

type UsuarioSesionView struct {
	Usuario      string `bson:"usuario"`
	NombrePerfil string `bson:"nombre_perfil"`
	RolSesion    string `bson:"rol_sesion"`
}

func (sesion *Sesion) GetUsuariosView() []UsuarioSesionView {
	usuarios := make([]UsuarioSesionView, 0, len(sesion.musicos))
	for _, musico := range sesion.musicos {
		usuarios = append(usuarios, UsuarioSesionView{
			Usuario:      musico.Usuario,
			NombrePerfil: musico.NombrePerfil,
			RolSesion:    musico.rolSesion,
		})
	}
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
		musico.Socket.Emit("cancionIniciada", sesion.compas, sesion.inicio.Format("2006-01-02 15:04:05.000"))
	}
}

func (app *Sesion) SalirSesion(musico *Musico) {
	if musico == nil {
		return
	}
	delete(app.musicos, musico.ID)
	if len(app.musicos) > 0 {
		for _, m := range app.musicos {
			if m.rolSesion == "director" {
				return // Al menos un director sigue en la sesión
			}
		}
		// Si no hay directores, el primero se convierte en director
		for _, m := range app.musicos {
			m.SetRolSesion("director")
			return
		}
	}
}
