package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/LuisWaldman/fogon-servidor/aplicacion"
	"github.com/gin-gonic/gin"
)

type NumeroCancionSesionController struct {
	aplicacion *aplicacion.Aplicacion
}

func NuevoNumeroCancionSesionController(aplicacion *aplicacion.Aplicacion) *NumeroCancionSesionController {
	return &NumeroCancionSesionController{aplicacion: aplicacion}
}

// Get obtiene el número de canción actual de la sesión
func (ncsc *NumeroCancionSesionController) Get(c *gin.Context) {
	user, _ := c.Get("userID")
	log.Println("LLEGO A NUMERO CANCION SESION GET", "method", c.Request.Method, "path", c.Request.URL.Path, "userID", user)

	musico, encuentra := ncsc.aplicacion.BuscarMusicoPorID(user.(int))
	if !encuentra {
		log.Println("No se encontró el músico con ID:", user)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se encontró el músico"})
		return
	}

	if !musico.TieneSesion() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No tienes una sesión activa"})
		return
	}

	nroCancion := musico.Sesion.GetNroCancion()
	log.Println("Número de canción actual:", nroCancion)
	c.JSON(http.StatusOK, gin.H{"nroCancion": nroCancion})
}

// Post establece el número de canción actual de la sesión
func (ncsc *NumeroCancionSesionController) Post(c *gin.Context) {
	user, _ := c.Get("userID")
	log.Println("LLEGO A NUMERO CANCION SESION POST", "method", c.Request.Method, "path", c.Request.URL.Path, "userID", user)

	musico, encuentra := ncsc.aplicacion.BuscarMusicoPorID(user.(int))
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
	log.Println("Número de canción actualizado a:", request.Numero)
	c.JSON(http.StatusOK, gin.H{"message": "Número de canción actualizado correctamente"})
}

// Put alternativa para actualizar usando query parameter
func (ncsc *NumeroCancionSesionController) Put(c *gin.Context) {
	user, _ := c.Get("userID")
	log.Println("LLEGO A NUMERO CANCION SESION PUT", "method", c.Request.Method, "path", c.Request.URL.Path, "userID", user)

	musico, encuentra := ncsc.aplicacion.BuscarMusicoPorID(user.(int))
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
	log.Println("Número de canción actualizado a:", numero)
	c.JSON(http.StatusOK, gin.H{"message": "Número de canción actualizado correctamente"})
}
