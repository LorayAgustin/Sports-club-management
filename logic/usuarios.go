package logic

import (
	"errors"
	"strings"
	"time"
)

type Usuario struct {
	Nombre          string    `json:"nombre"`
	Email           string    `json:"email"`
	Contrasena      string    `json:"contrasena"`
	FechaNacimiento time.Time `json:"fecha_nacimiento"`
	Rol             string    `json:"rol"`
}

func ValidateUsuario(usuario Usuario) error {
	if usuario.Nombre == "" {
		return errors.New("el nombre no puede estar vacío")
	}
	if usuario.Email == "" {
		return errors.New("el email no puede estar vacío")
	}
	if !strings.Contains(usuario.Email, "@") {
		return errors.New("el email debe contener un '@'")
	}
	if len(usuario.Contrasena) < 8 {
		return errors.New("la contraseña debe tener al menos 8 caracteres")
	}
	if usuario.FechaNacimiento.After(time.Now()) {
		return errors.New("la fecha de nacimiento no puede ser posterior a la fecha actual")
	}
	//Falta validacion de roles del usuario
	return nil
}
