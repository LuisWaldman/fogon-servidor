package controllers

import (
	"log"
	"net/http"

	"github.com/LuisWaldman/fogon-servidor/app"
	"github.com/LuisWaldman/fogon-servidor/modelo"
	servicios "github.com/LuisWaldman/fogon-servidor/servicios" // Adjust the import path as necessary
	"github.com/gin-gonic/gin"
)

type PerfilController struct {
	service    *servicios.PerfilServicio
	aplicacion *app.Aplicacion
}

func NuevoPerfilController(service *servicios.PerfilServicio, aplicacion *app.Aplicacion) *PerfilController {
	return &PerfilController{service: service, aplicacion: aplicacion} //{service: service}

}

func (sc *PerfilController) Get(c *gin.Context) {
	user, _ := c.Get("userID") // This is to ensure the middleware has run and set the userID
	log.Println("LLEGO A PERFIL GET", "method", c.Request.Method, "path", c.Request.URL.Path, "userID", user)
	musico, _ := sc.aplicacion.BuscarMusicoPorID(user.(int)) // Ensure user is of type string
	perfil, _ := sc.service.BuscarPorUsuario(musico.Usuario)
	c.JSON(http.StatusOK, perfil)
}

func (sc *PerfilController) Post(c *gin.Context) {

	perfil := modelo.Perfil{}
	if err := c.ShouldBindJSON(&perfil); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, _ := c.Get("userID") // This is to ensure the middleware has run and set the userID
	log.Println("LLEGO A PERFIL POST", "method", c.Request.Method, "path", c.Request.URL.Path, "userID", user)
	musico, _ := sc.aplicacion.BuscarMusicoPorID(user.(int)) // Ensure user is of type string
	perfil.Usuario = musico.Usuario                          // Set the user for the profile
	sc.service.CrearPerfil(perfil)
	c.JSON(http.StatusCreated, gin.H{"message": "Perfil creado exitosamente"})

}
