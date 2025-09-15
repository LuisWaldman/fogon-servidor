package controllers

import (
	"net/http"
	"strconv"

	"github.com/LuisWaldman/fogon-servidor/aplicacion" // Adjust the import path as necessary
	"github.com/gin-gonic/gin"
)

type RTCController struct {
	aplicacion *aplicacion.Aplicacion
}

func NuevoRTCController(aplicacion *aplicacion.Aplicacion) *RTCController {
	return &RTCController{aplicacion: aplicacion}

}

func (sc *RTCController) Get(c *gin.Context) {
	usuarioid := c.Query("usuarioid")
	if usuarioid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro 'usuarioid' requerido"})
		return
	}

	var request struct {
		SDP string `json:"sdp"`
	}

	usuarioIDInt, err := strconv.Atoi(usuarioid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	musico, _ := sc.aplicacion.BuscarMusicoPorID(usuarioIDInt)
	request.SDP = musico.SDP
	c.JSON(http.StatusOK, request)
}
func (sc *RTCController) Post(c *gin.Context) {

	var request struct {
		SDP string `json:"sdp"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing request body"})
		return
	}

	user, _ := c.Get("userID")                               // This is to ensure the middleware has run and set the userID
	musico, _ := sc.aplicacion.BuscarMusicoPorID(user.(int)) // Ensure user is of type string
	musico.SetSDP(request.SDP)
	c.JSON(http.StatusCreated, gin.H{"message": "Perfil creado exitosamente"})

}
