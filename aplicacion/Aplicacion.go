package aplicacion

import "log"

type Aplicacion struct {
	musicos  map[int]*Musico
	sesiones map[string]*Sesion
}

func NuevoAplicacion() *Aplicacion {
	return &Aplicacion{
		musicos:  make(map[int]*Musico),
		sesiones: make(map[string]*Sesion),
	}
}

type SesionView struct {
	Nombre   string  `bson:"nombre"`
	Latitud  float64 `bson:"latitud"`
	Longitud float64 `bson:"longitud"`
	Usuarios int     `bson:"usuarios"`
	Estado   string  `bson:"estado"`
}

func (app *Aplicacion) GetSesionView() []SesionView {
	sesiones := make([]SesionView, 0, len(app.sesiones))
	for _, sesion := range app.sesiones {
		log.Println("Procesando sesion:", sesion.nombre)
		sv := SesionView{
			Nombre:   sesion.nombre,
			Latitud:  sesion.latitud,
			Longitud: sesion.longitud,
			Usuarios: len(sesion.musicos),
			Estado:   sesion.estado,
		}
		sesiones = append(sesiones, sv)
	}
	return sesiones
}

func (app *Aplicacion) AgregarMusico(musico *Musico) {
	if musico == nil {
		return
	}
	musico.ID = len(app.musicos) + 1 // Assign a new ID based on the current size of the map
	app.musicos[musico.ID] = musico
}

func (app *Aplicacion) QuitarMusico(musico *Musico) {
	if musico == nil {
		return
	}
	delete(app.musicos, musico.ID)
}

func (app *Aplicacion) BuscarMusicoPorID(id int) (*Musico, bool) {
	musico, exists := app.musicos[id]
	if !exists {
		return nil, false
	}
	return musico, true
}

func (app *Aplicacion) CrearSesion(musico *Musico, sesion string, latitud float64, longitud float64) {
	// Check if the session already exists
	if _, exists := app.sesiones[sesion]; exists {
		musico.Socket.Emit("sesionFailed", "La sesion ya existe")
		return
	}

	// Create a new session
	newSesion := &Sesion{
		nombre:   sesion,
		latitud:  latitud,
		longitud: longitud,
	}
	app.sesiones[sesion] = newSesion
	app.UnirseSesion(musico, sesion)

}

func (app *Aplicacion) UnirseSesion(musico *Musico, sesion string) {
	// Check if the session already exists
	if _, exists := app.sesiones[sesion]; !exists {
		musico.Socket.Emit("sesionFailed", "La sesion no existe")
		return
	}
	musico.UnirseSesion(app.sesiones[sesion])
}
