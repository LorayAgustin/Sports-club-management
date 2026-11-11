package handlers

import (
	"log"
	"net/http"

	sqlc "cl_gst/db/sqlc"
	"cl_gst/logic"
	"cl_gst/views"
)

type PredioHandler struct {
	queries *sqlc.Queries
}

func NewPredioHandler(queries *sqlc.Queries) *PredioHandler {
	return &PredioHandler{queries: queries}
}

// GET /predios: Muestra la página con el formulario y la tabla de predios
func (h *PredioHandler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/predios" && r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	predios, err := h.queries.ListPredios(r.Context())
	if err != nil {
		log.Printf("Error al obtener predios: %v", err)
		http.Error(w, "Error interno en la base de datos", http.StatusInternalServerError)
		return
	}

	if err := views.IndexPage("Gestión del Club - Predios", predios).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar plantilla: %v", err)
		http.Error(w, "Error interno al generar la vista", http.StatusInternalServerError)
	}
}

// POST /predios: Guarda un nuevo predio y redirige a /predios
func (h *PredioHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	nombre := r.FormValue("nombre")
	direccion := r.FormValue("direccion")

	if nombre == "" || direccion == "" {
		http.Error(w, "El nombre y la dirección son obligatorios", http.StatusBadRequest)
		return
	}

	predioInput := logic.Predio{
		Nombre:    nombre,
		Direccion: direccion,
	}

	// Validar usando la capa /logic
	if err := logic.ValidatePredio(predioInput); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.queries.CreatePredio(r.Context(), sqlc.CreatePredioParams{
		Nombre:    nombre,
		Direccion: direccion,
	})
	if err != nil {
		log.Printf("Error al crear predio: %v", err)
		http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
		return
	}

	// Redirección PRG a la lista de predios
	http.Redirect(w, r, "/predios", http.StatusSeeOther)
}
