package controllers

import (
	"log"
	"net/http"

	"github.com/LuisWaldman/fogon-servidor/aplicacion" // Adjust the import path as necessary
	"github.com/gin-gonic/gin"
)

type UsuariosSesion struct {
	aplicacion *aplicacion.Aplicacion
}

func NuevoUsuariosSesion(aplicacion *aplicacion.Aplicacion) *UsuariosSesion {
	return &UsuariosSesion{aplicacion: aplicacion} //{service: service}
}

func (sc *UsuariosSesion) Get(c *gin.Context) {
	user, _ := c.Get("userID")
	log.Println("LLEGO A USUARIOS GET", "method", c.Request.Method, "path", c.Request.URL.Path, "userID", user)
	musico, _ := sc.aplicacion.BuscarMusicoPorID(user.(int))
	if !musico.TieneSesion() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No tienes una sesión activa"})
		return
	}

	usuarios := musico.Sesion.GetUsuariosView()
	log.Println("Usuarios en la sesión:", usuarios)
	c.JSON(http.StatusOK, usuarios)

}
