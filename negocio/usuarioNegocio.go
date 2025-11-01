package negocio

import (
	"github.com/LuisWaldman/fogon-servidor/modelo"
	"github.com/LuisWaldman/fogon-servidor/servicios"
)

type UsuarioNegocio struct {
	usuarioServicio *servicios.UsuarioServicio
	cancionServicio *servicios.CancionServicio
	listaServicio   *servicios.ListaServicio
	itemServicio    *servicios.ItemIndiceCancionServicio
}

func NuevoUsuarioNegocio(usuarioServicio *servicios.UsuarioServicio, cancionServicio *servicios.CancionServicio, listaServicio *servicios.ListaServicio, itemServicio *servicios.ItemIndiceCancionServicio) *UsuarioNegocio {
	return &UsuarioNegocio{
		usuarioServicio: usuarioServicio,
		cancionServicio: cancionServicio,
		listaServicio:   listaServicio,
		itemServicio:    itemServicio,
	}
}

func (n *UsuarioNegocio) CrearUsuario(nombreUsuario string) error {
	user := modelo.Usuario{
		Usuario: nombreUsuario,
	}
	return n.usuarioServicio.CrearUsuario(user)
}

func (n *UsuarioNegocio) BuscarPorUsuario(nombreUsuario string) (*modelo.Usuario, error) {
	return n.usuarioServicio.BuscarPorUsuario(nombreUsuario)
}

func (n *UsuarioNegocio) BorrarPorUsuario(nombreUsuario string) error {
	return n.usuarioServicio.BorrarPorUsuario(nombreUsuario)
}

func (n *UsuarioNegocio) GetCancionesPorUsuario(nombreUsuario string) ([]modelo.Cancion, error) {
	return n.cancionServicio.BuscarPorOwner(nombreUsuario)
}
