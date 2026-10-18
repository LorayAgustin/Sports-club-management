package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"Sports-club-management/db/sqlc"
	"Sports-club-management/logic"
)

type UserHandler struct {
	queries *sqlc.Queries
}

func NewUserHandler(queries *sqlc.Queries) *UserHandler {
	return &UserHandler{queries: queries}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/usuarios")
	path = strings.Trim(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.ListUsers(w, r)
		case http.MethodPost:
			h.CreateUser(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetUserById(w, r, int32(id))
	case http.MethodPut:
		h.UpdateUser(w, r, int32(id))
	case http.MethodDelete:
		h.DeleteUser(w, r, int32(id))
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// GET /usuarios
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.queries.ListUsuarios(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener la lista de usuarios", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// POST /usuarios
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var createUsuario logic.Usuario

	if err := json.NewDecoder(r.Body).Decode(&createUsuario); err != nil {
		// Error 400: Formato JSON inválido
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}

	if err := logic.ValidateUsuario(createUsuario); err != nil {
		// Error 400: Validación de datos fallida
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.queries.CreateUsuario(r.Context(), sqlc.CreateUsuarioParams{
		Nombre:          createUsuario.Nombre,
		Email:           createUsuario.Email,
		Contrasena:      createUsuario.Contrasena,
		Fechanacimiento: createUsuario.FechaNacimiento,
		Rol:             createUsuario.Rol,
	})
	if err != nil {
		// Error 500: Error interno del servidor
		http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
		return
	}
	// Estado 201: Usuario creado exitosamente
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GET /usuarios/{id}
func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request, id int32) {
	user, err := h.queries.GetUsuario(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Error 404: Usuario no encontrado
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}
		// Error 500: Error interno del servidor
		http.Error(w, "Error al obtener el usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// PUT /usuarios/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, id int32) {
	var updatedUsuario logic.Usuario

	if err := json.NewDecoder(r.Body).Decode(&updatedUsuario); err != nil {
		// Error 400: Formato JSON inválido
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}

	if err := logic.ValidateUsuario(updatedUsuario); err != nil {
		// Error 400: Validación de datos fallida
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.queries.UpdateUsuario(r.Context(), sqlc.UpdateUsuarioParams{
		ID:              id,
		Nombre:          updatedUsuario.Nombre,
		Email:           updatedUsuario.Email,
		Contrasena:      updatedUsuario.Contrasena,
		Fechanacimiento: updatedUsuario.FechaNacimiento,
		Rol:             updatedUsuario.Rol,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Error 404: Usuario no encontrado
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}
		// Error 500: Error interno del servidor
		http.Error(w, "Error al actualizar el usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// DELETE /usuarios/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, id int32) {
	err := h.queries.DeleteUsuario(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Error 404: Usuario no encontrado
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}
		// Error 500: Error interno del servidor
		http.Error(w, "Error al eliminar el usuario", http.StatusInternalServerError)
		return
	}
	// Estado 204: Usuario eliminado exitosamente
	w.WriteHeader(http.StatusNoContent)
}
