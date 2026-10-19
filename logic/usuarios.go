package logic

import (
	"context"
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

type UserDB interface {
	ExistsUserByName(ctx context.Context, Nombre string) (bool, error)
}

func ValidateUsuario(ctx context.Context, usuario Usuario, db UserDB) error {
	if usuario.Nombre == "" {
		return errors.New("el nombre no puede estar vacío")
	}

	exists, err := db.ExistsUserByName(ctx, usuario.Nombre)
	if err != nil {
		return errors.New("error al verificar el usuario en la base de datos")
	}
	if exists {
		return errors.New("el nombre de usuario ya está registrado")
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
	if usuario.Rol != "administrador" && usuario.Rol != "usuario" {
		return errors.New("el rol debe ser 'administrador' o 'usuario'")
	}

	return nil
}
