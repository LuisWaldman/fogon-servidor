package app

type Aplicacion struct {
	musicos map[int]*Musico
}

func NuevoAplicacion() *Aplicacion {
	return &Aplicacion{
		musicos: make(map[int]*Musico),
	}
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
