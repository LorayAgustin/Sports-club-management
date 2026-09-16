package test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	sqlc "cl_gst/db/sqlc"
)

func TestCRUDUsuarios(t *testing.T) {
	// 1. Obtener la cadena de conexión desde las variables de entorno
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		connStr = "postgres://admin_tpEspecial:admin_contra@localhost:5432/tpEspecial_db?sslmode=disable"
	}

	// 2. Abrir conexión con la base de datos de PostgreSQL
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("No se pudo conectar a PostgreSQL: %v", err)
	}

	// 3. Inicializar el cliente de sqlc
	queries := sqlc.New(db)
	ctx := context.Background()

	// --- A. CREAR (CreateUsuario) ---
	fechaNac := time.Date(1998, time.March, 15, 0, 0, 0, 0, time.UTC)
	nuevoUser, err := queries.CreateUsuario(ctx, sqlc.CreateUsuarioParams{
		Nombre:          "Gena Test",
		Email:           "gena_test@club.com",
		Contrasena:      "pass1234",
		Fechanacimiento: fechaNac,
		Rol:             "socio",
	})
	if err != nil {
		t.Fatalf("Fallo CreateUsuario: %v", err)
	}
	if nuevoUser.ID == 0 {
		t.Errorf("Se esperaba un ID autogenerado mayor a 0")
	}

	// --- B. LEER (GetUsuario) ---
	u, err := queries.GetUsuario(ctx, nuevoUser.ID)
	if err != nil {
		t.Fatalf("Fallo GetUsuario: %v", err)
	}
	if u.Nombre != "Gena Test" || u.Email != "gena_test@club.com" {
		t.Errorf("Los datos leídos no coinciden. Obtenido: %+v", u)
	}

	// --- C. ACTUALIZAR (UpdateUsuario) ---
	err = queries.UpdateUsuario(ctx, sqlc.UpdateUsuarioParams{
		ID:              nuevoUser.ID,
		Nombre:          "Gena Test Editado",
		Email:           "gena_editado@club.com",
		Contrasena:      "pass1234_nueva",
		Fechanacimiento: fechaNac,
		Rol:             "admin",
	})
	if err != nil {
		t.Fatalf("Fallo UpdateUsuario: %v", err)
	}

	// Verificación de la actualización
	uEditado, err := queries.GetUsuario(ctx, nuevoUser.ID)
	if err != nil {
		t.Fatalf("Fallo GetUsuario tras update: %v", err)
	}
	if uEditado.Nombre != "Gena Test Editado" || uEditado.Rol != "admin" {
		t.Errorf("No se aplicó la actualización correctamente. Obtenido: %+v", uEditado)
	}

	// --- D. LISTAR (ListUsuarios) ---
	lista, err := queries.ListUsuarios(ctx)
	if err != nil {
		t.Fatalf("Fallo ListUsuarios: %v", err)
	}
	if len(lista) == 0 {
		t.Errorf("Se esperaba al menos un usuario en la lista")
	}

	// --- E. ELIMINAR (DeleteUsuario) ---
	err = queries.DeleteUsuario(ctx, nuevoUser.ID)
	if err != nil {
		t.Fatalf("Fallo DeleteUsuario: %v", err)
	}

	// --- F. CONFIRMAR ELIMINACIÓN ---
	_, err = queries.GetUsuario(ctx, nuevoUser.ID)
	if err != sql.ErrNoRows {
		t.Errorf("Se esperaba sql.ErrNoRows tras eliminar el usuario, pero se obtuvo: %v", err)
	}
}
