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

func TestCRUDPracticas(t *testing.T) {
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

	// 1. Crear dependencias: Deporte -> Categoria y Predio
	dep, err := queries.CreateDeporte(ctx, "Deporte Practicas Test")
	if err != nil {
		t.Fatalf("Fallo crear deporte: %v", err)
	}
	defer queries.DeleteDeporte(ctx, dep.IDDeporte)

	cat, err := queries.CreateCategoria(ctx, sqlc.CreateCategoriaParams{
		Nombre:    "Categoria Practica Test",
		IDDeporte: dep.IDDeporte,
	})
	if err != nil {
		t.Fatalf("Fallo crear categoria: %v", err)
	}
	defer queries.DeleteCategoria(ctx, cat.IDCategoria)

	pred, err := queries.CreatePredio(ctx, sqlc.CreatePredioParams{
		Nombre:    "Predio Practica Test",
		Direccion: "Calle Practica 123",
	})
	if err != nil {
		t.Fatalf("Fallo crear predio: %v", err)
	}
	defer queries.DeletePredio(ctx, pred.IDPredio)

	// --- A. CREAR (CreatePractica) ---
	horaInicio := time.Now()
	horaFin := time.Now().Add(2 * time.Hour)

	nuevaPrac, err := queries.CreatePractica(ctx, sqlc.CreatePracticaParams{
		IDCategoria: cat.IDCategoria,
		IDPredio:    pred.IDPredio,
		Diasemana:   "Lunes",
		Horainicio:  horaInicio,
		Horafin:     horaFin,
	})
	if err != nil {
		t.Fatalf("Fallo CreatePractica: %v", err)
	}
	if nuevaPrac.IDPractica == 0 {
		t.Errorf("Se esperaba un IDPractica mayor a 0")
	}

	// --- B. LEER (GetPractica) ---
	p, err := queries.GetPractica(ctx, nuevaPrac.IDPractica)
	if err != nil {
		t.Fatalf("Fallo GetPractica: %v", err)
	}
	if p.Diasemana != "Lunes" {
		t.Errorf("El día de la semana no coincide: %s", p.Diasemana)
	}

	// --- C. ACTUALIZAR (UpdatePractica) ---
	err = queries.UpdatePractica(ctx, sqlc.UpdatePracticaParams{
		IDPractica:  nuevaPrac.IDPractica,
		IDCategoria: cat.IDCategoria,
		IDPredio:    pred.IDPredio,
		Diasemana:   "Miércoles",
		Horainicio:  horaInicio,
		Horafin:     horaFin,
	})
	if err != nil {
		t.Fatalf("Fallo UpdatePractica: %v", err)
	}

	pEditada, err := queries.GetPractica(ctx, nuevaPrac.IDPractica)
	if err != nil {
		t.Fatalf("Fallo GetPractica tras update: %v", err)
	}
	if pEditada.Diasemana != "Miércoles" {
		t.Errorf("No se actualizó el día: %s", pEditada.Diasemana)
	}

	// --- D. LISTAR (ListPracticas) ---
	lista, err := queries.ListPracticas(ctx)
	if err != nil {
		t.Fatalf("Fallo ListPracticas: %v", err)
	}
	if len(lista) == 0 {
		t.Errorf("Se esperaba al menos una práctica en la lista")
	}

	// --- E. ELIMINAR (DeletePractica) ---
	err = queries.DeletePractica(ctx, nuevaPrac.IDPractica)
	if err != nil {
		t.Fatalf("Fallo DeletePractica: %v", err)
	}

	// --- F. CONFIRMAR ELIMINACIÓN ---
	_, err = queries.GetPractica(ctx, nuevaPrac.IDPractica)
	if err != sql.ErrNoRows {
		t.Errorf("Se esperaba sql.ErrNoRows tras eliminar la práctica, pero se obtuvo: %v", err)
	}
}
