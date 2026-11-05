package logic

import (
	"errors"
	"strings"
)

type Deporte struct {
	Nombre string `json:"nombre"`
}

func ValidateDeporte(deporte Deporte) error {
	nombre := strings.TrimSpace(deporte.Nombre)
	if nombre == "" {
		return errors.New("el nombre del deporte no puede estar vacío")
	}
	if len(nombre) < 3 {
		return errors.New("el nombre del deporte debe tener al menos 3 caracteres")
	}

	return nil
}
