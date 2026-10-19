package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	sqlc "cl_gst/db/sqlc"
	"cl_gst/handlers"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	ctx := context.Background()

	connStr := "host=localhost port=5432 user=admin_tpEspecial password=admin_contra dbname=tpEspecial_db sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Fallo de conexion a la DB: %v", err)
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
