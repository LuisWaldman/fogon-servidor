package logueadores

import (
	"github.com/LuisWaldman/fogon-servidor/servicios"
)

type UserPassLogeador struct {
	usuarioServicio *servicios.UsuarioServicio
}

// NewUserPassLogeador crea una nueva instancia de UserPassLogeador
func NewUserPassLogeador(usuarioServicio *servicios.UsuarioServicio) *UserPassLogeador {
	return &UserPassLogeador{
		usuarioServicio: usuarioServicio,
	}
}

func (l *UserPassLogeador) Login(par_1 string, par_2 string) bool {
	usuario, _ := l.usuarioServicio.BuscarPorUsuario(par_1)
	if usuario.Encontrado == false {
		return false // Usuario no encontrado
		/*
			usuaio := modelo.Usuario{
				Usuario:   par_1,
				Clave:     par_2,
				Modologin: "UserPass",
			}

			l.usuarioServicio.CrearUsuario(usuaio)
		*/
	}

	// Si el usuario ya existe, verificamos la clave
	if usuario.Clave == par_2 {
		return true // Login exitoso
	}
	return false // Clave incorrecta
}
