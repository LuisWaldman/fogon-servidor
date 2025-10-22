package controllers

import (
	"log"
	"net/http"

	"github.com/LuisWaldman/fogon-servidor/aplicacion"
	modelo "github.com/LuisWaldman/fogon-servidor/modelo"
	"github.com/LuisWaldman/fogon-servidor/servicios"

	"github.com/gin-gonic/gin"
)

type ListaController struct {
	listaServicio *servicios.ListaServicio
	aplicacion    *aplicacion.Aplicacion
}

func NuevoListaController(listaServicio *servicios.ListaServicio, aplicacion *aplicacion.Aplicacion) *ListaController {
	return &ListaController{
		listaServicio: listaServicio,
		aplicacion:    aplicacion,
	}
}

// Post crea una nueva lista
func (controller *ListaController) Post(c *gin.Context) {
	var lista modelo.Lista
	if err := c.ShouldBindJSON(&lista); err != nil {
		log.Println("Error al decodificar JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	if lista.Nombre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'nombre' requerido"})
		return
	}

	user, _ := c.Get("userID")
	musico, _ := controller.aplicacion.BuscarMusicoPorID(user.(int))
	lista.Owner = musico.Usuario

	// Verificar si ya existe una lista con el mismo nombre para este usuario
	existeLista, err := controller.listaServicio.BuscarPorNombreYOwner(lista.Nombre, lista.Owner)
	if err != nil {
		log.Println("Error verificando lista existente:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	if existeLista != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Ya existe una lista con ese nombre"})
		return
	}

	nuevaLista := modelo.NuevaLista(lista.Nombre, lista.Owner)
	err = controller.listaServicio.CrearLista(nuevaLista)
	if err != nil {
		log.Println("Error guardando lista:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	log.Println("Lista creada exitosamente:", lista.Nombre, "Owner:", lista.Owner)
	c.JSON(http.StatusOK, gin.H{"message": "Lista creada exitosamente", "id": nuevaLista.ID})
}

// Get obtiene una lista por ID
func (controller *ListaController) Get(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro 'id' requerido"})
		return
	}

	lista, err := controller.listaServicio.BuscarPorID(id)
	if err != nil {
		log.Println("Error obteniendo lista:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	if lista == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lista no encontrada"})
		return
	}

	c.JSON(http.StatusOK, lista)
}

// GetByOwner obtiene todas las listas de un usuario
func (controller *ListaController) GetByOwner(c *gin.Context) {
	user, _ := c.Get("userID")
	musico, _ := controller.aplicacion.BuscarMusicoPorID(user.(int))
	owner := musico.Usuario

	// Check if owner is provided in the query param, if yes, use it instead of the user's owner
	if ownerParam := c.Query("owner"); ownerParam != "" {
		owner = ownerParam
	}

	listas, err := controller.listaServicio.BuscarPorOwner(owner)
	if err != nil {
		log.Println("Error obteniendo listas por owner:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	if len(listas) == 0 {
		c.JSON(http.StatusOK, []string{})
		return
	}

	nombreListas := make([]string, len(listas))
	for i, lista := range listas {
		nombreListas[i] = lista.Nombre
	}
	c.JSON(http.StatusOK, nombreListas)
}

// Put actualiza/renombra una lista
func (controller *ListaController) Put(c *gin.Context) {
	var request struct {
		Nombre      string `json:"nombre" binding:"required"`
		NuevoNombre string `json:"nuevoNombre" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Error al decodificar JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido o campos requeridos faltantes"})
		return
	}

	// Obtener el owner del usuario autenticado
	user, _ := c.Get("userID")
	musico, _ := controller.aplicacion.BuscarMusicoPorID(user.(int))
	owner := musico.Usuario

	// Buscar la lista por nombre y owner
	lista, err := controller.listaServicio.BuscarPorNombreYOwner(request.Nombre, owner)
	if err != nil {
		log.Println("Error obteniendo lista:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	if lista == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lista no encontrada"})
		return
	}

	// Verificar si ya existe una lista con el nuevo nombre
	if request.Nombre != request.NuevoNombre {
		existeLista, err := controller.listaServicio.BuscarPorNombreYOwner(request.NuevoNombre, owner)
		if err != nil {
			log.Println("Error verificando lista existente:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
			return
		}
		if existeLista != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Ya existe una lista con ese nombre"})
			return
		}
	}

	// Actualizar la lista con el nuevo nombre
	actualizacion := map[string]interface{}{
		"nombre": request.NuevoNombre,
	}

	err = controller.listaServicio.ActualizarLista(lista.ID.Hex(), actualizacion)
	if err != nil {
		log.Println("Error actualizando lista:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	log.Println("Lista actualizada exitosamente:", request.Nombre, "->", request.NuevoNombre)
	c.JSON(http.StatusOK, gin.H{"message": "Lista actualizada exitosamente"})
}

// Delete elimina una lista
func (controller *ListaController) Delete(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro 'id' requerido"})
		return
	}

	// Verificar que la lista existe y pertenece al usuario
	lista, err := controller.listaServicio.BuscarPorID(id)
	if err != nil {
		log.Println("Error obteniendo lista:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	if lista == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lista no encontrada"})
		return
	}

	user, _ := c.Get("userID")
	musico, _ := controller.aplicacion.BuscarMusicoPorID(user.(int))
	if lista.Owner != musico.Usuario {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tiene permisos para eliminar esta lista"})
		return
	}

	err = controller.listaServicio.BorrarPorID(id)
	if err != nil {
		log.Println("Error eliminando lista:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	log.Println("Lista eliminada exitosamente:", id)
	c.JSON(http.StatusOK, gin.H{"message": "Lista eliminada exitosamente"})
}
