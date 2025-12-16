package interfaces

import (
	"github.com/LuisWaldman/fogon-servidor/modelo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// IPerfilServicio define la interfaz para el servicio de perfiles
type IPerfilServicio interface {
	CrearPerfil(user modelo.Perfil) error
	BuscarPorUsuario(usuario string) (*modelo.Perfil, error)
	BorrarPorUsuario(usuario string) error
}

// IUsuarioServicio define la interfaz para el servicio de usuarios
type IUsuarioServicio interface {
	CrearUsuario(user modelo.Usuario) error
	BuscarPorUsuario(usuario string) (*modelo.Usuario, error)
	BorrarPorUsuario(usuario string) error
	ActualizarUsuario(user *modelo.Usuario) error
}

// ICancionServicio define la interfaz para el servicio de canciones
type ICancionServicio interface {
	CrearCancion(cancion *modelo.Cancion) error
	BuscarPorNombre(nombreArchivo string) (*modelo.Cancion, error)
	BuscarPorNombreYOwner(nombreArchivo string, owner string) (*modelo.Cancion, error)
	BuscarPorOwner(owner string) ([]modelo.Cancion, error)
	BorrarPorNombre(nombreArchivo string) error
	BorrarPorNombreYOwner(nombreArchivo string, owner string) error
}

// IListaServicio define la interfaz para el servicio de listas
type IListaServicio interface {
	CrearLista(nombre string, owner string) error
	BuscarPorNombreYOwner(nombre string, owner string) (*modelo.Lista, error)
	ActualizarLista(lista *modelo.Lista) error
	BorrarPorID(id string) error
}

// IItemIndiceCancionServicio define la interfaz para el servicio de items de índice
type IItemIndiceCancionServicio interface {
	AgregarCancion(item *modelo.ItemIndiceCancion) error
	GetCancionesPorListaID(listaID bson.ObjectID) []*modelo.ItemIndiceCancion
	BorrarPorListaID(id string) error
	BorrarPorID(id string) error
}
