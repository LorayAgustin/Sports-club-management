package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	sqlc "cl_gst/db/sqlc"
	"cl_gst/logic"
)

type CategoriaHandler struct {
	queries *sqlc.Queries
}

func NewCategoriaHandler(queries *sqlc.Queries) *CategoriaHandler {
	return &CategoriaHandler{queries: queries}
}

func (h *CategoriaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/categorias")
	path = strings.Trim(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.ListCategorias(w, r)
		case http.MethodPost:
			h.CreateCategoria(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetCategoriaById(w, r, int32(id))
	case http.MethodPut:
		h.UpdateCategoria(w, r, int32(id))
	case http.MethodDelete:
		h.DeleteCategoria(w, r, int32(id))
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// GET /categorias
func (h *CategoriaHandler) ListCategorias(w http.ResponseWriter, r *http.Request) {
	categorias, err := h.queries.ListCategorias(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener la lista de categorías", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categorias)
}

// POST /categorias
func (h *CategoriaHandler) CreateCategoria(w http.ResponseWriter, r *http.Request) {
	if h.queries == nil {
		http.Error(w, "Queries no inicializadas", http.StatusInternalServerError)
		return
	}

	var createCategoria logic.Cateagoria
	if err := json.NewDecoder(r.Body).Decode(&createCategoria); err != nil {
		http.Error(w, "Formato JSON invalido:", http.StatusBadRequest)
		return
	}

	if err := logic.ValidateCategoria(createCategoria); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	categoria, err := h.queries.CreateCategoria(r.Context(), sqlc.CreateCategoriaParams{
		Nombre:    createCategoria.Nombre,
		IDDeporte: int32(createCategoria.IDDeporte),
	})
	if err != nil {
		log.Printf("Error en CreateCategoria (DB): %v", err)
		http.Error(w, "Error al crear la categoría", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(categoria)
}

// GET /categorias/{id}
func (h *CategoriaHandler) GetCategoriaById(w http.ResponseWriter, r *http.Request, id int32) {
	categoria, err := h.queries.GetCategoria(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Categoría no encontrada", http.StatusNotFound)
		}
		http.Error(w, "Error al obtener la categoría", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoria)
}

// PUT /categorias/{id}
func (h *CategoriaHandler) UpdateCategoria(w http.ResponseWriter, r *http.Request, id int32) {

	categoriaActual, err := h.queries.GetCategoria(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Categoría no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al verificar la categoría en la base de datos", http.StatusInternalServerError)
		return
	}

	var updatedCategoria logic.Categoria
	if err := json.NewDecoder(r.Body).Decode(&updatedCategoria); err != nil {
		http.Error(w, "Formato JSON inválido: ", http.StatusBadRequest)
		return
	}

	if err := logic.ValidateCategoria(updatedCategoria); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.queries.UpdateCategoria(r.Context(), sqlc.UpdateCategoriaParams{
		IDCategoria: id,
		Nombre:      updatedCategoria.Nombre,
		IDDeporte:   int32(updatedCategoria.IDDeporte),
	})
	if err != nil {

		http.Error(w, "Error al actualizar la categoría", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedCategoria)
}

// DELETE /categorias/{id}
func (h *CategoriaHandler) DeleteCategoria(w http.ResponseWriter, r *http.Request, id int32) {
	_, err := h.queries.GetCategoria(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Categoría no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener la categoría", http.StatusInternalServerError)
		return
	}

	err = h.queries.DeleteCategoria(r.Context(), id)
	if err != nil {
		http.Error(w, "Error al eliminar la categoría", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
