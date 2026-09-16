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

func TestCRUDInscripciones(t *testing.T) {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		connStr = "postgres://admin_tpEspecial:admin_contra@localhost:5432/tpEspecial_db?sslmode=disable"
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)
	ctx := context.Background()

	// 1. Crear dependencias (Usuario -> Socio) y (Deporte -> Categoria)
	usr, err := queries.CreateUsuario(ctx, sqlc.CreateUsuarioParams{
		Nombre:          "User Inscripcion Test",
		Email:           "user_inscripcion@club.com",
		Contrasena:      "pass123",
		Fechanacimiento: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Rol:             "socio",
	})
	if err != nil {
		t.Fatalf("Fallo usuario base: %v", err)
	}
	defer queries.DeleteUsuario(ctx, usr.ID)

	socio, err := queries.CreateSocio(ctx, sqlc.CreateSocioParams{
		IDUsuario: usr.ID,
		Dni:       "99887766",
		Fechaalta: time.Now(),
	})
	if err != nil {
		t.Fatalf("Fallo socio base: %v", err)
	}
	defer queries.DeleteSocio(ctx, socio.IDSocio)

	dep, err := queries.CreateDeporte(ctx, "Deporte Inscripcion Test")
	if err != nil {
		t.Fatalf("Fallo deporte base: %v", err)
	}
	defer queries.DeleteDeporte(ctx, dep.IDDeporte)

	cat, err := queries.CreateCategoria(ctx, sqlc.CreateCategoriaParams{
		Nombre:    "Categoria Inscripcion Test",
		IDDeporte: dep.IDDeporte,
	})
	if err != nil {
		t.Fatalf("Fallo categoria base: %v", err)
	}
	defer queries.DeleteCategoria(ctx, cat.IDCategoria)

	// --- A. CREAR (CreateInscripcion) ---
	nuevaInsc, err := queries.CreateInscripcion(ctx, sqlc.CreateInscripcionParams{
		IDSocio:          socio.IDSocio,
		IDCategoria:      cat.IDCategoria,
		Fechainscripcion: time.Now(),
		Estado:           "activa",
		Fechabaja:        sql.NullTime{Valid: false},
	})
	if err != nil {
		t.Fatalf("Fallo CreateInscripcion: %v", err)
	}
	if nuevaInsc.IDInscripcion == 0 {
		t.Errorf("Se esperaba un IDInscripcion mayor a 0")
	}

	// --- B. LEER (GetInscripcion) ---
	ins, err := queries.GetInscripcion(ctx, nuevaInsc.IDInscripcion)
	if err != nil {
		t.Fatalf("Fallo GetInscripcion: %v", err)
	}
	if ins.Estado != "activa" {
		t.Errorf("El estado de la inscripción no coincide: %s", ins.Estado)
	}

	// --- C. ACTUALIZAR (UpdateInscripcion) ---
	err = queries.UpdateInscripcion(ctx, sqlc.UpdateInscripcionParams{
		IDInscripcion:    nuevaInsc.IDInscripcion,
		IDSocio:          socio.IDSocio,
		IDCategoria:      cat.IDCategoria,
		Fechainscripcion: nuevaInsc.Fechainscripcion,
		Estado:           "suspendida",
		Fechabaja:        sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		t.Fatalf("Fallo UpdateInscripcion: %v", err)
	}

	insEditada, err := queries.GetInscripcion(ctx, nuevaInsc.IDInscripcion)
	if err != nil {
		t.Fatalf("Fallo GetInscripcion tras update: %v", err)
	}
	if insEditada.Estado != "suspendida" {
		t.Errorf("No se actualizó el estado: %s", insEditada.Estado)
	}

	// --- D. LISTAR (ListInscripciones) ---
	lista, err := queries.ListInscripciones(ctx)
	if err != nil {
		t.Fatalf("Fallo ListInscripciones: %v", err)
	}
	if len(lista) == 0 {
		t.Errorf("Se esperaba al menos una inscripción en la lista")
	}

	// --- E. ELIMINAR (DeleteInscripcion) ---
	err = queries.DeleteInscripcion(ctx, nuevaInsc.IDInscripcion)
	if err != nil {
		t.Fatalf("Fallo DeleteInscripcion: %v", err)
	}

	// --- F. CONFIRMAR ELIMINACIÓN ---
	_, err = queries.GetInscripcion(ctx, nuevaInsc.IDInscripcion)
	if err != sql.ErrNoRows {
		t.Errorf("Se esperaba sql.ErrNoRows tras eliminar la inscripción, pero se obtuvo: %v", err)
	}
}
