package logic

import (
	"errors"
	"strings"
	"time"
)

type Practica struct {
	IDCategoria int       `json:"id_categoria"`
	IDPredio    int       `json:"id_predio"`
	DiaSemana   string    `json:"dia_semana"`
	HoraInicio  time.Time `json:"hora_inicio"`
	HoraFin     time.Time `json:"hora_fin"`
}

var diasValidos = map[string]bool{
	"lunes":     true,
	"martes":    true,
	"miércoles": true,
	"miercoles": true,
	"jueves":    true,
	"viernes":   true,
	"sábado":    true,
	"sabado":    true,
	"domingo":   true,
}

func ValidatePractica(practica Practica) error {
	if practica.IDCategoria <= 0 {
		return errors.New("El ID de la categoría debe ser válido")
	}
	if practica.IDPredio <= 0 {
		return errors.New("El ID del predio debe ser válido")
	}
	dia := strings.ToLower(strings.TrimSpace(practica.DiaSemana))
	if !diasValidos[dia] {
		return errors.New("El día de la semana no es válido")
	}
	if practica.HoraInicio.IsZero() || practica.HoraFin.IsZero() {
		return errors.New("Las horas de inicio y fin no pueden estar vacías")
	}
	if !practica.HoraFin.After(practica.HoraInicio) {
		return errors.New("La hora de fin debe ser posterior a la hora de inicio")
	}
	return nil
}
