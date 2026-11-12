package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	sqlc "cl_gst/db/sqlc"
	"cl_gst/logic"
	"cl_gst/views"
)

type InscripcionHandler struct {
	queries *sqlc.Queries
}

func NewInscripcionHandler(queries *sqlc.Queries) *InscripcionHandler {
	return &InscripcionHandler{queries: queries}
}

// Handler central para la ruta /inscripciones (despacha GET o POST según el método HTTP)
func (h *InscripcionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Index(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// GET /inscripciones: Muestra la página con el formulario y la tabla de inscripciones
func (h *InscripcionHandler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/inscripciones" {
		http.NotFound(w, r)
		return
	}

	inscripciones, err := h.queries.ListInscripciones(r.Context())
	if err != nil {
		log.Printf("Error al obtener inscripciones: %v", err)
		http.Error(w, "Error interno en la base de datos", http.StatusInternalServerError)
		return
	}

	socios, err := h.queries.ListSocios(r.Context())
	if err != nil {
		log.Printf("Error al obtener socios: %v", err)
		http.Error(w, "Error interno en la base de datos", http.StatusInternalServerError)
		return
	}

	categorias, err := h.queries.ListCategorias(r.Context())
	if err != nil {
		log.Printf("Error al obtener categorías: %v", err)
		http.Error(w, "Error interno en la base de datos", http.StatusInternalServerError)
		return
	}

	if err := views.InscripcionesPage("Gestión del Club - Inscripciones", inscripciones, socios, categorias).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar plantilla: %v", err)
		http.Error(w, "Error interno al generar la vista", http.StatusInternalServerError)
	}
}

// POST /inscripciones: Guarda una nueva inscripción y redirige a /inscripciones
func (h *InscripcionHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	idSocioStr := r.FormValue("id_socio")
	idCategoriaStr := r.FormValue("id_categoria")
	fechaInscripcionStr := r.FormValue("fecha_inscripcion")
	estado := r.FormValue("estado")

	if idSocioStr == "" || idCategoriaStr == "" || fechaInscripcionStr == "" || estado == "" {
		http.Error(w, "Los campos socio, categoría, fecha de inscripción y estado son obligatorios", http.StatusBadRequest)
		return
	}

	idSocio, err := strconv.Atoi(idSocioStr)
	if err != nil {
		http.Error(w, "El ID del socio debe ser un número entero válido", http.StatusBadRequest)
		return
	}

	idCategoria, err := strconv.Atoi(idCategoriaStr)
	if err != nil {
		http.Error(w, "El ID de la categoría debe ser un número entero válido", http.StatusBadRequest)
		return
	}

	fechaInscripcion, err := time.Parse("2006-01-02", fechaInscripcionStr)
	if err != nil {
		http.Error(w, "Formato de fecha de inscripción inválido", http.StatusBadRequest)
		return
	}

	inscripcionInput := logic.Inscripcion{
		IDSocio:          int32(idSocio),
		IDCategoria:      int32(idCategoria),
		FechaInscripcion: fechaInscripcion,
		Estado:           estado,
	}

	// Validar usando la capa /logic
	if err := logic.ValidateInscripcion(inscripcionInput); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.queries.CreateInscripcion(r.Context(), sqlc.CreateInscripcionParams{
		IDSocio:          int32(idSocio),
		IDCategoria:      int32(idCategoria),
		Fechainscripcion: fechaInscripcion,
		Estado:           estado,
	})
	if err != nil {
		log.Printf("Error al crear inscripción: %v", err)
		http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
		return
	}

	// Redirección PRG a la lista de inscripciones
	http.Redirect(w, r, "/inscripciones", http.StatusSeeOther)
}
