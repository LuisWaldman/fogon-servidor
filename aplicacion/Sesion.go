package aplicacion

import (
	"log"
	"strings"
	"sync"

	"github.com/LuisWaldman/fogon-servidor/modelo"
	// Adjust the import path as necessary
)

type Sesion struct {
	nombre     string
	musicos    map[int]*Musico
	estado     string
	inicio     float64
	compas     int
	Mutex      *sync.Mutex
	cancion    modelo.Cancion
	lista      []modelo.ItemIndiceCancion
	nroCancion int
}

func mismaRed(ipsA, ipsB []string) bool {
	for _, ipA := range ipsA {
		for _, ipB := range ipsB {
			if strings.HasPrefix(ipA, "192.168.") && strings.HasPrefix(ipB, "192.168.") {
				if strings.Split(ipA, ".")[2] == strings.Split(ipB, ".")[2] {
					return true
				}
			}
		}
	}
	return false
}

func (sesion *Sesion) NuevoSDP(musico *Musico) {

	sesion.Mutex.Lock()
	for _, musicoSess := range sesion.musicos {
		if (musico.ID != musicoSess.ID) && mismaRed(musico.IPs, musicoSess.IPs) {

			musico.Socket.Emit("sincronizarRTC", musicoSess.ID)
			break
		}
	}
	sesion.Mutex.Unlock()
}

func NuevaSesion(nombre string) *Sesion {
	return &Sesion{
		nombre:  nombre,
		musicos: make(map[int]*Musico),
		Mutex:   &sync.Mutex{},
	}
}

func (sesion *Sesion) GetCancion() modelo.Cancion {
	return sesion.cancion
}
func (sesion *Sesion) SetCancion(cancion modelo.Cancion) {
	sesion.cancion = cancion
	sesion.Mutex.Lock()
	for _, musico := range sesion.musicos {
		log.Print("Actualizando canción para el músico:", musico.ID)
		musico.Socket.Emit("cancionActualizada")
	}
	sesion.Mutex.Unlock()
}

func (sesion *Sesion) MensajeSesion(msj string) {
	sesion.Mutex.Lock()
	for _, musicos := range sesion.musicos {
		musicos.emit("mensajesesion", msj)
	}
	sesion.Mutex.Unlock()
}

func (sesion *Sesion) ActualizarUsuarios() {
	sesion.Mutex.Lock()
	println("Actualizando usuarios en la sesión:", sesion.nombre)
	for _, musicos := range sesion.musicos {
		musicos.emit("actualizarusuarios")
	}
	sesion.Mutex.Unlock()
}

func (sesion *Sesion) SincronizarReproduccion(compas int, time float64) {
	sesion.compas = compas
	sesion.inicio = time
	sesion.Mutex.Lock()
	for _, musico := range sesion.musicos {
		musico.emit("sincronizar", compas, sesion.inicio)
	}
	sesion.Mutex.Unlock()
}

func (sesion *Sesion) IniciarReproduccion(compas int, time float64) {
	sesion.compas = compas
	sesion.estado = "reproduciendo"
	sesion.inicio = time
	sesion.Mutex.Lock()
	for _, musico := range sesion.musicos {
		musico.emit("cancionIniciada", compas, sesion.inicio)
	}
	sesion.Mutex.Unlock()
}
func (sesion *Sesion) NroMusico(musico *Musico) int {
	sesion.Mutex.Lock()
	defer sesion.Mutex.Unlock()
	nro := 0
	for _, m := range sesion.musicos {
		if m.ID < musico.ID {
			nro++
		}
	}
	return nro

}

func (sesion *Sesion) CambiarEstado(musico *Musico, estado string) {
	nro := sesion.NroMusico(musico)
	sesion.estado = estado
	sesion.Mutex.Lock()

	for _, musico := range sesion.musicos {
		musico.emit("cancionCambioEstado", estado, nro)

	}
	sesion.Mutex.Unlock()

}

func (sesion *Sesion) ActualizarCompas(compas int) {
	sesion.compas = compas
	sesion.Mutex.Lock()
	for _, musico := range sesion.musicos {
		musico.emit("compasActualizado", compas)
	}
	sesion.Mutex.Unlock()
}

type UsuarioSesionView struct {
	ID        int    `bson:"id"`
	Usuario   string `bson:"usuario"`
	Perfil    *modelo.Perfil
	RolSesion string `bson:"rolSesion"`
}

func (sesion *Sesion) GetUsuariosView() []UsuarioSesionView {
	usuarios := make([]UsuarioSesionView, 0, len(sesion.musicos))
	sesion.Mutex.Lock()
	for _, musico := range sesion.musicos {
		usuarios = append(usuarios, UsuarioSesionView{
			ID:        musico.ID,
			Perfil:    musico.Perfil,
			Usuario:   musico.Usuario,
			RolSesion: musico.rolSesion,
		})
	}
	sesion.Mutex.Unlock()
	return usuarios
}

