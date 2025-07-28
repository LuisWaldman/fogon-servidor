package servicios

import (
	"sync"
	"time"
)

var (
	ntpServicioInstance *NTPServicio
	once                sync.Once
)

type NTPServicio struct {
	horaActual  time.Time
	ultimoDelta time.Time
	Delta       time.Duration
}

func NuevoNTPServicio() *NTPServicio {
	once.Do(func() {
		ntpServicioInstance = &NTPServicio{
			ultimoDelta: time.Now().Add(-time.Hour * 24), // Inicializar con un valor pasado
			horaActual:  time.Now(),
			Delta:       0,
		}
		//go ntpServicioInstance.ActualizarHora()
	})
	return ntpServicioInstance
}

/*
	func (s *NTPServicio) ActualizarHora() {
		s.horaActual = time.Now()
		if time.Since(s.ultimoDelta) < time.Minute {
			return // No actualizar si ya se actualizó en el último minuto
		}
		ntpTime, _ := ntp.Time("pool.ntp.org")
		s.ultimoDelta = ntpTime
		s.Delta = time.Until(ntpTime)
	}
*/
func (s *NTPServicio) Get() time.Time {

	//go s.ActualizarHora()
	return time.Now()
}
