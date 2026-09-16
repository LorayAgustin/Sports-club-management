package test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	sqlc "cl_gst/db/sqlc"
)

func TestCRUDCategorias(t *testing.T) {
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

	// 1. Crear un deporte base para la clave foránea
	deporteBase, err := queries.CreateDeporte(ctx, "Deporte Para Categoria Test")
	if err != nil {
		t.Fatalf("Fallo al crear deporte base: %v", err)
	}
	defer queries.DeleteDeporte(ctx, deporteBase.IDDeporte)

	// --- A. CREAR (CreateCategoria) ---
	nuevaCat, err := queries.CreateCategoria(ctx, sqlc.CreateCategoriaParams{
		Nombre:    "Sub-18 Masculino Test",
		IDDeporte: deporteBase.IDDeporte,
	})
	if err != nil {
		t.Fatalf("Fallo CreateCategoria: %v", err)
	}
	if nuevaCat.IDCategoria == 0 {
		t.Errorf("Se esperaba un IDCategoria mayor a 0")
	}

	// --- B. LEER (GetCategoria) ---
	c, err := queries.GetCategoria(ctx, nuevaCat.IDCategoria)
	if err != nil {
		t.Fatalf("Fallo GetCategoria: %v", err)
	}
	if c.Nombre != "Sub-18 Masculino Test" {
		t.Errorf("El nombre de la categoría no coincide: %s", c.Nombre)
	}

	// --- C. ACTUALIZAR (UpdateCategoria) ---
	err = queries.UpdateCategoria(ctx, sqlc.UpdateCategoriaParams{
		IDCategoria: nuevaCat.IDCategoria,
		Nombre:      "Sub-20 Masculino Editado",
		IDDeporte:   deporteBase.IDDeporte,
	})
	if err != nil {
		t.Fatalf("Fallo UpdateCategoria: %v", err)
	}

	cEditada, err := queries.GetCategoria(ctx, nuevaCat.IDCategoria)
	if err != nil {
		t.Fatalf("Fallo GetCategoria tras update: %v", err)
	}
	if cEditada.Nombre != "Sub-20 Masculino Editado" {
		t.Errorf("No se actualizó correctamente: %s", cEditada.Nombre)
	}

	// --- D. LISTAR (ListCategorias) ---
	lista, err := queries.ListCategorias(ctx)
	if err != nil {
		t.Fatalf("Fallo ListCategorias: %v", err)
	}
	if len(lista) == 0 {
		t.Errorf("Se esperaba al menos una categoría en la lista")
	}

	// --- E. ELIMINAR (DeleteCategoria) ---
	err = queries.DeleteCategoria(ctx, nuevaCat.IDCategoria)
	if err != nil {
		t.Fatalf("Fallo DeleteCategoria: %v", err)
	}

	// --- F. CONFIRMAR ELIMINACIÓN ---
	_, err = queries.GetCategoria(ctx, nuevaCat.IDCategoria)
	if err != sql.ErrNoRows {
		t.Errorf("Se esperaba sql.ErrNoRows tras eliminar la categoría, pero se obtuvo: %v", err)
	}
}
