package controllers

import (
	"net/http"

	"github.com/LuisWaldman/fogon-servidor/aplicacion"
	"github.com/LuisWaldman/fogon-servidor/modelo"
	"github.com/LuisWaldman/fogon-servidor/negocio"

	"github.com/gin-gonic/gin"
)

type ItemCancionesListasController struct {
	usuarioNegocio *negocio.UsuarioNegocio
	aplicacion     *aplicacion.Aplicacion
}

func NuevoItemCancionesListasController(usuarioNegocio *negocio.UsuarioNegocio, aplicacion *aplicacion.Aplicacion) *ItemCancionesListasController {
	return &ItemCancionesListasController{
		usuarioNegocio: usuarioNegocio,
		aplicacion:     aplicacion,
	}
}

func (controller *ItemCancionesListasController) GetCancionesPorUsuario(c *gin.Context) {
	owner := c.Query("owner")
	if owner == "" {
		// Si no se proporciona owner, usar el del token de autenticación
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro 'owner' requerido"})
			return
		}
		owner = userID.(string)
	}

	listas := controller.usuarioNegocio.GetCancionesPorUsuario(owner)

	c.JSON(http.StatusOK, listas)

}

func (controller *ItemCancionesListasController) GetCancionesLista(c *gin.Context) {
	nombreLista := c.Query("lista")
	owner := c.Query("owner")
	if owner == "" {
		// Si no se proporciona owner, usar el del token de autenticación
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro 'owner' requerido"})
			return
		}
		owner = userID.(string)
	}

	listas := controller.usuarioNegocio.GetCancionesLista(nombreLista, owner)
	c.JSON(http.StatusOK, listas)

}

func (controller *ItemCancionesListasController) PostCancionesLista(c *gin.Context) {
	nombreLista := c.Query("lista")
	owner := c.Query("owner")

	if owner == "" {
		// Si no se proporciona owner, usar el del token de autenticación
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro 'owner' requerido"})
			return
		}
		owner = userID.(string)
	}
	cancion := "cancionEjemplo"
	banda := "bandaEjemplo"
	item := modelo.NewItemIndiceCancion(cancion, banda)

	controller.usuarioNegocio.AgregarCancionALista(nombreLista, owner, item)
	c.JSON(http.StatusOK, gin.H{"status": "canción agregada"})
}
