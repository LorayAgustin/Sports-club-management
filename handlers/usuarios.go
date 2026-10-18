package handlers

import (
	"encoding/json"
	"net/http"

	"Sports-club-management/db/sqlc"
	"Sports-club-management/logic"
)

type UserHandler struct {
	queries *sqlc.Queries
}

func NewUserHandler(queries *sqlc.Queries) *UserHandler {
	return &UserHandler{queries: queries}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input logic.Usuario

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Formato JSON invalido", http.StatusBadRequest)
		return
	}

	if err := logic.ValidateUsuario(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.queries.CreateUsuario(r.Context(), sqlc.CreateUsuarioParams{
		Nombre:          input.Nombre,
		Email:           input.Email,
		Contrasena:      input.Contrasena,
		Fechanacimiento: input.FechaNacimiento,
		Rol:             input.Rol,
	})
	if err != nil {
		http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
