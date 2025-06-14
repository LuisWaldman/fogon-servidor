package app

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your-secret-key")

type Emitter interface {
	Emit(ev string, args ...any) error
}

type Musico struct {
	ID      int
	Usuario string
	Socket  Emitter
}

func (player *Musico) Login(modo string, par_1 string, par_2 string) {

	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   fmt.Sprintf("%d", player.ID), // Using player ID as subject
		// You can add custom claims here
		// "name": player.Name,
		// "modo": modo,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		// Handle error, maybe send an error message to the client
		fmt.Println("Error generating JWT:", err)
		player.emit("loginError", "Failed to generate token")
		return
	}

	player.Usuario = par_1 // Assuming par_1 is the username or identifier
	err = player.emit("loginSuccess", map[string]string{"token": tokenString})
	if err != nil {
		fmt.Println("Error sending token:", err)
	}
}

func NuevoMusico(socket Emitter) *Musico {
	return &Musico{
		ID:     0, // Default ID, should be set after login
		Socket: socket,
	}
}

func (player *Musico) emit(ev string, args ...any) error {
	if player.Socket == nil {
		return nil
	}

	return player.Socket.Emit(ev, args...)
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
