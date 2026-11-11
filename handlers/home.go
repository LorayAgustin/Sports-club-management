package handlers

import (
	"log"
	"net/http"

	"cl_gst/views"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

// GET /: Renderiza el menú principal sin consultar la BD
func (h *HomeHandler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if err := views.IndexPage("Gestión del Club - Inicio").Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar inicio: %v", err)
		http.Error(w, "Error interno al renderizar la vista", http.StatusInternalServerError)
	}
}