func (sesion *Sesion) AgregarMusico(musico *Musico) {
	if musico == nil {
		return
	}
	if sesion.musicos == nil {
		sesion.musicos = make(map[int]*Musico)
		musico.SetRolSesion("director")
	} else {
		if len(sesion.musicos) == 0 {
			musico.SetRolSesion("director")
		}
	}
	sesion.musicos[musico.ID] = musico
	if sesion.estado == "reproduciendo" {
		musico.Socket.Emit("cancionIniciada", sesion.compas, sesion.inicio)
	}
}

func (sesion *Sesion) SalirSesion(musico *Musico) {
	if musico == nil {
		return
	}
	sesion.Mutex.Lock()
	delete(sesion.musicos, musico.ID)
	if len(sesion.musicos) > 0 {
		for _, m := range sesion.musicos {
			if m.rolSesion == "director" {
				sesion.Mutex.Unlock()
				return // Al menos un director sigue en la sesión
			}
		}
		// Si no hay directores, el primero se convierte en director
		for _, m := range sesion.musicos {
			m.SetRolSesion("director")
			sesion.Mutex.Unlock()
			return
		}
	}
	sesion.Mutex.Unlock()
}

// Tocar agrega una canción a la lista en la posición siguiente a nroCancion y le suma 1
// Si la lista está vacía, la agrega y pone nroCancion = 0
func (sesion *Sesion) Tocar(cancion modelo.ItemIndiceCancion) {
	sesion.Mutex.Lock()

	if len(sesion.lista) == 0 {
		// Si la lista está vacía, agregar la canción y establecer nroCancion = 0
		sesion.lista = append(sesion.lista, cancion)
		sesion.nroCancion = 0
	} else {
		// Insertar la canción en la posición siguiente a nroCancion
		insertPos := sesion.nroCancion + 1
		if insertPos > len(sesion.lista) {
			insertPos = len(sesion.lista)
		}

		// Crear nueva slice con espacio para la nueva canción
		nuevaLista := make([]modelo.ItemIndiceCancion, len(sesion.lista)+1)
		copy(nuevaLista[:insertPos], sesion.lista[:insertPos])
		nuevaLista[insertPos] = cancion
		copy(nuevaLista[insertPos+1:], sesion.lista[insertPos:])

		sesion.lista = nuevaLista
		sesion.nroCancion++
	}

	// Notificar a todos los músicos que la lista cambió
	for _, musico := range sesion.musicos {
		musico.emit("listacambiada")
	}
	sesion.Mutex.Unlock()
}

// TocarNro cambia el número de canción actual
func (sesion *Sesion) TocarNro(numero int) {
	sesion.Mutex.Lock()

	if numero >= 0 && numero < len(sesion.lista) {
		sesion.nroCancion = numero
		// Notificar a todos los músicos que el número cambió
		for _, musico := range sesion.musicos {
			musico.emit("nrocambiado")
		}
	}
	sesion.Mutex.Unlock()
}

// SetLista establece la lista de canciones
func (sesion *Sesion) SetLista(lista []modelo.ItemIndiceCancion) {
	sesion.Mutex.Lock()

	sesion.lista = lista
	// Resetear nroCancion si la nueva lista está vacía o el índice actual es inválido
	if len(lista) == 0 || sesion.nroCancion >= len(lista) {
		sesion.nroCancion = 0
	}

	// Notificar a todos los músicos que la lista cambió
	for _, musico := range sesion.musicos {
		musico.emit("listacambiada")
	}
	sesion.Mutex.Unlock()
}

// AgregarItem agrega una canción al final de la lista sin cambiar el número de canción actual
func (sesion *Sesion) AgregarItem(cancion modelo.ItemIndiceCancion) {
	sesion.Mutex.Lock()

	// Agregar la canción al final de la lista
	sesion.lista = append(sesion.lista, cancion)

	// Notificar a todos los músicos que la lista cambió
	for _, musico := range sesion.musicos {
		musico.emit("listacambiada")
	}
	sesion.Mutex.Unlock()
}

// GetLista obtiene la lista de canciones
func (sesion *Sesion) GetLista() []modelo.ItemIndiceCancion {
	sesion.Mutex.Lock()
	defer sesion.Mutex.Unlock()

	// Retornar una copia de la lista para evitar modificaciones concurrentes
	lista := make([]modelo.ItemIndiceCancion, len(sesion.lista))
	copy(lista, sesion.lista)
	return lista
}

// GetNroCancion obtiene el número de canción actual
func (sesion *Sesion) GetNroCancion() int {
	sesion.Mutex.Lock()
	defer sesion.Mutex.Unlock()
	return sesion.nroCancion
}
