package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	sqlc "cl_gst/db/sqlc"
	"cl_gst/handlers"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	ctx := context.Background()

	// 1. Leer la cadena de conexión desde la variable de entorno DB_URL (Docker)
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		// Fallback para ejecución local fuera de Docker
		connStr = "postgres://admin_tpEspecial:admin_contra@localhost:5432/tpEspecial_db?sslmode=disable"
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Fallo de conexión a la DB: %v", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Fallo de ping a la DB: %v", err)
	}

	log.Println("Conexión exitosa a la base de datos")

	queries := sqlc.New(db)
	usuarioHandler := handlers.NewUserHandler(queries)

	http.Handle("/usuarios/", usuarioHandler)

	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
