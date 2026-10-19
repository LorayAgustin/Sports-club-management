package logic

import (
	"context"
	"errors"
	"regexp" // biblioteca para expresiones regulares para uso de validacion de email.
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
	if len(usuario.Nombre) <= 3 {
		return errors.New("el nombre no puede estar vacío")
	}
	exists, err := db.ExistsUserByName(ctx, usuario.Nombre)
	if err != nil {
		return errors.New("error al verificar el usuario en la base de datos")
	}
	if exists {
		return errors.New("el nombre de usuario ya está registrado")
	}
	emailExReg := regexp.MustCompile(`^[a-zA-Z0-9.\_%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$\`)
	email := strings.TrimSpace(usuario.Email) // Eliminamos espacios en blanco al inicio y al final del email
	if email == "" {                          // Verificamos que no esté vacío
		return errors.New("el email del usuario no puede estar vacío")
	}
	if !emailExReg.MatchString(email) { // Validamos el formato (nombre + @ + proveedor + .ext)
		return errors.New("el formato del email es inválido (debe ser tipo ejemplo@dominio.com)")
	}
	if len(usuario.Contrasena) < 8 || strings.Contains(usuario.Contrasena, " ") || !strings.ContainsAny(usuario.Contrasena, "0123456789") || !strings.ContainsAny(usuario.Contrasena, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") || !strings.ContainsAny(usuario.Contrasena, "!@#$%^&*()-+") {
		return errors.New("la contraseña debe tener al menos 8 caracteres, sin espacios y contener al menos un número, una letra mayúscula y un carácter especial")
	}
	if usuario.FechaNacimiento.After(time.Now()) {
		return errors.New("la fecha de nacimiento no puede ser posterior a la fecha actual")
	}
	if usuario.Rol != "administrador" && usuario.Rol != "usuario" {
		return errors.New("el rol debe ser 'administrador' o 'usuario'")
	}
	return nil
}
