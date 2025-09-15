package controllers

import (
	"net/http"

	"github.com/LuisWaldman/fogon-servidor/aplicacion" // Adjust the import path as necessary
	"github.com/gin-gonic/gin"
)

type AnswerRTCController struct {
	aplicacion *aplicacion.Aplicacion
}

func NuevoAnswerRTCController(aplicacion *aplicacion.Aplicacion) *AnswerRTCController {
	return &AnswerRTCController{aplicacion: aplicacion}

}

func (sc *AnswerRTCController) Post(c *gin.Context) {
	var request struct {
		SDP       string `json:"sdp"`
		UsuarioID int    `json:"usuarioid"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing request body"})
		return
	}

	musico, _ := sc.aplicacion.BuscarMusicoPorID(request.UsuarioID)
	musico.Socket.Emit("answerRTC", request.SDP)
	c.JSON(http.StatusCreated, gin.H{"message": "Answer enviada exitosamente"})

}
