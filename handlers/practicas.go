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

type PracticaHandler struct {
	queries *sqlc.Queries
}

func NewPracticaHandler(queries *sqlc.Queries) *PracticaHandler {
	return &PracticaHandler{queries: queries}
}

// Handler central para la ruta /practicas (despacha GET o POST según el método HTTP)
func (h *PracticaHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Index(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// GET /practicas: Muestra la página con el formulario y la tabla de prácticas
func (h *PracticaHandler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/practicas" {
		http.NotFound(w, r)
		return
	}

	practicas, err := h.queries.ListPracticas(r.Context())
	if err != nil {
		log.Printf("Error al obtener prácticas: %v", err)
		http.Error(w, "Error interno en la base de datos", http.StatusInternalServerError)
		return
	}

	categorias, err := h.queries.ListCategorias(r.Context())
	if err != nil {
		log.Printf("Error al obtener categorías: %v", err)
		http.Error(w, "Error interno en la base de datos", http.StatusInternalServerError)
		return
	}

	predios, err := h.queries.ListPredios(r.Context())
	if err != nil {
		log.Printf("Error al obtener predios: %v", err)
		http.Error(w, "Error interno en la base de datos", http.StatusInternalServerError)
		return
	}

	if err := views.PracticasPage("Gestión del Club - Prácticas", practicas, categorias, predios).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar plantilla: %v", err)
		http.Error(w, "Error interno al generar la vista", http.StatusInternalServerError)
	}
}

// POST /practicas: Guarda una nueva práctica y redirige a /practicas
func (h *PracticaHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	idCategoriaStr := r.FormValue("id_categoria")
	idPredioStr := r.FormValue("id_predio")
	diaSemana := r.FormValue("dia_semana")
	horaInicio := r.FormValue("hora_inicio")
	horaFin := r.FormValue("hora_fin")

	if idCategoriaStr == "" || idPredioStr == "" || diaSemana == "" || horaInicio == "" || horaFin == "" {
		http.Error(w, "Todos los campos del formulario son obligatorios", http.StatusBadRequest)
		return
	}

	idCategoria, err := strconv.Atoi(idCategoriaStr)
	if err != nil {
		http.Error(w, "El ID de la categoría debe ser un número entero válido", http.StatusBadRequest)
		return
	}

	idPredio, err := strconv.Atoi(idPredioStr)
	if err != nil {
		http.Error(w, "El ID del predio debe ser un número entero válido", http.StatusBadRequest)
		return
	}
	// 1. Obtener los strings del formulario
	horaInicioStr := r.FormValue("hora_inicio")
	horaFinStr := r.FormValue("hora_fin")

	// 2. Parsear los textos al tipo time.Time usando el formato "15:04"
	horaInicioParsed, err := time.Parse("15:04", horaInicioStr)
	if err != nil {
		http.Error(w, "Formato de hora de inicio inválido (se esperaba HH:MM)", http.StatusBadRequest)
		return
	}

	horaFinParsed, err := time.Parse("15:04", horaFinStr)
	if err != nil {
		http.Error(w, "Formato de hora de fin inválido (se esperaba HH:MM)", http.StatusBadRequest)
		return
	}

	practicaInput := logic.Practica{
		IDCategoria: int(idCategoria),
		IDPredio:    int(idPredio),
		DiaSemana:   diaSemana,
		HoraInicio:  horaInicioParsed,
		HoraFin:     horaFinParsed,
	}

	// Validar usando la capa /logic
	if err := logic.ValidatePractica(practicaInput); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.queries.CreatePractica(r.Context(), sqlc.CreatePracticaParams{
		IDCategoria: int32(idCategoria),
		IDPredio:    int32(idPredio),
		Diasemana:   diaSemana,
		Horainicio:  horaInicioParsed,
		Horafin:     horaFinParsed,
	})

	if err != nil {
		log.Printf("Error al crear práctica: %v", err)
		http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
		return
	}

	// Redirección PRG a la lista de prácticas
	http.Redirect(w, r, "/practicas", http.StatusSeeOther)
}
