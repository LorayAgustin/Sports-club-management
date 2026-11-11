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

	if err := views.IndexPage("Gestión del Club - Categorías", categorias).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar plantilla: %v", err)
		http.Error(w, "Error interno al generar la vista", http.StatusInternalServerError)
	}
}

// POST /categorias: Guarda una nueva categoría y redirige a /categorias
func (h *CategoriaHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

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

	// Redirección PRG a la lista de categorías
	http.Redirect(w, r, "/categorias", http.StatusSeeOther)
}
