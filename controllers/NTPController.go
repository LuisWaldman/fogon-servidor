package controllers

import (
	"github.com/LuisWaldman/fogon-servidor/servicios"
	"github.com/gin-gonic/gin"
)

type NTPController struct {
	service *servicios.NTPServicio
}

func NuevoNTPController(service *servicios.NTPServicio) *NTPController {
	return &NTPController{
		service: service,
	}
}

func (sc *NTPController) Get(c *gin.Context) {
	hora, err := sc.service.Get()
	if err != nil {
		c.JSON(500, gin.H{"error": "Error al obtener la hora"})
		return
	}
	c.JSON(200, gin.H{"hora": hora.Format("2006-01-02T15:04:05.000Z")})
}
