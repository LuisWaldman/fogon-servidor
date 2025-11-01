package controllers

import (
	"log"
	"net/http"

	"github.com/LuisWaldman/fogon-servidor/aplicacion"
	modelo "github.com/LuisWaldman/fogon-servidor/modelo"
	"github.com/LuisWaldman/fogon-servidor/negocio"

	"github.com/gin-gonic/gin"
)

type ListaController struct {
	usuarioNegocio *negocio.UsuarioNegocio
	aplicacion     *aplicacion.Aplicacion
}

func NuevoListaController(listaNegocio *negocio.UsuarioNegocio, aplicacion *aplicacion.Aplicacion) *ListaController {
	return &ListaController{
		usuarioNegocio: listaNegocio,
		aplicacion:     aplicacion,
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
	err := controller.usuarioNegocio.AgregarLista(lista.Nombre, lista.Owner)
	if err != nil {
		log.Println("Error verificando lista existente:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	log.Println("Lista creada exitosamente:", lista.Nombre, "Owner:", lista.Owner)
	c.JSON(http.StatusOK, gin.H{"message": "Lista creada exitosamente"})
}

// GetByOwner obtiene todas las listas de un usuario
func (controller *ListaController) Get(c *gin.Context) {
	user, _ := c.Get("userID")
	musico, _ := controller.aplicacion.BuscarMusicoPorID(user.(int))
	owner := musico.Usuario

	// Check if owner is provided in the query param, if yes, use it instead of the user's owner
	if ownerParam := c.Query("owner"); ownerParam != "" {
		owner = ownerParam
	}

	listas, err := controller.usuarioNegocio.GetListasPorUsuario(owner)
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
	err := controller.usuarioNegocio.RenombrarLista(request.Nombre, request.NuevoNombre, owner)
	if err != nil {
		log.Println("Error renombrando lista:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lista actualizada exitosamente"})
}

// Delete elimina una lista
func (controller *ListaController) Delete(c *gin.Context) {
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

	// Verificar que la lista existe y pertenece al usuario
	err := controller.usuarioNegocio.BorrarLista(lista.Nombre, lista.Owner)
	if err != nil {
		log.Println("Error obteniendo lista:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lista eliminada exitosamente"})
}
