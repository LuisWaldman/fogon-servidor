package aplicacion

import (
	"fmt"
	"strconv"
	"time"

	"github.com/LuisWaldman/fogon-servidor/aplicacion/logueadores"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your-secret-key")

type Emitter interface {
	Emit(ev string, args ...any) error
}

type Musico struct {
	ID           int
	Usuario      string
	Socket       Emitter
	logRepo      logueadores.LogeadorRepository
	NombrePerfil string
	Sesion       *Sesion
	rolSesion    string
}

func (musico *Musico) UnirseSesion(sesion *Sesion) {
	musico.Sesion = sesion
	musico.rolSesion = "default" // Default role for a musician
	musico.emit("ensesion", sesion.nombre)
	sesion.AgregarMusico(musico)
}

// SetRolSesion sets the role of the musician in the session
func (musico *Musico) SetRolSesion(rol string) {
	if musico.Sesion == nil {
		musico.emit("error", "No session joined")
		return
	}
	musico.rolSesion = rol
	musico.emit("rolSesion", rol)
}

func (musico *Musico) SalirSesion() {
	if musico.Sesion == nil {
		musico.emit("error", "No session joined")
		return
	}
	musico.Sesion.SalirSesion(musico)
	musico.Sesion = nil
	musico.rolSesion = "default"
	musico.emit("salirsesion", "You have left the session")
}

func (musico *Musico) IniciarReproduccion(compas int, delay float64) {
	if musico.Sesion == nil {
		musico.emit("error", "No session joined")
		return
	}
	musico.Sesion.IniciarReproduccion(compas, delay)
}

func (musico *Musico) DetenerReproduccion() {
	if musico.Sesion == nil {
		musico.emit("error", "No session joined")
		return
	}
	musico.Sesion.DetenerReproduccion()
}

func (musico *Musico) ActualizarCompas(compas int) {
	if musico.Sesion == nil {
		musico.emit("error", "No session joined")
		return
	}
	musico.Sesion.ActualizarCompas(compas)
}

func (musico *Musico) MensajeSesion(msj string) {
	if musico.Sesion == nil {
		musico.emit("error", "No session joined")
		return
	}
	musico.Sesion.MensajeSesion(msj)
}

func (musico *Musico) ActualizarCancion(nmCancion string) {
	if musico.Sesion == nil {
		musico.emit("error", "No session joined")
		return
	}
	musico.Sesion.ActualizarCancion(nmCancion)
}

func (musico *Musico) GenerarToken() {

	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   fmt.Sprintf("%d", musico.ID), // Using player ID as subject
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// Handle error, maybe send an error message to the client
		fmt.Println("Error generating JWT:", err)
		musico.emit("loginFailed", "Failed to generate token")
		return
	}
	musico.emit("conectado", map[string]string{"token": tokenString})
}

func (musico *Musico) Login(modo string, par_1 string, par_2 string) {
	if !musico.logRepo.Login(modo, par_1, par_2) {
		musico.emit("loginFailed", "Failed to generate token")
		return
	}
	musico.Usuario = par_1 // Assuming par_1 is the username or identifier
	err := musico.emit("loginSuccess", "")
	if err != nil {
		fmt.Println("Error sending token:", err)
	}
}

func NuevoMusico(socket Emitter, logRepo logueadores.LogeadorRepository) *Musico {
	return &Musico{
		ID:      0, // Default ID, should be set after login
		Socket:  socket,
		logRepo: logRepo,
	}
}

func (player *Musico) emit(ev string, args ...any) error {
	if player.Socket == nil {
		return nil
	}

	return player.Socket.Emit(ev, args...)
}

func (musico *Musico) TieneSesion() bool {
	return musico.Sesion != nil
}

func VerifyToken(tokenString string) (int, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return 0, fmt.Errorf("error parsing token: %w", err)
	}

	// Extract claims
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		// Convert the subject (user ID) back to an integer
		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			return 0, fmt.Errorf("invalid user ID in token: %w", err)
		}
		return userID, nil
	}

	return 0, fmt.Errorf("invalid token")
}
