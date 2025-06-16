package controllers

import (
	"net/http"

	"github.com/LuisWaldman/fogon-servidor/aplicacion" // Adjust the import path as necessary
	"github.com/gin-gonic/gin"
)

type SesionesController struct {
	aplicacion *aplicacion.Aplicacion
}

func NuevoSesionesController(aplicacion *aplicacion.Aplicacion) *SesionesController {
	return &SesionesController{aplicacion: aplicacion} //{service: service}
}

func (sc *SesionesController) Get(c *gin.Context) {
	c.JSON(http.StatusOK, sc.aplicacion.GetSesionView())
}
