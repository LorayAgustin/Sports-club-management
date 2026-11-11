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

	// Leer la cadena de conexión desde la variable de entorno DB_URL (Docker)
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

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Fallo de ping a la DB: %v", err)
	}

	log.Println("Conexión exitosa a la base de datos")

	queries := sqlc.New(db)

	usuarioHandler := handlers.NewUserHandler(queries)
	deporteHandler := handlers.NewDeporteHandler(queries)
	prediosHandler := handlers.NewPredioHandler(queries)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Usuarios
	http.Handle("/usuarios", usuarioHandler)
	http.Handle("/usuarios/", usuarioHandler)

	// Deportes
	// Ruta principal GET / para listar y mostrar el formulario
	http.HandleFunc("/", deporteHandler.Index)
	// Ruta POST /deportes para recibir el envío del formulario
	http.HandleFunc("/deportes", deporteHandler.Create)

	//Predios
	// Ruta principal GET / para listar y mostrar el formulario
	http.HandleFunc("/", prediosHandler.Create)
	// Ruta POST /deporte para recibir el envio del formulario
	http.HandleFunc("/predios", prediosHandler.Create)

	log.Println("Servidor escuchando en http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
