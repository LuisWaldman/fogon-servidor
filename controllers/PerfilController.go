package controllers

import (
	"log"
	"net/http"

	"github.com/LuisWaldman/fogon-servidor/aplicacion"
	"github.com/LuisWaldman/fogon-servidor/modelo"
	servicios "github.com/LuisWaldman/fogon-servidor/servicios" // Adjust the import path as necessary
	"github.com/gin-gonic/gin"
)

type PerfilController struct {
	service    *servicios.PerfilServicio
	aplicacion *aplicacion.Aplicacion
}

func NuevoPerfilController(service *servicios.PerfilServicio, aplicacion *aplicacion.Aplicacion) *PerfilController {
	return &PerfilController{service: service, aplicacion: aplicacion} //{service: service}

}

func (sc *PerfilController) Get(c *gin.Context) {
	user, _ := c.Get("userID")
	log.Println("LLEGO A PERFIL GET", "userID", user)
	musico, encuentra := sc.aplicacion.BuscarMusicoPorID(user.(int))
	if !encuentra {
		log.Println("No se encontró el músico con ID:", user)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se encontró el músico"})
		return
	}
	perfil, _ := sc.service.BuscarPorUsuario(musico.Usuario)
	musico.Perfil = perfil // Associate the profile with the musician
	c.JSON(http.StatusOK, perfil)
}

func (sc *PerfilController) Post(c *gin.Context) {

	perfil := modelo.Perfil{}
	if err := c.ShouldBindJSON(&perfil); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, _ := c.Get("userID")                               // This is to ensure the middleware has run and set the userID
	musico, _ := sc.aplicacion.BuscarMusicoPorID(user.(int)) // Ensure user is of type string
	//perfil.Usuario = musico.Usuario
	//sc.service.CrearPerfil(perfil)
	musico.ActualizarPerfil(&perfil) // Associate the profile with the musician
	c.JSON(http.StatusCreated, gin.H{"message": "Perfil creado exitosamente"})

}
