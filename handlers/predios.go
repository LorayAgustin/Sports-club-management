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

type PredioHandler struct {
	queries *sqlc.Queries
}

func NewPredioHandler(queries *sqlc.Queries) *PredioHandler {
	return &PredioHandler{queries: queries}
}

func (h *PredioHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/predios")
	path = strings.Trim(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.ListPredios(w, r)
		case http.MethodPost:
			h.CreatePredio(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "ID de predio inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetPredioById(w, r, int32(id))
	case http.MethodPut:
		h.UpdatePredio(w, r, int32(id))
	case http.MethodDelete:
		h.DeletePredio(w, r, int32(id))
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// POST /predios
func (h *PredioHandler) CreatePredio(w http.ResponseWriter, r *http.Request) {
	if h.queries == nil {
		http.Error(w, "Error interno del servidor: consultas no inicializadas", http.StatusInternalServerError)
		return
	}

	var createPredio logic.Predio
	if err := json.NewDecoder(r.Body).Decode(&createPredio); err != nil {
		http.Error(w, "Formato JSON inválido: ", http.StatusBadRequest)
		return
	}

	if err := logic.ValidatePredio(createPredio); err != nil {
		http.Error(w, "Error de validación: "+err.Error(), http.StatusBadRequest)
		return
	}

	predio, err := h.queries.CreatePredio(r.Context(), sqlc.CreatePredioParams{
		Nombre:    createPredio.Nombre,
		Direccion: createPredio.Direccion,
	})
	if err != nil {
		log.Printf("Error al crear el predio: %v", err)
		http.Error(w, "Error al crear el predio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(predio)
}

// GET /predios
func (h *PredioHandler) ListPredios(w http.ResponseWriter, r *http.Request) {
	predios, err := h.queries.ListPredios(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener la lista de predios", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(predios)
}

// GET /predios/{id}
func (h *PredioHandler) GetPredioById(w http.ResponseWriter, r *http.Request, id int32) {
	predio, err := h.queries.GetPredio(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Predio no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener el predio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(predio)
}

// PUT /predios/{id}
func (h *PredioHandler) UpdatePredio(w http.ResponseWriter, r *http.Request, id int32) {
	_, err := h.queries.GetPredio(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Predio no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener el predio", http.StatusInternalServerError)
		return
	}

	var updatedPredio logic.Predio
	if err := json.NewDecoder(r.Body).Decode(&updatedPredio); err != nil {
		http.Error(w, "Formato JSON inválido: ", http.StatusBadRequest)
		return
	}

	if err := logic.ValidatePredio(updatedPredio); err != nil {
		http.Error(w, "Error de validación: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = h.queries.UpdatePredio(r.Context(), sqlc.UpdatePredioParams{
		IDPredio:  id,
		Nombre:    updatedPredio.Nombre,
		Direccion: updatedPredio.Direccion,
	})
	if err != nil {
		http.Error(w, "Error al actualizar el predio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedPredio)
}

// DELETE /predios/{id}
func (h *PredioHandler) DeletePredio(w http.ResponseWriter, r *http.Request, id int32) {
	_, err := h.queries.GetPredio(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Predio no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener el predio", http.StatusInternalServerError)
		return
	}

	if err := h.queries.DeletePredio(r.Context(), id); err != nil {
		http.Error(w, "Error al eliminar el predio", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
