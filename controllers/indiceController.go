package controllers

import (
	"log"
	"net/http"

	"github.com/LuisWaldman/fogon-servidor/aplicacion"
	"github.com/LuisWaldman/fogon-servidor/servicios"

	"github.com/gin-gonic/gin"
)

type IndiceController struct {
	indiceServicio *servicios.IndiceServicio
	aplicacion     *aplicacion.Aplicacion
}

func NuevoIndiceController(indiceServicio *servicios.IndiceServicio, aplicacion *aplicacion.Aplicacion) *IndiceController {
	return &IndiceController{
		indiceServicio: indiceServicio,
		aplicacion:     aplicacion,
	}
}

// GetByOwner obtiene todos los índices de canciones para un owner específico
func (controller *IndiceController) GetByOwner(c *gin.Context) {
	user, _ := c.Get("userID")
	musico, _ := controller.aplicacion.BuscarMusicoPorID(user.(int))
	owner := musico.Usuario

	// Check if owner is provided in the query param, if yes, use it instead of the user's owner
	if ownerParam := c.Query("owner"); ownerParam != "" {
		owner = ownerParam
	}

	indices, err := controller.indiceServicio.BuscarPorOwner(owner)
	if err != nil {
		log.Println("Error obteniendo índices por owner:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	c.JSON(http.StatusOK, indices)
}

// GetByName obtiene un índice específico por nombre de archivo
func (controller *IndiceController) GetByName(c *gin.Context) {
	nombre := c.Query("nombre")
	if nombre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro 'nombre' requerido"})
		return
	}

	indice, err := controller.indiceServicio.BuscarPorNombre(nombre)
	if err != nil {
		log.Println("Error obteniendo índice por nombre:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	if indice == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Índice no encontrado"})
		return
	}

	c.JSON(http.StatusOK, indice)
}

// GetByNameAndOwner obtiene un índice específico por nombre de archivo y owner
func (controller *IndiceController) GetByNameAndOwner(c *gin.Context) {
	nombre := c.Query("nombre")
	owner := c.Query("owner")

	if nombre == "" || owner == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros 'nombre' y 'owner' requeridos"})
		return
	}

	indice, err := controller.indiceServicio.BuscarPorNombreYOwner(nombre, owner)
	if err != nil {
		log.Println("Error obteniendo índice por nombre y owner:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	if indice == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Índice no encontrado"})
		return
	}

	c.JSON(http.StatusOK, indice)
}

// GetAll obtiene todos los índices
func (controller *IndiceController) GetAll(c *gin.Context) {
	indices, err := controller.indiceServicio.ListarTodos()
	if err != nil {
		log.Println("Error obteniendo todos los índices:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	c.JSON(http.StatusOK, indices)
}

// Delete elimina un índice por nombre y owner
func (controller *IndiceController) Delete(c *gin.Context) {
	nombre := c.Query("nombre")
	owner := c.Query("owner")

	if nombre == "" {
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

	err := controller.indiceServicio.BorrarPorNombreYOwner(nombre, owner)
	if err != nil {
		log.Println("Error eliminando índice:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	log.Println("Índice eliminado exitosamente:", nombre, "Owner:", owner)
	c.JSON(http.StatusOK, gin.H{"message": "Índice eliminado exitosamente"})
}
