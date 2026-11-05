package logic

import (
	"errors"
	"strings"
)

type Categoria struct {
	Nombre    string `json:"nombre"`
	IDDeporte int    `json:"id_deporte"`
}

func ValidateCategoria(categoria Categoria) error {
	nombre := strings.TrimSpace(categoria.Nombre)
	if nombre == "" {
		return errors.New("El nombre de la categoría no puede estar vacío")
	}
	if len(nombre) < 3 {
		return errors.New("El nombre de la categoría debe tener al menos 3 caracteres")
	}
	if categoria.IDDeporte <= 0 {
		return errors.New("El ID del deporte asociado a la categoría debe ser valido")
	}

	return nil
}
