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

type DeporteHandler struct {
	queries *sqlc.Queries
}

func NewDeporteHandler(queries *sqlc.Queries) *DeporteHandler {
	return &DeporteHandler{queries: queries}
}

func (h *DeporteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/deportes")
	path = strings.Trim(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.ListDeporte(w, r)
		case http.MethodPost:
			h.CreateDeporte(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "ID de deporte inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetDeporteById(w, r, int32(id))
	case http.MethodPut:
		h.UpdateDeporte(w, r, int32(id))
	case http.MethodDelete:
		h.DeleteDeporte(w, r, int32(id))
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// POST /deportes
func (h *DeporteHandler) CreateDeporte(w http.ResponseWriter, r *http.Request) {
	if h.queries == nil {
		http.Error(w, "Error interno: consultas no inicializado", http.StatusInternalServerError)
		return
	}

	var createDeporte logic.Deporte
	if err := json.NewDecoder(r.Body).Decode(&createDeporte); err != nil {
		http.Error(w, "Formato JSON inválido: ", http.StatusBadRequest)
		return
	}

	if err := logic.ValidateDeporte(createDeporte); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deporte, err := h.queries.CreateDeporte(r.Context(), createDeporte.Nombre)
	if err != nil {
		log.Printf("Error en CreateDeporte (DB): %v", err)
		http.Error(w, "Error al crear el deporte: ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(deporte)
}

// GET /deportes
func (h *DeporteHandler) ListDeporte(w http.ResponseWriter, r *http.Request) {
	deportes, err := h.queries.ListDeportes(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener la lista de deportes: ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deportes)
}

// DELETE /deportes/{id}
func (h *DeporteHandler) DeleteDeporte(w http.ResponseWriter, r *http.Request, id int32) {
	_, err := h.queries.GetDeporte(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Deporte no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener el deporte: ", http.StatusInternalServerError)
		return
	}

	if err := h.queries.DeleteDeporte(r.Context(), id); err != nil {
		http.Error(w, "Error al eliminar el deporte: ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /deportes/{id}
func (h *DeporteHandler) GetDeporteById(w http.ResponseWriter, r *http.Request, id int32) {
	deporte, err := h.queries.GetDeporte(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Deporte no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener el deporte: ", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deporte)
}

// PUT /deportes/{id}
func (h *DeporteHandler) UpdateDeporte(w http.ResponseWriter, r *http.Request, id int32) {
	_, err := h.queries.GetDeporte(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Deporte no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener el deporte: ", http.StatusInternalServerError)
		return
	}

	var updatedDeporte logic.Deporte
	if err := json.NewDecoder(r.Body).Decode(&updatedDeporte); err != nil {
		http.Error(w, "Formato JSON inválido: ", http.StatusBadRequest)
		return
	}

	if err := logic.ValidateDeporte(updatedDeporte); err != nil {
		http.Error(w, "Error de validación: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = h.queries.UpdateDeporte(r.Context(), sqlc.UpdateDeporteParams{
		IDDeporte: id,
		Nombre:    updatedDeporte.Nombre,
	})
	if err != nil {
		http.Error(w, "Error al actualizar el deporte: ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedDeporte)
}
