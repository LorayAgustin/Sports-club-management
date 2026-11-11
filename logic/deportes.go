package logic

import (
	"errors"
	"strings"
)

type Deporte struct {
	Nombre string `json:"nombre"`
}

func ValidateDeporte(deporte *Deporte) error {
	deporte.Nombre = strings.TrimSpace(deporte.Nombre)
	if deporte.Nombre == "" {
		return errors.New("el nombre del deporte no puede estar vacío")
	}
	if len(deporte.Nombre) < 3 {
		return errors.New("el nombre del deporte debe tener al menos 3 caracteres")
	}

	deporte.Nombre = strings.ToLower(deporte.Nombre)

	return nil
}
