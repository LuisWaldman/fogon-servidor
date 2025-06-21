package servicios

import (
	"time"

	"github.com/beevik/ntp"
)

type NTPServicio struct {
	horaActual time.Time
}

func NuevoNTPServicio() *NTPServicio {
	tor := &NTPServicio{
		horaActual: time.Now(),
	}
	tor.ActualizarHora() // Inicializar la hora actual al crear el servicio
	return tor
}

func (s *NTPServicio) ActualizarHora() {
	ntpTime, _ := ntp.Time("pool.ntp.org")
	s.horaActual = ntpTime

}

func (s *NTPServicio) Get() (time.Time, error) {

	go s.ActualizarHora()
	return s.horaActual, nil
}
