package controllers

import (
	"log"
	"net/http"

	modelo "github.com/LuisWaldman/fogon-servidor/modelo"
	"github.com/LuisWaldman/fogon-servidor/servicios"

	"github.com/LuisWaldman/fogon-servidor/aplicacion"
	"github.com/gin-gonic/gin"
)

type CancionController struct {
	cancionServicio *servicios.CancionServicio
	aplicacion      *aplicacion.Aplicacion
}

func NuevoCancionController(cancionServicio *servicios.CancionServicio, aplicacion *aplicacion.Aplicacion) *CancionController {
	return &CancionController{
		cancionServicio: cancionServicio,
		aplicacion:      aplicacion,
	}
}

func (controller *CancionController) Get(c *gin.Context) {
	nombreArchivo := c.Query("nombre")

	if nombreArchivo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro 'nombre' requerido"})
		return
	}

	var cancion *modelo.Cancion
	var err error
	user, _ := c.Get("userID")
	musico, _ := controller.aplicacion.BuscarMusicoPorID(user.(int))
	owner := musico.Usuario

	// Check if owner is provided in the query param, if yes, use it instead of the user's owner
	if ownerParam := c.Query("owner"); ownerParam != "" {
		owner = ownerParam
	}

	if owner != "" {
		// Buscar por nombre y owner si se proporciona owner
		cancion, err = controller.cancionServicio.BuscarPorNombreYOwner(nombreArchivo, owner)
	} else {
		// Buscar solo por nombre (comportamiento anterior para compatibilidad)
		cancion, err = controller.cancionServicio.BuscarPorNombre(nombreArchivo)
	}

	if err != nil {
		log.Println("Error obteniendo canción:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	if cancion == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Canción no encontrada"})
		return
	}

	c.JSON(http.StatusOK, cancion)
}

func (controller *CancionController) Post(c *gin.Context) {
	var cancion modelo.Cancion
	if err := c.ShouldBindJSON(&cancion); err != nil {
		log.Println("Error al decodificar JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	if cancion.NombreArchivo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'nombreArchivo' requerido"})
		return
	}

	user, _ := c.Get("userID")
	musico, _ := controller.aplicacion.BuscarMusicoPorID(user.(int))
	cancion.Owner = musico.Usuario

	err := controller.cancionServicio.CrearCancion(cancion)
	if err != nil {
		log.Println("Error guardando canción:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	log.Println("Canción guardada exitosamente:", cancion.NombreArchivo, "Owner:", cancion.Owner)
	c.JSON(http.StatusOK, gin.H{"message": "Canción guardada exitosamente"})
}

// Delete elimina una canción por nombre y owner
func (controller *CancionController) Delete(c *gin.Context) {
	nombreArchivo := c.Query("nombre")
	owner := c.Query("owner")

	if nombreArchivo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro 'nombre' requerido"})
		return
	}

	// Si no se proporciona owner, usar el del token de autenticación
	if owner == "" {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Owner requerido"})
			return
		}
		owner = userID.(string)
	}

	err := controller.cancionServicio.BorrarPorNombreYOwner(nombreArchivo, owner)
	if err != nil {
		log.Println("Error eliminando canción:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	log.Println("Canción eliminada exitosamente:", nombreArchivo, "Owner:", owner)
	c.JSON(http.StatusOK, gin.H{"message": "Canción eliminada exitosamente"})
}
