package test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	sqlc "cl_gst/db/sqlc"
)

func TestCRUDPredios(t *testing.T) {
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

	// --- A. CREAR (CreatePredio) ---
	nuevoPredio, err := queries.CreatePredio(ctx, sqlc.CreatePredioParams{
		Nombre:    "Predio Central Test",
		Direccion: "Av. Del Deporte 1234",
	})
	if err != nil {
		t.Fatalf("Fallo CreatePredio: %v", err)
	}
	if nuevoPredio.IDPredio == 0 {
		t.Errorf("Se esperaba un IDPredio mayor a 0")
	}

	// --- B. LEER (GetPredio) ---
	p, err := queries.GetPredio(ctx, nuevoPredio.IDPredio)
	if err != nil {
		t.Fatalf("Fallo GetPredio: %v", err)
	}
	if p.Nombre != "Predio Central Test" {
		t.Errorf("El nombre leído no coincide: %s", p.Nombre)
	}

	// --- C. ACTUALIZAR (UpdatePredio) ---
	err = queries.UpdatePredio(ctx, sqlc.UpdatePredioParams{
		IDPredio:  nuevoPredio.IDPredio,
		Nombre:    "Predio Central Editado",
		Direccion: "Av. Del Deporte 5678",
	})
	if err != nil {
		t.Fatalf("Fallo UpdatePredio: %v", err)
	}

	pEditado, err := queries.GetPredio(ctx, nuevoPredio.IDPredio)
	if err != nil {
		t.Fatalf("Fallo GetPredio tras update: %v", err)
	}
	if pEditado.Nombre != "Predio Central Editado" {
		t.Errorf("No se actualizó correctamente: %s", pEditado.Nombre)
	}

	// --- D. LISTAR (ListPredios) ---
	lista, err := queries.ListPredios(ctx)
	if err != nil {
		t.Fatalf("Fallo ListPredios: %v", err)
	}
	if len(lista) == 0 {
		t.Errorf("Se esperaba al menos un predio en la lista")
	}

	// --- E. ELIMINAR (DeletePredio) ---
	err = queries.DeletePredio(ctx, nuevoPredio.IDPredio)
	if err != nil {
		t.Fatalf("Fallo DeletePredio: %v", err)
	}

	// --- F. CONFIRMAR ELIMINACIÓN ---
	_, err = queries.GetPredio(ctx, nuevoPredio.IDPredio)
	if err != sql.ErrNoRows {
		t.Errorf("Se esperaba sql.ErrNoRows tras eliminar el predio, pero se obtuvo: %v", err)
	}
}
