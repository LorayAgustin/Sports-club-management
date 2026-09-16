-- name: CreatePredio :one
INSERT INTO predios (id_Predio, nombre, direccion)
    VALUES ($1, $2, $3) RETURNING *;

-- name: GetPredio :one
SELECT * FROM predios WHERE id_Predio = $1;

-- name: ListPredios :many
SELECT * FROM predios;

-- name: UpdatePredio :exec
UPDATE predios SET
        nombre = $2,
        direccion = $3
    WHERE id_Predio = $1;

-- name: DeletePredio :exec
DELETE FROM predios WHERE id_Predio = $1;  