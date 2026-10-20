package logic

import (
	"context"
	"errors"
	"regexp" // Biblioteca para expresiones regulares usada para la verificacion del email.
	"strings"
	"time"
)

// Compilación única de la regex al cargar el paquete
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

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
	// 1. Validar Nombre
	nombre := strings.TrimSpace(usuario.Nombre)
	if nombre == "" {
		return errors.New("el nombre del usuario no puede estar vacío")
	}
	if len(nombre) < 3 {
		return errors.New("el nombre debe tener al menos 3 caracteres")
	}

	// 2. Verificar existencia en la DB de forma segura
	if db != nil {
		exists, err := db.ExistsUserByName(ctx, nombre)
		if err != nil {
			return errors.New("error al verificar el usuario en la base de datos")
		}
		if exists {
			return errors.New("el nombre de usuario ya está registrado")
		}
	}

	// 3. Validar Email
	email := strings.TrimSpace(usuario.Email)
	if email == "" {
		return errors.New("el email del usuario no puede estar vacío")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("el formato del email es inválido (debe ser tipo ejemplo@dominio.com)")
	}

	// 4. Validar Contraseña
	if len(usuario.Contrasena) < 8 || strings.Contains(usuario.Contrasena, " ") ||
		!strings.ContainsAny(usuario.Contrasena, "0123456789") ||
		!strings.ContainsAny(usuario.Contrasena, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") ||
		!strings.ContainsAny(usuario.Contrasena, "!@#/$%^&*()-+") {
		return errors.New("la contraseña debe tener al menos 8 caracteres, sin espacios y contener al menos un número, una letra mayúscula y un carácter especial")
	}

	// 5. Validar Fecha de Nacimiento
	if usuario.FechaNacimiento.IsZero() {
		return errors.New("la fecha de nacimiento es requerida y debe enviarse en formato ISO (ej: 2000-01-01T00:00:00Z)")
	}
	if usuario.FechaNacimiento.After(time.Now()) {
		return errors.New("la fecha de nacimiento no puede ser posterior a la fecha actual")
	}

	// 6. Validar Rol
	rol := strings.ToLower(strings.TrimSpace(usuario.Rol))
	if rol != "administrador" && rol != "usuario" && rol != "socio" && rol != "empleado" {
		return errors.New("el rol especificado no es válido")
	}

	return nil
}
