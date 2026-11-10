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

type SocioHandler struct {
	queries *sqlc.Queries
}

func NewSocioHandler(queries *sqlc.Queries) *SocioHandler {
	return &SocioHandler{queries: queries}
}

func (h *SocioHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/socios")
	path = strings.Trim(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.ListSocios(w, r)
		case http.MethodPost:
			h.CreateSocio(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "ID de socio inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetSocioById(w, r, int32(id))
	case http.MethodPut:
		h.UpdateSocio(w, r, int32(id))
	case http.MethodDelete:
		h.DeleteSocio(w, r, int32(id))
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// POST /socios
func (h *SocioHandler) CreateSocio(w http.ResponseWriter, r *http.Request) {
	if h.queries == nil {
		http.Error(w, "Error interno: consultas no inicializado", http.StatusInternalServerError)
		return
	}

	var createSocio logic.Socio
	if err := json.NewDecoder(r.Body).Decode(&createSocio); err != nil {
		http.Error(w, "Formato JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := logic.ValidateSocio(createSocio); err != nil {
		http.Error(w, "Validación de socio fallida: "+err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.queries.GetUsuario(r.Context(), createSocio.IDUsuario)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "El usuario no existe", http.StatusBadRequest)
			return
		}
		http.Error(w, "Error al verificar el usuario", http.StatusInternalServerError)
		return
	}

	var fechaBaja sql.NullTime
	if createSocio.FechaBaja != nil && !createSocio.FechaBaja.IsZero() {
		fechaBaja = sql.NullTime{Time: *createSocio.FechaBaja, Valid: true}
	}

	socio, err := h.queries.CreateSocio(r.Context(), sqlc.CreateSocioParams{
		IDUsuario: createSocio.IDUsuario,
		Dni:       createSocio.DNI,
		Fechaalta: createSocio.FechaAlta,
		Fechabaja: fechaBaja,
	})
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "socios_dni_key") || strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "El DNI ya está en uso", http.StatusBadRequest)
			return
		}
		log.Printf("Error en CreateSocio (DB): %v", err)
		http.Error(w, "Error al crear el socio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(socio)
}

// GET /socios
func (h *SocioHandler) ListSocios(w http.ResponseWriter, r *http.Request) {
	socios, err := h.queries.ListSocios(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener la lista de socios", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(socios)
}

// GET /socios/{id}
func (h *SocioHandler) GetSocioById(w http.ResponseWriter, r *http.Request, id int32) {
	socio, err := h.queries.GetSocio(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Socio no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener el socio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(socio)
}

// PUT /socios/{id}
func (h *SocioHandler) UpdateSocio(w http.ResponseWriter, r *http.Request, id int32) {
	_, err := h.queries.GetSocio(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Socio no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener el socio", http.StatusInternalServerError)
		return
	}

	var updatedSocio logic.Socio
	if err := json.NewDecoder(r.Body).Decode(&updatedSocio); err != nil {
		http.Error(w, "Formato JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := logic.ValidateSocio(updatedSocio); err != nil {
		http.Error(w, "Validación de socio fallida: "+err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.queries.GetUsuario(r.Context(), updatedSocio.IDUsuario)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "El usuario no existe", http.StatusBadRequest)
			return
		}
		http.Error(w, "Error al verificar el usuario", http.StatusInternalServerError)
		return
	}

	var fechaBaja sql.NullTime
	if updatedSocio.FechaBaja != nil && !updatedSocio.FechaBaja.IsZero() {
		fechaBaja = sql.NullTime{Time: *updatedSocio.FechaBaja, Valid: true}
	}

	err = h.queries.UpdateSocio(r.Context(), sqlc.UpdateSocioParams{
		IDSocio:   id,
		IDUsuario: updatedSocio.IDUsuario,
		Dni:       updatedSocio.DNI,
		Fechaalta: updatedSocio.FechaAlta,
		Fechabaja: fechaBaja,
	})
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "socios_dni_key") || strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "El DNI ya está en uso", http.StatusBadRequest)
			return
		}
		http.Error(w, "Error al actualizar el socio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedSocio)
}

// DELETE /socios/{id}
func (h *SocioHandler) DeleteSocio(w http.ResponseWriter, r *http.Request, id int32) {
	_, err := h.queries.GetSocio(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Socio no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener el socio", http.StatusInternalServerError)
		return
	}

	err = h.queries.DeleteSocio(r.Context(), id)
	if err != nil {
		http.Error(w, "Error al eliminar el socio", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
