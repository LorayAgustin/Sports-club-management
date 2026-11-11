package handlers

import (
	"log"
	"net/http"

	sqlc "cl_gst/db/sqlc"
	"cl_gst/logic"
	"cl_gst/views"
)

type DeporteHandler struct {
	queries *sqlc.Queries
}

func NewDeporteHandler(queries *sqlc.Queries) *DeporteHandler {
	return &DeporteHandler{queries: queries}
}

// Handler central para la ruta /deportes (despacha GET o POST según el método HTTP)
func (h *DeporteHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Index(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// GET /: Muestra la página con el formulario y la tabla de deportes
func (h *DeporteHandler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/deportes" {
		http.NotFound(w, r)
		return
	}

	deportes, err := h.queries.ListDeportes(r.Context())
	if err != nil {
		log.Printf("Error al obtener deportes: %v", err)
		http.Error(w, "Error interno en la base de datos", http.StatusInternalServerError)
		return
	}

	if err := views.DeportesPage("Gestión del Club - Deportes", deportes).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar plantilla: %v", err)
		http.Error(w, "Error interno al generar la vista", http.StatusInternalServerError)
	}
}

// POST /deportes: Guarda un nuevo deporte y redirige a /
func (h *DeporteHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	nombre := r.FormValue("nombre")
	if nombre == "" {
		http.Error(w, "El nombre es obligatorio", http.StatusBadRequest)
		return
	}

	deporteInput := logic.Deporte{
		Nombre: r.FormValue("nombre"),
	}
	// Validar usando la capa /logic
	if err := logic.ValidateDeporte(&deporteInput); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.queries.CreateDeporte(r.Context(), nombre)
	if err != nil {
		log.Printf("Error al crear deporte: %v", err)
		http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
		return
	}

	// Redirección PRG a la página principal
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
