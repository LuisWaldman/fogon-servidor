package negocio

import (
	"github.com/LuisWaldman/fogon-servidor/modelo"
	"github.com/LuisWaldman/fogon-servidor/servicios"
)

type UsuarioNegocio struct {
	usuarioServicio *servicios.UsuarioServicio
	cancionServicio *servicios.CancionServicio
	listaServicio   *servicios.ListaServicio
}

func NuevoUsuarioNegocio(usuarioServicio *servicios.UsuarioServicio, cancionServicio *servicios.CancionServicio, listaServicio *servicios.ListaServicio) *UsuarioNegocio {
	return &UsuarioNegocio{
		usuarioServicio: usuarioServicio,
		cancionServicio: cancionServicio,
		listaServicio:   listaServicio,
	}
}

func (n *UsuarioNegocio) CrearUsuario(nombreUsuario string) error {
	user := modelo.Usuario{
		Usuario: nombreUsuario,
	}
	return n.usuarioServicio.CrearUsuario(user)
}

func (n *UsuarioNegocio) BuscarPorNombre(nombreUsuario string) (*modelo.Usuario, error) {
	return n.usuarioServicio.BuscarPorUsuario(nombreUsuario)
}

func (n *UsuarioNegocio) GetCancionesPorUsuario(nombreUsuario string) ([]modelo.Cancion, error) {
	return n.cancionServicio.BuscarPorOwner(nombreUsuario)
}
