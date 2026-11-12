package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	sqlc "cl_gst/db/sqlc"
	"cl_gst/logic"
	"cl_gst/views"
)

type SocioHandler struct {
	queries *sqlc.Queries
}

func NewSocioHandler(queries *sqlc.Queries) *SocioHandler {
	return &SocioHandler{queries: queries}
}

// Handler central para la ruta /socios (despacha GET o POST según el método HTTP)
func (h *SocioHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Index(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// GET /socios: Consulta la base de datos y renderiza la vista
func (h *SocioHandler) Index(w http.ResponseWriter, r *http.Request) {
	// Validar que la ruta sea exactamente "/" o "/socios"
	if r.URL.Path != "/" && r.URL.Path != "/socios" {
		http.NotFound(w, r)
		return
	}

	socios, err := h.queries.ListSocios(r.Context())
	if err != nil {
		log.Printf("Error al obtener socios: %v", err)
		http.Error(w, "Error interno en la base de datos", http.StatusInternalServerError)
		return
	}

	if err := views.SociosPage("Gestión del Club - Socios", socios).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar plantilla: %v", err)
		http.Error(w, "Error interno al generar la vista", http.StatusInternalServerError)
	}
}

// POST /socios: Procesa los datos del formulario, valida e inserta un nuevo socio
func (h *SocioHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Leer los datos enviados por el <form action="/socios" method="POST">
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	// Parsear el id_Usuario desde el <select name="id_usuario">
	idUsuarioStr := r.FormValue("id_usuario")
	if idUsuarioStr == "" {
		idUsuarioStr = r.FormValue("usuario_id")
	}
	idUsuario, err := strconv.Atoi(idUsuarioStr)
	if err != nil || idUsuario <= 0 {
		http.Error(w, "Debe seleccionar un usuario válido", http.StatusBadRequest)
		return
	}

	// Obtener el DNI desde el <input name="dni">
	dni := r.FormValue("dni")

	// Parsear la fecha de alta desde el <input type="date" name="fechaAlta"> (o time.Now() si no se envía)
	fechaAltaStr := r.FormValue("fechaAlta")
	var fechaAlta time.Time
	if fechaAltaStr != "" {
		var err error
		fechaAlta, err = time.Parse("2006-01-02", fechaAltaStr)
		if err != nil {
			http.Error(w, "Formato de fecha de alta inválido", http.StatusBadRequest)
			return
		}
	} else {
		fechaAlta = time.Now()
	}

	socioInput := logic.Socio{
		IDUsuario: int32(idUsuario),
		DNI:       dni,
		FechaAlta: fechaAlta,
	}

	// Validar con las reglas de negocio de la capa /logic
	if err := logic.ValidateSocio(socioInput); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Guardar en PostgreSQL usando los parámetros generados por sqlc
	_, err = h.queries.CreateSocio(r.Context(), sqlc.CreateSocioParams{
		IDUsuario: int32(idUsuario),
		Dni:       dni,
		Fechaalta: fechaAlta,
	})
	if err != nil {
		log.Printf("Error al crear socio: %v", err)
		http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
		return
	}

	// Redirección PRG (Post/Redirect/Get)
	http.Redirect(w, r, "/socios", http.StatusSeeOther)
}
