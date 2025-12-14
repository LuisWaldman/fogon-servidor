package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LuisWaldman/fogon-servidor/server"
)

func main() {
	// Crear contexto para manejo de shutdown graceful
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Crear el servidor con todas las dependencias
	srv, err := server.NewServer(ctx)
	if err != nil {
		log.Fatalf("Error al crear el servidor: %v", err)
	}

	// Canal para escuchar señales del sistema
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Goroutine para manejar shutdown graceful
	go func() {
		<-sigChan
		log.Println("Señal de shutdown recibida, cerrando servidor...")

		cancel() // Cancelar el contexto

		// Dar tiempo para shutdown graceful
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Error durante el shutdown: %v", err)
		}

		os.Exit(0)
	}()

	// Iniciar el servidor
	if err := srv.Start(); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
