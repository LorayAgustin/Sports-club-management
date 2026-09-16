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

func TestCRUDSocios(t *testing.T) {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		connStr = "postgres://admin_tpEspecial:admin_contra@localhost:5432/tpEspecial_db?sslmode=disable"
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("No se pudo conectar a PostgreSQL: %v", err)
	}

	queries := sqlc.New(db)
	ctx := context.Background()

	// 1. Crear un usuario de prueba para cumplir la clave foránea id_Usuario
	userBase, err := queries.CreateUsuario(ctx, sqlc.CreateUsuarioParams{
		Nombre:          "Usuario Para Socio",
		Email:           "socio_test@club.com",
		Contrasena:      "pass123",
		Fechanacimiento: time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC),
		Rol:             "socio",
	})
	if err != nil {
		t.Fatalf("Fallo al crear usuario base para socio: %v", err)
	}
	defer queries.DeleteUsuario(ctx, userBase.ID)

	// --- A. CREAR (CreateSocio) ---
	fechaAlta := time.Now()
	nuevoSocio, err := queries.CreateSocio(ctx, sqlc.CreateSocioParams{
		IDUsuario: userBase.ID,
		Dni:       "12345678",
		Fechaalta: fechaAlta,
	})
	if err != nil {
		t.Fatalf("Fallo CreateSocio: %v", err)
	}
	if nuevoSocio.IDSocio == 0 {
		t.Errorf("Se esperaba un IDSocio autogenerado mayor a 0")
	}

	// --- B. LEER (GetSocio) ---
	s, err := queries.GetSocio(ctx, nuevoSocio.IDSocio)
	if err != nil {
		t.Fatalf("Fallo GetSocio: %v", err)
	}
	if s.Dni != "12345678" || s.IDUsuario != userBase.ID {
		t.Errorf("Los datos leídos del socio no coinciden. Obtenido: %+v", s)
	}

	// --- C. ACTUALIZAR (UpdateSocio) ---
	err = queries.UpdateSocio(ctx, sqlc.UpdateSocioParams{
		IDSocio:   nuevoSocio.IDSocio,
		IDUsuario: userBase.ID,
		Dni:       "87654321",
		Fechaalta: fechaAlta,
	})
	if err != nil {
		t.Fatalf("Fallo UpdateSocio: %v", err)
	}

	sEditado, err := queries.GetSocio(ctx, nuevoSocio.IDSocio)
	if err != nil {
		t.Fatalf("Fallo GetSocio tras update: %v", err)
	}
	if sEditado.Dni != "87654321" {
		t.Errorf("No se aplicó la actualización correctamente. Obtenido: %+v", sEditado)
	}

	// --- D. LISTAR (ListSocios) ---
	lista, err := queries.ListSocios(ctx)
	if err != nil {
		t.Fatalf("Fallo ListSocios: %v", err)
	}
	if len(lista) == 0 {
		t.Errorf("Se esperaba al menos un socio en la lista")
	}

	// --- E. ELIMINAR (DeleteSocio) ---
	err = queries.DeleteSocio(ctx, nuevoSocio.IDSocio)
	if err != nil {
		t.Fatalf("Fallo DeleteSocio: %v", err)
	}

	// --- F. CONFIRMAR ELIMINACIÓN ---
	_, err = queries.GetSocio(ctx, nuevoSocio.IDSocio)
	if err != sql.ErrNoRows {
		t.Errorf("Se esperaba sql.ErrNoRows tras eliminar el socio, pero se obtuvo: %v", err)
	}
}
