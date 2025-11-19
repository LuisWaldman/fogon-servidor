package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/LuisWaldman/fogon-servidor/aplicacion"
	"github.com/LuisWaldman/fogon-servidor/modelo"
	"github.com/gin-gonic/gin"
)

type ReproductorSesionController struct {
	aplicacion *aplicacion.Aplicacion
}

func NuevoReproductorSesionController(aplicacion *aplicacion.Aplicacion) *ReproductorSesionController {
	return &ReproductorSesionController{aplicacion: aplicacion}
}

// PostTocar agrega una canción a la lista en la siguiente posición y la toca
func (rsc *ReproductorSesionController) PostTocar(c *gin.Context) {
	user, _ := c.Get("userID")
	log.Println("LLEGO A REPRODUCTOR SESION POST TOCAR", "method", c.Request.Method, "path", c.Request.URL.Path, "userID", user)

	musico, encuentra := rsc.aplicacion.BuscarMusicoPorID(user.(int))
	if !encuentra {
		log.Println("No se encontró el músico con ID:", user)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se encontró el músico"})
		return
	}

	if !musico.TieneSesion() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No tienes una sesión activa"})
		return
	}

	var cancion modelo.ItemIndiceCancion
	if err := c.ShouldBindJSON(&cancion); err != nil {
		log.Println("Error al parsear la canción:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de canción inválido"})
		return
	}

	musico.Sesion.Tocar(cancion)
	log.Println("Canción agregada y puesta en reproducción:", cancion.Cancion, "de", cancion.Banda)
	c.JSON(http.StatusOK, gin.H{"message": "Canción agregada y en reproducción"})
}

// PostTocarNro cambia el número de canción actual para reproducir
func (rsc *ReproductorSesionController) PostTocarNro(c *gin.Context) {
	user, _ := c.Get("userID")
	log.Println("LLEGO A REPRODUCTOR SESION POST TOCAR NRO", "method", c.Request.Method, "path", c.Request.URL.Path, "userID", user)

	musico, encuentra := rsc.aplicacion.BuscarMusicoPorID(user.(int))
	if !encuentra {
		log.Println("No se encontró el músico con ID:", user)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se encontró el músico"})
		return
	}

	if !musico.TieneSesion() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No tienes una sesión activa"})
		return
	}

	// Leer el número desde el body JSON
	var request struct {
		Numero int `json:"numero"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Error al parsear el número:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de número inválido"})
		return
	}

	musico.Sesion.TocarNro(request.Numero)
	log.Println("Cambiado a canción número:", request.Numero)
	c.JSON(http.StatusOK, gin.H{"message": "Cambiado a canción número " + strconv.Itoa(request.Numero)})
}

// PutTocarNro alternativa usando query parameter para cambiar número de canción
func (rsc *ReproductorSesionController) PutTocarNro(c *gin.Context) {
	user, _ := c.Get("userID")
	log.Println("LLEGO A REPRODUCTOR SESION PUT TOCAR NRO", "method", c.Request.Method, "path", c.Request.URL.Path, "userID", user)

	musico, encuentra := rsc.aplicacion.BuscarMusicoPorID(user.(int))
	if !encuentra {
		log.Println("No se encontró el músico con ID:", user)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se encontró el músico"})
		return
	}

	if !musico.TieneSesion() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No tienes una sesión activa"})
		return
	}

	// Leer el número desde query parameter
	numeroStr := c.Query("numero")
	if numeroStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro 'numero' requerido"})
		return
	}

	numero, err := strconv.Atoi(numeroStr)
	if err != nil {
		log.Println("Error al convertir número:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Número inválido"})
		return
	}

	musico.Sesion.TocarNro(numero)
	log.Println("Cambiado a canción número:", numero)
	c.JSON(http.StatusOK, gin.H{"message": "Cambiado a canción número " + strconv.Itoa(numero)})
}
