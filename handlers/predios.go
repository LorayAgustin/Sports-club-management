package handlers

import (
	"log"
	"net/http"
	"strconv"

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

// Handler central para la ruta /deportes (despacha GET o POST según el método HTTP)
func (h *PredioHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

// GET /predios: Renderiza la página completa inicial (Layout + Formulario + Lista)
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

	if err := views.PrediosPage("Gestión del Club - Predios", predios).Render(r.Context(), w); err != nil {
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
		Nombre:    predioInput.Nombre,
		Direccion: predioInput.Direccion,
	})
	if err != nil {
		log.Printf("Error al crear predio: %v", err)
		http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
		return
	}

	// Consultar la lista actualizada de predios
	predios, err := h.queries.ListPredios(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener la lista de predios", http.StatusInternalServerError)
		return
	}

	// Renderizar el 1er fragmento HTML: Lista actualizada (apunta a hx-target="#predios-list")
	if err := views.PrediosList(predios).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar lista parcial de predios: %v", err)
		http.Error(w, "Error al actualizar la lista", http.StatusInternalServerError)
		return
	}

	// Renderizar el 2do fragmento HTML: Formulario limpio OOB (reemplaza #predio-form)
	if err := views.PredioFormOOB().Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar formulario OOB: %v", err)
	}
}

// DELETE /predios/{id}: Borra en BD y devuelve 200 OK vacío para que HTMX remueva la fila
func (h *PredioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		http.Error(w, "ID de predio no válido", http.StatusBadRequest)
		return
	}
	err = h.queries.DeletePredio(r.Context(), int32(id))
	if err != nil {
		log.Printf("Error al eliminar predio: %v", err)
		http.Error(w, "Error al eliminar en la base de datos", http.StatusInternalServerError)
		return
	}

	// Respuesta vacía con estado 200 OK
	w.WriteHeader(http.StatusOK)
}
