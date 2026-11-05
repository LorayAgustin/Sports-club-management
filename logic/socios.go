package logic

import (
	"errors"
	"strings"
	"time"
)

type Socio struct {
	IDUsuario int32      `json:"id_usuario"`
	DNI       string     `json:"dni"`
	FechaAlta time.Time  `json:"fecha_alta"`
	FechaBaja *time.Time `json:"fecha_baja,omitempty"`
}

func ValidateSocio(socio Socio) error {
	if socio.IDUsuario <= 0 {
		return errors.New("El ID del usuario debe ser valido")
	}
	dni := strings.TrimSpace(socio.DNI)
	if dni == "" {
		return errors.New("El DNI del socio no puede estar vacío")
	}
	if len(dni) < 7 || len(dni) > 10 {
		return errors.New("El DNI del socio debe tener entre 7 y 10 caracteres")
	}
	if socio.FechaAlta.IsZero() {
		return errors.New("La fecha de alta del socio no puede estar vacía")
	}
	if socio.FechaBaja != nil && !socio.FechaBaja.IsZero() {
		if socio.FechaBaja.Before(socio.FechaAlta) {
			return errors.New("La fecha de baja no puede ser anterior a la fecha de alta")
		}
	}
	return nil
}
