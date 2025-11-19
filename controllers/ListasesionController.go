package controllers

import (
	"log"
	"net/http"

	"github.com/LuisWaldman/fogon-servidor/aplicacion"
	modelo "github.com/LuisWaldman/fogon-servidor/modelo"

	"github.com/gin-gonic/gin"
)

type ListaSesionController struct {
	aplicacion *aplicacion.Aplicacion
}

func NuevoListaSesionController(aplicacion *aplicacion.Aplicacion) *ListaSesionController {
	return &ListaSesionController{
		aplicacion: aplicacion,
	}
}

// Get obtiene la lista de canciones de la sesión actual
func (controller *ListaSesionController) Get(c *gin.Context) {
	user, _ := c.Get("userID")
	log.Println("LLEGO A LISTA SESION GET", "method", c.Request.Method, "path", c.Request.URL.Path, "userID", user)

	musico, encuentra := controller.aplicacion.BuscarMusicoPorID(user.(int))
	if !encuentra {
		log.Println("No se encontró el músico con ID:", user)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se encontró el músico"})
		return
	}

	if !musico.TieneSesion() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No tienes una sesión activa"})
		return
	}

	lista := musico.Sesion.GetLista()
	log.Println("Lista de canciones en la sesión:", len(lista), "canciones")
	c.JSON(http.StatusOK, lista)
}

// Post establece la lista de canciones de la sesión actual
func (controller *ListaSesionController) Post(c *gin.Context) {
	user, _ := c.Get("userID")
	log.Println("LLEGO A LISTA SESION POST", "method", c.Request.Method, "path", c.Request.URL.Path, "userID", user)

	musico, encuentra := controller.aplicacion.BuscarMusicoPorID(user.(int))
	if !encuentra {
		log.Println("No se encontró el músico con ID:", user)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se encontró el músico"})
		return
	}

	if !musico.TieneSesion() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No tienes una sesión activa"})
		return
	}

	var lista []modelo.ItemIndiceCancion
	if err := c.ShouldBindJSON(&lista); err != nil {
		log.Println("Error al parsear la lista:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de lista inválido"})
		return
	}

	musico.Sesion.SetLista(lista)
	log.Println("Lista de canciones actualizada:", len(lista), "canciones")
	c.JSON(http.StatusOK, gin.H{"message": "Lista actualizada correctamente"})
}

// PostItem agrega un item individual a la lista de canciones de la sesión
func (controller *ListaSesionController) PostItem(c *gin.Context) {
	user, _ := c.Get("userID")
	log.Println("LLEGO A LISTA SESION POST ITEM", "method", c.Request.Method, "path", c.Request.URL.Path, "userID", user)

	musico, encuentra := controller.aplicacion.BuscarMusicoPorID(user.(int))
	if !encuentra {
		log.Println("No se encontró el músico con ID:", user)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se encontró el músico"})
		return
	}

	if !musico.TieneSesion() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No tienes una sesión activa"})
		return
	}

	var item modelo.ItemIndiceCancion
	if err := c.ShouldBindJSON(&item); err != nil {
		log.Println("Error al parsear el item:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de item inválido"})
		return
	}

	musico.Sesion.AgregarItem(item)
	log.Println("Item agregado a la lista:", item.Cancion, "de", item.Banda)
	c.JSON(http.StatusOK, gin.H{"message": "Item agregado correctamente a la lista"})
}
