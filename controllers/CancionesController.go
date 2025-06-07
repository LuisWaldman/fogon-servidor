package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CancionesController struct {
	//service *services.SongService
}

func NuevoCancionesController() *CancionesController {
	return &CancionesController{} //{service: service}
}

func (sc *CancionesController) GetSongs(c *gin.Context) {
	log.Println("Fetching songs", "method", c.Request.Method, "path", c.Request.URL.Path)
	c.JSON(http.StatusOK, map[string]string{"songs": "Song1"})
}
