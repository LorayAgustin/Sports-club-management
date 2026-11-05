package logic

import (
	"errors"
	"time"
)

type Inscripcion struct {
	IDSocio          int32      `json:"id_socio"`
	IDCategoria      int32      `json:"id_categoria"`
	FechaInscripcion time.Time  `json:"fecha_inscripcion"`
	FechaBaja        *time.Time `json:"fecha_baja,omitempty"`
	Estado           string     `json:"estado"`
}

func ValidateInscripcion(inscripcion Inscripcion) error {
	if inscripcion.IDSocio <= 0 {
		return errors.New("El ID del socio debe ser valido")
	}
	if inscripcion.IDCategoria <= 0 {
		return errors.New("El ID de la categoría debe ser valido")
	}
	if inscripcion.FechaInscripcion.IsZero() {
		return errors.New("La fecha de inscripción no puede estar vacía")
	}
	if inscripcion.FechaBaja != nil && !inscripcion.FechaBaja.IsZero() {
		if inscripcion.FechaBaja.Before(inscripcion.FechaInscripcion) {
			return errors.New("La fecha de baja no puede ser anterior a la fecha de inscripción")
		}
	}
	if inscripcion.Estado != "activo" && inscripcion.Estado != "inactivo" {
		return errors.New("El estado de la inscripción debe ser \"activo\" o \"inactivo\"")
	}
	return nil
}
