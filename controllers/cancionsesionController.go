package controllers

import (
	"log"
	"net/http"

	"github.com/LuisWaldman/fogon-servidor/aplicacion"
	modelo "github.com/LuisWaldman/fogon-servidor/modelo"

	"github.com/gin-gonic/gin"
)

type CancionSesionController struct {
	aplicacion *aplicacion.Aplicacion
}

func NuevoCancionSesionController(aplicacion *aplicacion.Aplicacion) *CancionSesionController {
	return &CancionSesionController{
		aplicacion: aplicacion,
	}
}

func (controller *CancionSesionController) Get(c *gin.Context) {
	user, _ := c.Get("userID")
	log.Println("LLEGO A USUARIOS GET", "method", c.Request.Method, "path", c.Request.URL.Path, "userID", user)
	musico, encuentra := controller.aplicacion.BuscarMusicoPorID(user.(int))
	if !encuentra {
		log.Println("No se encontró el músico con ID:", user)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se encontró el músico"})
		return
	}
	if !musico.TieneSesion() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No tienes una sesión activa"})
		return
	}

	c.JSON(http.StatusOK, musico.Sesion.GetCancion())
}

func (controller *CancionSesionController) Post(c *gin.Context) {
	var cancion modelo.Cancion
	if err := c.ShouldBindJSON(&cancion); err != nil {
		log.Println("Error al decodificar JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	user, _ := c.Get("userID")
	log.Println("LLEGO A USUARIOS GET", "method", c.Request.Method, "path", c.Request.URL.Path, "userID", user)
	musico, encuentra := controller.aplicacion.BuscarMusicoPorID(user.(int))
	if !encuentra {
		log.Println("No se encontró el músico con ID:", user)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se encontró el músico"})
		return
	}
	if !musico.TieneSesion() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No tienes una sesión activa"})
		return
	}
	musico.Sesion.SetCancion(cancion)

	c.JSON(http.StatusOK, gin.H{"message": "Canción guardada exitosamente"})
}
