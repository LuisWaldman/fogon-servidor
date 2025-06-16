package aplicacion

type Sesion struct {
	nombre   string
	latitud  float64
	longitud float64
	musicos  map[int]*Musico
	estado   string
}

func (sesion *Sesion) AgregarMusico(musico *Musico) {
	if musico == nil {
		return
	}
	if sesion.musicos == nil {
		sesion.musicos = make(map[int]*Musico)
	}
	sesion.musicos[musico.ID] = musico
}

func (app *Sesion) QuitarMusico(musico *Musico) {
	if musico == nil {
		return
	}
	delete(app.musicos, musico.ID)
}
