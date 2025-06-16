package aplicacion

type Sesion struct {
	sesion   string
	latitud  float64
	longitud float64
	musicos  map[int]*Musico
}

func (app *Sesion) AgregarMusico(musico *Musico) {
	if musico == nil {
		return
	}
	musico.ID = len(app.musicos) + 1 // Assign a new ID based on the current size of the map
	app.musicos[musico.ID] = musico
}

func (app *Sesion) QuitarMusico(musico *Musico) {
	if musico == nil {
		return
	}
	delete(app.musicos, musico.ID)
}
