package handlers

import (
	"log"
	"net/http"
	"strconv" //se encarga de realizar conversiones entre tipos de datos básicos y sus representaciones en formato de texto (string).
	"time"

	sqlc "cl_gst/db/sqlc"
	"cl_gst/logic"
	"cl_gst/views"
)

type UsuarioHandler struct {
	queries *sqlc.Queries
}

func NewUsuarioHandler(queries *sqlc.Queries) *UsuarioHandler {
	return &UsuarioHandler{queries: queries}
}

// Handle redirige peticiones GET, POST o DELETE
func (h *UsuarioHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

// GET /usuarios: Consulta la base de datos y renderiza la vista
func (h *UsuarioHandler) Index(w http.ResponseWriter, r *http.Request) {
	// Validar que la ruta sea exactamente "/" o "/usuarios"
	if r.URL.Path != "/" && r.URL.Path != "/usuarios" {
		http.NotFound(w, r)
		return
	}

	usuarios, err := h.queries.ListUsuarios(r.Context())
	if err != nil {
		log.Printf("Error al obtener usuarios: %v", err)
		http.Error(w, "Error interno en la base de datos", http.StatusInternalServerError)
		return
	}

	if err := views.UsuariosPage("Gestión del Club - Usuarios", usuarios).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar plantilla: %v", err)
		http.Error(w, "Error interno al generar la vista", http.StatusInternalServerError)
	}
}

// POST /usuarios: Alta vía HTMX
func (h *UsuarioHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Leer los datos enviados por el <form action="/usuarios" method="POST">
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	// Obtener y parsear la fecha de nacimiento
	fechaStr := r.FormValue("fechaNacimiento")
	if fechaStr == "" {
		http.Error(w, "La fecha de nacimiento es obligatoria", http.StatusBadRequest)
		return
	}

	fechaNac, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		http.Error(w, "Formato de fecha inválido", http.StatusBadRequest)
		return
	}

	userInput := logic.Usuario{
		Nombre:          r.FormValue("nombre"),
		Email:           r.FormValue("email"),
		Contrasena:      r.FormValue("contrasena"),
		FechaNacimiento: fechaNac,
		Rol:             r.FormValue("rol"),
	}

	// Validar con las reglas de negocio de la capa /logic
	if err := logic.ValidateUsuario(r.Context(), userInput, h.queries); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Guardar en PostgreSQL usando la consulta generada por sqlc
	_, err = h.queries.CreateUsuario(r.Context(), sqlc.CreateUsuarioParams{
		Nombre:          userInput.Nombre,
		Email:           userInput.Email,
		Contrasena:      userInput.Contrasena,
		Fechanacimiento: userInput.FechaNacimiento,
		Rol:             userInput.Rol,
	})
	if err != nil {
		log.Printf("Error al crear usuario: %v", err)
		http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
		return
	}

	// Consultar la lista actualizada de usuarios
	usuarios, err := h.queries.ListUsuarios(r.Context())
	if err != nil {
		log.Printf("Error al obtener lista de usuarios: %v", err)
		http.Error(w, "Error al actualizar la lista", http.StatusInternalServerError)
		return
	}

	// Renderizar el fragmento de la lista (reemplaza #usuarios-list-container)
	if err := views.UsuariosList(usuarios).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar lista: %v", err)
		http.Error(w, "Error interno al generar respuesta", http.StatusInternalServerError)
		return
	}

	// Renderizar el formulario OOB (reemplaza #usuario-form-container fuera de banda)
	if err := views.UsuarioFormOOB().Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar formulario OOB: %v", err)
	}
}

// DELETE (Borrado dinámico)
func (h *UsuarioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.queries.DeleteUsuario(r.Context(), int32(id)); err != nil {
		log.Printf("Error al eliminar usuario %d: %v", id, err)
		http.Error(w, "Error al eliminar usuario", http.StatusInternalServerError)
		return
	}

	// Devuelve 200 OK con contenido vacío para que HTMX remueva la fila
	w.WriteHeader(http.StatusOK)
}
