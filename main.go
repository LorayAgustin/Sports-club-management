package main

import (
	"context"
	"database/sql"
	"log"

	sqlc "cl_gst/db/sqlc"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	connStr := "host=localhost port=5432 user=admin_tpEspecial password=admin_contra dbname=tpEspecial_db sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Fallo de conexion a la DB: %v", err)
	}

	defer db.Close()

	ctx := context.Background()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Fallo de ping a la DB: %v", err)
	}

	log.Println("Conexión exitosa a la base de datos")

	queries := sqlc.New(db)

}
