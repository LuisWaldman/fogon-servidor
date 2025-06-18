package aplicacion

type Sesion struct {
	nombre   string
	latitud  float64
	longitud float64
	musicos  map[int]*Musico
	estado   string
}

func (sesion *Sesion) MensajeSesion(msj string) {
	for _, sesion := range sesion.musicos {
		sesion.Socket.Emit("mensajesesion", msj)
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
