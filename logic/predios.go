package logic

import (
	"errors"
	"strings"
)

type Predio struct {
	Nombre    string `json:"nombre"`
	Direccion string `json:"direccion"`
}

func ValidatePredio(predio Predio) error {
	if strings.TrimSpace(predio.Nombre) == "" {
		return errors.New("El nombre del predio no puede estar vacío")
	}
	if strings.TrimSpace(predio.Direccion) == "" {
		return errors.New("La dirección del predio no puede estar vacía")
	}
	return nil
}
