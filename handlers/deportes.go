package handlers

import (
	"log"
	"net/http"
	"strconv"

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
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// GET /: Muestra la página con el formulario y la tabla de deportes
func (h *DeporteHandler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/deportes" && r.URL.Path != "/" {
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

	deportes, err := h.queries.ListDeportes(r.Context())
	if err != nil {
		log.Printf("Error al consultar deportes actualizados: %v", err)
		http.Error(w, "Error interno en la base de datos", http.StatusInternalServerError)
		return
	}

	if err := views.DeportesList(deportes).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar plantilla: %v", err)
		http.Error(w, "Error interno al generar la vista", http.StatusInternalServerError)
		return
	}

	if err := views.DeporteFormOOB().Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar formulario OOB: %v", err)
	}
}

// DELETE /deportes?id=X: Elimina una disciplina de la BD y responde 200 OK
func (h *DeporteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		http.Error(w, "ID de deporte no válido", http.StatusBadRequest)
		return
	}

	err = h.queries.DeleteDeporte(r.Context(), int32(id))
	if err != nil {
		log.Printf("Error al eliminar deporte: %v", err)
		http.Error(w, "Error al eliminar en la base de datos", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
