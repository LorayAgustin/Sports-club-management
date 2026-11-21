package handlers

import (
	"log"
	"net/http"
	"strconv"

	sqlc "cl_gst/db/sqlc"
	"cl_gst/logic"
	"cl_gst/views"
)

type CategoriaHandler struct {
	queries *sqlc.Queries
}

func NewCategoriaHandler(queries *sqlc.Queries) *CategoriaHandler {
	return &CategoriaHandler{queries: queries}
}

// Handler central para la ruta /categorias (despacha GET o POST según el método HTTP)
func (h *CategoriaHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Index(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// GET /categorias: Muestra la página con el formulario y la tabla de categorías
func (h *CategoriaHandler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/categorias" && r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	categorias, err := h.queries.ListCategorias(r.Context())
	if err != nil {
		log.Printf("Error al obtener categorías: %v", err)
		http.Error(w, "Error interno en la base de datos", http.StatusInternalServerError)
		return
	}

	deportes, err := h.queries.ListDeportes(r.Context())
	if err != nil {
		log.Printf("Error al obtener deportes: %v", err)
		http.Error(w, "Error interno en la base de datos", http.StatusInternalServerError)
		return
	}

	if err := views.CategoriasPage("Gestión del Club - Categorías", categorias, deportes).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar plantilla: %v", err)
		http.Error(w, "Error interno al generar la vista", http.StatusInternalServerError)
	}
}

// POST /categorias: Guarda una nueva categoría y redirige a /categorias
func (h *CategoriaHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	nombre := r.FormValue("nombre")
	idDeporteStr := r.FormValue("id_deporte")

	if nombre == "" || idDeporteStr == "" {
		http.Error(w, "El nombre y el deporte asociado son obligatorios", http.StatusBadRequest)
		return
	}

	idDeporte, err := strconv.Atoi(idDeporteStr)
	if err != nil {
		http.Error(w, "El ID del deporte debe ser un número entero válido", http.StatusBadRequest)
		return
	}

	categoriaInput := logic.Categoria{
		Nombre:    nombre,
		IDDeporte: int(idDeporte),
	}
	// Validar usando la capa /logic
	if err := logic.ValidateCategoria(categoriaInput); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.queries.CreateCategoria(r.Context(), sqlc.CreateCategoriaParams{
		Nombre:    nombre,
		IDDeporte: int32(idDeporte),
	})
	if err != nil {
		log.Printf("Error al crear categoría: %v", err)
		http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
		return
	}

	categorias, err := h.queries.ListCategorias(r.Context())
	if err != nil {
		log.Printf("Error al obtener categorías actualizadas: %v", err)
		http.Error(w, "Error al actualizar la lista", http.StatusInternalServerError)
		return
	}

	deportes, err := h.queries.ListDeportes(r.Context())
	if err != nil {
		log.Printf("Error al obtener deportes: %v", err)
		http.Error(w, "Error al actualizar la lista", http.StatusInternalServerError)
		return
	}

	// Devuelve la lista actualizada a #categorias-list-container
	if err := views.CategoriasList(categorias, deportes).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar lista: %v", err)
		http.Error(w, "Error interno al generar respuesta", http.StatusInternalServerError)
		return
	}

	// Devuelve el formulario limpio OOB a #categoria-form-container
	if err := views.CategoriaFormOOB(deportes).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar formulario OOB: %v", err)
	}
}

// DELETE
func (h *CategoriaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "El ID de la categoría es obligatorio", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.queries.DeleteCategoria(r.Context(), int32(id))
	if err != nil {
		log.Printf("Error al eliminar categoría: %v", err)
		http.Error(w, "Error al eliminar en la base de datos", http.StatusInternalServerError)
		return
	}

	// Responder con 200 OK y cuerpo vacío para que HTMX reemplace 'closest tr' por nada (elimina la fila)
	w.WriteHeader(http.StatusOK)
}
