package controllers

import (
	"log"
	"net/http"

	"github.com/LuisWaldman/fogon-servidor/aplicacion"
	modelo "github.com/LuisWaldman/fogon-servidor/modelo"
	"github.com/LuisWaldman/fogon-servidor/servicios"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListaCancionController struct {
	listaCancionServicio *servicios.ListaCancionServicio
	listaServicio        *servicios.ListaServicio
	indiceServicio       *servicios.IndiceServicio
	aplicacion           *aplicacion.Aplicacion
}

func NuevoListaCancionController(
	listaCancionServicio *servicios.ListaCancionServicio,
	listaServicio *servicios.ListaServicio,
	indiceServicio *servicios.IndiceServicio,
	aplicacion *aplicacion.Aplicacion) *ListaCancionController {
	return &ListaCancionController{
		listaCancionServicio: listaCancionServicio,
		listaServicio:        listaServicio,
		indiceServicio:       indiceServicio,
		aplicacion:           aplicacion,
	}
}

// Post agrega una canción (ItemIndiceCancion) a una lista
func (controller *ListaCancionController) Post(c *gin.Context) {
	var request struct {
		ListaID           string                   `json:"listaId" binding:"required"`
		ItemIndiceCancion modelo.ItemIndiceCancion `json:"itemIndiceCancion" binding:"required"`
		Orden             int                      `json:"orden"`
		Notas             string                   `json:"notas"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Error al decodificar JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	// Verificar que la lista existe
	lista, err := controller.listaServicio.BuscarPorID(request.ListaID)
	if err != nil {
		log.Println("Error obteniendo lista:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	if lista == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lista no encontrada"})
		return
	}

	// Verificar permisos del usuario
	user, _ := c.Get("userID")
	musico, _ := controller.aplicacion.BuscarMusicoPorID(user.(int))
	if lista.Owner != musico.Usuario {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tiene permisos para modificar esta lista"})
		return
	}

	listaID, err := primitive.ObjectIDFromHex(request.ListaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de lista inválido"})
		return
	}

	listaCancion := modelo.NuevaListaCancion(listaID, request.ItemIndiceCancion, request.Orden)
	if request.Notas != "" {
		listaCancion.Notas = request.Notas
	}

	err = controller.listaCancionServicio.AgregarCancion(listaCancion)
	if err != nil {
		log.Println("Error agregando canción a lista:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	log.Println("Canción agregada a lista exitosamente")
	c.JSON(http.StatusOK, gin.H{"message": "Canción agregada a lista exitosamente", "id": listaCancion.ID})
}

// GetByLista obtiene todas las canciones de una lista específica
func (controller *ListaCancionController) GetByLista(c *gin.Context) {
	listaIDStr := c.Query("listaId")
	if listaIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro 'listaId' requerido"})
		return
	}

	listaID, err := primitive.ObjectIDFromHex(listaIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de lista inválido"})
		return
	}

	// Verificar que la lista existe
	lista, err := controller.listaServicio.BuscarPorID(listaIDStr)
	if err != nil {
		log.Println("Error obteniendo lista:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	if lista == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lista no encontrada"})
		return
	}

	canciones, err := controller.listaCancionServicio.ObtenerCancionesPorLista(listaID)
	if err != nil {
		log.Println("Error obteniendo canciones de lista:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	c.JSON(http.StatusOK, canciones)
}

// Put cambia el orden de una canción en la lista
func (controller *ListaCancionController) Put(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro 'id' requerido"})
		return
	}

	var request struct {
		Orden int    `json:"orden"`
		Notas string `json:"notas"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Error al decodificar JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	// Verificar que la canción existe
	listaCancion, err := controller.listaCancionServicio.BuscarPorID(id)
	if err != nil {
		log.Println("Error obteniendo canción de lista:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	if listaCancion == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Canción no encontrada en lista"})
		return
	}

	// Verificar que la lista pertenece al usuario
	lista, err := controller.listaServicio.BuscarPorID(listaCancion.ListaID.Hex())
	if err != nil {
		log.Println("Error obteniendo lista:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	if lista == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lista no encontrada"})
		return
	}

	user, _ := c.Get("userID")
	musico, _ := controller.aplicacion.BuscarMusicoPorID(user.(int))
	if lista.Owner != musico.Usuario {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tiene permisos para modificar esta lista"})
		return
	}

	// Cambiar orden si se especifica
	if request.Orden > 0 && request.Orden != listaCancion.Orden {
		err = controller.listaCancionServicio.CambiarOrden(id, request.Orden)
		if err != nil {
			log.Println("Error cambiando orden:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
			return
		}
	}

	// Actualizar notas si se especifican
	if request.Notas != listaCancion.Notas {
		err = controller.listaCancionServicio.ActualizarNotas(id, request.Notas)
		if err != nil {
			log.Println("Error actualizando notas:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
			return
		}
	}

	log.Println("Canción actualizada en lista exitosamente")
	c.JSON(http.StatusOK, gin.H{"message": "Canción actualizada exitosamente"})
}

// Delete elimina una canción de la lista
func (controller *ListaCancionController) Delete(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro 'id' requerido"})
		return
	}

	// Verificar que la canción existe
	listaCancion, err := controller.listaCancionServicio.BuscarPorID(id)
	if err != nil {
		log.Println("Error obteniendo canción de lista:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	if listaCancion == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Canción no encontrada en lista"})
		return
	}

	// Verificar que la lista pertenece al usuario
	lista, err := controller.listaServicio.BuscarPorID(listaCancion.ListaID.Hex())
	if err != nil {
		log.Println("Error obteniendo lista:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	if lista == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lista no encontrada"})
		return
	}

	user, _ := c.Get("userID")
	musico, _ := controller.aplicacion.BuscarMusicoPorID(user.(int))
	if lista.Owner != musico.Usuario {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tiene permisos para modificar esta lista"})
		return
	}

	err = controller.listaCancionServicio.EliminarCancion(id)
	if err != nil {
		log.Println("Error eliminando canción de lista:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	log.Println("Canción eliminada de lista exitosamente")
	c.JSON(http.StatusOK, gin.H{"message": "Canción eliminada de lista exitosamente"})
}

// PostByIndice agrega una canción a una lista basándose en un ItemIndiceCancion existente
func (controller *ListaCancionController) PostByIndice(c *gin.Context) {
	var request struct {
		ListaID       string `json:"listaId" binding:"required"`
		NombreArchivo string `json:"nombreArchivo" binding:"required"`
		Owner         string `json:"owner" binding:"required"`
		Orden         int    `json:"orden"`
		Notas         string `json:"notas"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Error al decodificar JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	// Buscar el ItemIndiceCancion
	itemIndice, err := controller.indiceServicio.BuscarPorNombreYOwner(request.NombreArchivo, request.Owner)
	if err != nil {
		log.Println("Error obteniendo índice:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	if itemIndice == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Canción no encontrada en índice"})
		return
	}

	// Verificar que la lista existe y pertenece al usuario
	lista, err := controller.listaServicio.BuscarPorID(request.ListaID)
	if err != nil {
		log.Println("Error obteniendo lista:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	if lista == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lista no encontrada"})
		return
	}

	user, _ := c.Get("userID")
	musico, _ := controller.aplicacion.BuscarMusicoPorID(user.(int))
	if lista.Owner != musico.Usuario {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tiene permisos para modificar esta lista"})
		return
	}

	listaID, err := primitive.ObjectIDFromHex(request.ListaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de lista inválido"})
		return
	}

	listaCancion := modelo.NuevaListaCancion(listaID, *itemIndice, request.Orden)
	if request.Notas != "" {
		listaCancion.Notas = request.Notas
	}

	err = controller.listaCancionServicio.AgregarCancion(listaCancion)
	if err != nil {
		log.Println("Error agregando canción a lista:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	log.Println("Canción agregada a lista desde índice exitosamente")
	c.JSON(http.StatusOK, gin.H{"message": "Canción agregada a lista exitosamente", "id": listaCancion.ID})
}
