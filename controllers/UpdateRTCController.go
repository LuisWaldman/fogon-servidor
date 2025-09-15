package controllers

import (
	"net/http"

	"github.com/LuisWaldman/fogon-servidor/aplicacion" // Adjust the import path as necessary
	"github.com/gin-gonic/gin"
)

type UpdateRTCController struct {
	aplicacion *aplicacion.Aplicacion
}

func NuevoUpdateRTCController(aplicacion *aplicacion.Aplicacion) *UpdateRTCController {
	return &UpdateRTCController{aplicacion: aplicacion}

}

func (sc *UpdateRTCController) Post(c *gin.Context) {

	var request struct {
		SDP string `json:"sdp"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing request body"})
		return
	}

	user, _ := c.Get("userID")                               // This is to ensure the middleware has run and set the userID
	musico, _ := sc.aplicacion.BuscarMusicoPorID(user.(int)) // Ensure user is of type string
	musico.UpdateSDP(request.SDP)
	c.JSON(http.StatusCreated, gin.H{"message": "Perfil creado exitosamente"})

}
