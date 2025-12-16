package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/LuisWaldman/fogon-servidor/aplicacion"
	"github.com/gin-gonic/gin"
)

// corsMiddleware maneja CORS para el servidor
func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", s.config.Site)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Manejar preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

// authMiddleware maneja la autenticación JWT para el servidor
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.RequestURI(), "/socket.io/") {
			c.Next()
			return
		}

		if c.Request.Method == "GET" && c.Request.URL.Path == "/ntp" {
			c.Next()
			return
		}

		// Obtener el token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			c.Abort()
			return
		}

		// Extraer el token eliminando "Bearer "
		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := aplicacion.VerifyToken(token, s.config.JWTSecret)
		if err != nil {
			log.Println("Error al verificar el token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		c.Set("userID", userID) // Almacenar el ID de usuario para su uso posterior
		c.Set("token", token)   // Puedes almacenar el token para su uso posterior

		c.Next()
	}
}
