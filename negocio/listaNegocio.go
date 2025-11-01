package negocio

import (
	"github.com/LuisWaldman/fogon-servidor/modelo"
	"github.com/LuisWaldman/fogon-servidor/servicios"
)

type ListaNegocio struct {
	cancionServicio    *servicios.CancionServicio
	listaServicio      *servicios.ListaServicio
	itemindiceServicio *servicios.ItemIndiceCancionServicio
}

func NuevoListaNegocio(cancionServicio *servicios.CancionServicio, listaServicio *servicios.ListaServicio, itemindiceServicio *servicios.ItemIndiceCancionServicio) *ListaNegocio {
	return &ListaNegocio{
		cancionServicio:    cancionServicio,
		listaServicio:      listaServicio,
		itemindiceServicio: itemindiceServicio,
	}
}

func (n *ListaNegocio) NuevaLista(nombre string, owner string) error {
	return n.listaServicio.CrearLista(nombre, owner)
}

func (n *ListaNegocio) BorrarPorID(id string) error {
	err := n.itemindiceServicio.BorrarPorListaID(id)
	if err != nil {
		return err
	}
	return n.listaServicio.BorrarPorID(id)
}

func (n *ListaNegocio) NuevaListaForzarCreacion(nombre string, owner string) error {
	lista, err := n.listaServicio.BuscarPorNombreYOwner(nombre, owner)
	if err != nil {
		return err
	}
	if lista != nil {
		err = n.BorrarPorID(lista.ID.Hex())

		if err != nil {
			return err
		}
	}
	return n.listaServicio.CrearLista(nombre, owner)
}
func (n *ListaNegocio) BorrarLista(nombreLista string, nombreUsuario string) error {
	lista, _ := n.listaServicio.BuscarPorNombreYOwner(nombreLista, nombreUsuario)
	if lista != nil {
		return n.BorrarPorID(lista.ID.Hex())
	}
	return nil
}

func (n *ListaNegocio) GetLista(nombreLista string, nombreUsuario string) (*modelo.Lista, error) {
	return n.listaServicio.BuscarPorNombreYOwner(nombreLista, nombreUsuario)
}

func (n *ListaNegocio) AgregarCancionALista(nombreLista string, nombreUsuario string, item *modelo.ItemIndiceCancion) error {
	lista, err := n.listaServicio.BuscarPorNombreYOwner(nombreLista, nombreUsuario)
	if err != nil {
		return err
	}
	lista.TotalCanciones++
	item.ListaID = lista.ID
	err = n.listaServicio.ActualizarLista(lista)
	if err != nil {
		return err
	}
	return n.itemindiceServicio.AgregarCancion(item)
}
