package controllers

import (
	"log"
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
		user, _ := c.Get("userID")
		musico, _ := controller.aplicacion.BuscarMusicoPorID(user.(int))
		owner = musico.Usuario
	}

	listas := controller.usuarioNegocio.GetCancionesPorUsuario(owner)

	c.JSON(http.StatusOK, listas)

}

func (controller *ItemCancionesListasController) GetCancionesLista(c *gin.Context) {
	nombreLista := c.Query("lista")
	owner := c.Query("owner")
	if owner == "" {
		user, _ := c.Get("userID")
		musico, _ := controller.aplicacion.BuscarMusicoPorID(user.(int))
		owner = musico.Usuario
	}

	listas := controller.usuarioNegocio.GetCancionesLista(nombreLista, owner)
	c.JSON(http.StatusOK, listas)

}

func (controller *ItemCancionesListasController) PostCancionesLista(c *gin.Context) {
	nombreLista := c.Query("lista")
	owner := c.Query("owner")

	if owner == "" {
		// Si no se proporciona owner, usar el del token de autenticación
		user, _ := c.Get("userID")
		musico, encuentra := controller.aplicacion.BuscarMusicoPorID(user.(int))
		if !encuentra {
			log.Println("No se encontró el músico con ID:", user)
			c.JSON(http.StatusBadRequest, gin.H{"error": "No se encontró el músico"})
			return
		}
		owner = musico.Usuario
	}

	var item modelo.ItemIndiceCancion
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	err := controller.usuarioNegocio.AgregarCancionALista(nombreLista, owner, &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "canción agregada"})
}
