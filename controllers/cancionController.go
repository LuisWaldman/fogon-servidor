package controllers

import (
	"log"
	"net/http"

	modelo "github.com/LuisWaldman/fogon-servidor/modelo"
	"github.com/LuisWaldman/fogon-servidor/servicios"

	"github.com/gin-gonic/gin"
)

type CancionController struct {
	cancionServicio *servicios.CancionServicio
}

func NuevoCancionController(cancionServicio *servicios.CancionServicio) *CancionController {
	return &CancionController{
		cancionServicio: cancionServicio,
	}
}

func (controller *CancionController) Get(c *gin.Context) {
	nombreArchivo := c.Query("nombre")
	if nombreArchivo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro 'nombre' requerido"})
		return
	}

	cancion, err := controller.cancionServicio.BuscarPorNombre(nombreArchivo)
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

	err := controller.cancionServicio.CrearCancion(cancion)
	if err != nil {
		log.Println("Error guardando canción:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	log.Println("Canción guardada exitosamente:", cancion.NombreArchivo)
	c.JSON(http.StatusOK, gin.H{"message": "Canción guardada exitosamente"})
}
