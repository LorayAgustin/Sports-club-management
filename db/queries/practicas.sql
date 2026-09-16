-- name: CreatePractica :one
INSERT INTO practicas (id_Categoria, id_Predio, diaSemana, horaInicio, horaFin)
    VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetPractica :one
SELECT * FROM practicas WHERE id_Practica = $1;

-- name: ListPracticas :many
SELECT * FROM practicas;

-- name: UpdatePractica :exec
UPDATE practicas SET 
        id_Categoria = $2,
        id_Predio = $3,
        diaSemana = $4,
        horaInicio = $5,
        horaFin = $6
    WHERE id_Practica = $1;

-- name: DeletePractica :exec
DELETE FROM practicas WHERE id_Practica = $1;