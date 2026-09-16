-- name: CreateInscripcion :one
INSERT INTO inscripciones (id_Socio, id_Categoria, fechaInscripcion, estado, fechaBaja)
    VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetInscripcion :one
SELECT * FROM inscripciones WHERE id_Inscripcion = $1;

-- name: ListInscripciones :many
SELECT * FROM inscripciones;

-- name: UpdateInscripcion :exec
UPDATE inscripciones SET
        id_Socio = $2,
        id_Categoria = $3,
        fechaInscripcion = $4,
        estado = $5,
        fechaBaja = $6
    WHERE id_Inscripcion = $1;

-- name: DeleteInscripcion :exec
DELETE FROM inscripciones WHERE id_Inscripcion = $1;