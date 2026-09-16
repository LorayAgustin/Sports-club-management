package test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	sqlc "cl_gst/db/sqlc"
)

func TestCRUDDeportes(t *testing.T) {
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

	// --- A. CREAR (CreateDeporte) ---
	nuevoDeporte, err := queries.CreateDeporte(ctx, "Basquetbol Test")
	if err != nil {
		t.Fatalf("Fallo CreateDeporte: %v", err)
	}
	if nuevoDeporte.IDDeporte == 0 {
		t.Errorf("Se esperaba un IDDeporte autogenerado mayor a 0")
	}

	// --- B. LEER (GetDeporte) ---
	d, err := queries.GetDeporte(ctx, nuevoDeporte.IDDeporte)
	if err != nil {
		t.Fatalf("Fallo GetDeporte: %v", err)
	}
	if d.Nombre != "Basquetbol Test" {
		t.Errorf("El nombre del deporte leído no coincide. Obtenido: %s", d.Nombre)
	}

	// --- C. ACTUALIZAR (UpdateDeporte) ---
	err = queries.UpdateDeporte(ctx, sqlc.UpdateDeporteParams{
		IDDeporte: nuevoDeporte.IDDeporte,
		Nombre:    "Basquetbol Editado",
	})
	if err != nil {
		t.Fatalf("Fallo UpdateDeporte: %v", err)
	}

	dEditado, err := queries.GetDeporte(ctx, nuevoDeporte.IDDeporte)
	if err != nil {
		t.Fatalf("Fallo GetDeporte tras update: %v", err)
	}
	if dEditado.Nombre != "Basquetbol Editado" {
		t.Errorf("No se aplicó la actualización correctamente. Obtenido: %s", dEditado.Nombre)
	}

	// --- D. LISTAR (ListDeportes) ---
	lista, err := queries.ListDeportes(ctx)
	if err != nil {
		t.Fatalf("Fallo ListDeportes: %v", err)
	}
	if len(lista) == 0 {
		t.Errorf("Se esperaba al menos un deporte en la lista")
	}

	// --- E. ELIMINAR (DeleteDeporte) ---
	err = queries.DeleteDeporte(ctx, nuevoDeporte.IDDeporte)
	if err != nil {
		t.Fatalf("Fallo DeleteDeporte: %v", err)
	}

	// --- F. CONFIRMAR ELIMINACIÓN ---
	_, err = queries.GetDeporte(ctx, nuevoDeporte.IDDeporte)
	if err != sql.ErrNoRows {
		t.Errorf("Se esperaba sql.ErrNoRows tras eliminar el deporte, pero se obtuvo: %v", err)
	}
}
