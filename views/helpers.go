package views

import (
	sqlc "cl_gst/db/sqlc"
	"fmt"
)

// ObtenerNombreDeporte busca el nombre descriptivo de un deporte por su ID Deporte
func ObtenerNombreDeporte(idDeporte int32, deportes []sqlc.Deporte) string {
	for _, d := range deportes {
		if d.IDDeporte == idDeporte {
			return d.Nombre
		}
	}
	return fmt.Sprintf("Usuario #%d", idDeporte)
}

// ObtenerNombreCategoria busca el detalle descriptivo de una categoría por su ID
func ObtenerDetalleCategoria(idCat int32, categorias []sqlc.Categoria, deportes []sqlc.Deporte) string {
	for _, c := range categorias {
		if c.IDCategoria == idCat {
			return fmt.Sprintf("Deporte: %s - %s - #%d", ObtenerNombreDeporte(c.IDDeporte, deportes), c.Nombre, c.IDCategoria)
		}
	}
	return "Sin categoría"
}

// ObtenerNombrePredio busca el nombre descriptivo de un predio por su ID
func ObtenerNombrePredio(idPredio int32, predios []sqlc.Predio) string {
	for _, p := range predios {
		if p.IDPredio == idPredio {
			return fmt.Sprintf("%s (#%d)", p.Nombre, p.IDPredio)
		}
	}
	return "Sin predio"
}

// ObtenerNombreUsuario busca el nombre completo de un usuario por su ID
func ObtenerNombreUsuario(idUsuario int32, usuarios []sqlc.Usuario) string {
	for _, u := range usuarios {
		if u.ID == idUsuario {
			return u.Nombre
		}
	}
	return fmt.Sprintf("Usuario #%d", idUsuario)
}

// ObtenerDetalleSocio busca el DNI y datos del socio por su ID
func ObtenerDetalleSocio(idSocio int32, socios []sqlc.Socio, usuarios []sqlc.Usuario) string {
	for _, s := range socios {
		if s.IDSocio == idSocio {
			nombreUser := ObtenerNombreUsuario(s.IDUsuario, usuarios)
			return fmt.Sprintf("DNI: %s - %s - #%d", s.Dni, nombreUser, s.IDSocio)
		}
	}
	return fmt.Sprintf("Socio #%d", idSocio)
}
