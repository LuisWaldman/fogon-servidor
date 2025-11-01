package negocio

import (
	"github.com/LuisWaldman/fogon-servidor/modelo"
	"github.com/LuisWaldman/fogon-servidor/servicios"
)

type UsuarioNegocio struct {
	usuarioServicio *servicios.UsuarioServicio
	cancionServicio *servicios.CancionServicio
	listaNegocio    *ListaNegocio
}

func NuevoUsuarioNegocio(usuarioServicio *servicios.UsuarioServicio, cancionServicio *servicios.CancionServicio, listaServicio *servicios.ListaServicio, itemServicio *servicios.ItemIndiceCancionServicio) *UsuarioNegocio {
	return &UsuarioNegocio{
		usuarioServicio: usuarioServicio,
		cancionServicio: cancionServicio,
		listaNegocio:    NuevoListaNegocio(cancionServicio, listaServicio, itemServicio),
	}
}

func (n *UsuarioNegocio) CrearUsuario(nombreUsuario string) error {
	user := modelo.Usuario{
		Usuario: nombreUsuario,
	}
	n.listaNegocio.NuevaListaForzarCreacion(nombreUsuario, "FOGON@FOGON")
	return n.usuarioServicio.CrearUsuario(user)
}

func (n *UsuarioNegocio) BuscarPorUsuario(nombreUsuario string) (*modelo.Usuario, error) {
	return n.usuarioServicio.BuscarPorUsuario(nombreUsuario)
}

func (n *UsuarioNegocio) GetCancionesPorUsuario(nombreUsuario string) []*modelo.ItemIndiceCancion {
	canciones, _ := n.listaNegocio.GetListaCanciones(nombreUsuario, "FOGON@FOGON")
	return canciones
}

func (n *UsuarioNegocio) GetCancionesLista(nombreLista string, owner string) []*modelo.ItemIndiceCancion {
	canciones, _ := n.listaNegocio.GetListaCanciones(nombreLista, owner)
	return canciones
}

func (n *UsuarioNegocio) BorrarPorUsuario(nombreUsuario string) error {
	return n.usuarioServicio.BorrarPorUsuario(nombreUsuario)
}

func (n *UsuarioNegocio) BorrarLista(nombreLista string, owner string) error {
	return n.listaNegocio.BorrarLista(nombreLista, owner)
}

func (n *UsuarioNegocio) AgregarCancion(nombreUsuario string, cancion *modelo.Cancion) error {
	n.listaNegocio.AgregarCancionALista(nombreUsuario, "FOGON@FOGON", modelo.BuildFromCancion(cancion))
	return n.cancionServicio.CrearCancion(cancion)
}

func (n *UsuarioNegocio) AgregarCancionALista(nombreLista string, nombreUsuario string, item *modelo.ItemIndiceCancion) error {
	return n.listaNegocio.AgregarCancionALista(nombreLista, nombreUsuario, item)
}

func (n *UsuarioNegocio) AgregarLista(nombreLista string, nombreUsuario string) error {
	user, _ := n.BuscarPorUsuario(nombreUsuario)
	if user == nil {
		return nil
	}
	user.Listas = append(user.Listas, nombreLista)
	n.usuarioServicio.ActualizarUsuario(user)
	return n.listaNegocio.NuevaLista(nombreLista, nombreUsuario)
}

func (n *UsuarioNegocio) GetListasPorUsuario(nombreUsuario string) ([]string, error) {
	user, _ := n.BuscarPorUsuario(nombreUsuario)
	if user == nil {
		return nil, nil
	}
	return user.Listas, nil
}

func (n *UsuarioNegocio) RenombrarLista(nombreActual string, nuevoNombre string, nombreUsuario string) error {
	lista, err := n.listaNegocio.GetLista(nombreActual, nombreUsuario)
	if err != nil {
		return err
	}

	lista.Nombre = nuevoNombre
	return n.listaNegocio.listaServicio.ActualizarLista(lista)
}
