package controllers

import (
	"log"
	"net/http"

	"github.com/LuisWaldman/fogon-servidor/modelo"
	servicios "github.com/LuisWaldman/fogon-servidor/servicios" // Adjust the import path as necessary
	"github.com/gin-gonic/gin"
)

type PerfilController struct {
	service *servicios.PerfilServicio
}

func NuevoPerfilController(service *servicios.PerfilServicio) *PerfilController {
	return &PerfilController{service: service} //{service: service}
}

func (sc *PerfilController) Get(c *gin.Context) {
	log.Println("LLEGO A PERFIL GET", "method", c.Request.Method, "path", c.Request.URL.Path)
	perfil, _ := sc.service.BuscarPorUsuario("servicio1")
	c.JSON(http.StatusOK, perfil)
}

func (sc *PerfilController) Post(c *gin.Context) {

	perfil := modelo.Perfil{}
	if err := c.ShouldBindJSON(&perfil); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sc.service.CrearPerfil(perfil)
	c.JSON(http.StatusCreated, gin.H{"message": "Perfil creado exitosamente"})

}
