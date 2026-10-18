package logic

import (
	"errors"
	"strings"
	"time"
)

type Usuario struct {
	nombre          string
	email           string
	contrasena      string
	fechaNacimiento time.Time
	rol             string
}

func ValidateUsuario(usuario Usuario) error {
	if usuario.nombre == "" {
		return errors.New("el nombre no puede estar vacío")
	}
	if usuario.email == "" {
		return errors.New("el email no puede estar vacío")
	}
	if !strings.Contains(usuario.email, "@") {
		return errors.New("el email debe contener un '@'")
	}
	if len(usuario.contrasena) < 8 {
		return errors.New("la contraseña debe tener al menos 8 caracteres")
	}
	if usuario.fechaNacimiento.After(time.Now()) {
		return errors.New("la fecha de nacimiento no puede ser posterior a la fecha actual")
	}
	//Falta validacion de roles del usuario
	return nil
}
