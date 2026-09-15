-- name: CreateDeporte :one
INSERT INTO deportes (nombre)   
    VALUES ($1) RETURNING *;

-- name: GetDeporte :one
SELECT * FROM deportes WHERE id_Deporte = $1;

-- name: ListDeportes :many
SELECT * FROM deportes;

-- name: UpdateDeporte :exec
UPDATE deportes SET 
        nombre = $2 
    WHERE id_Deporte = $1;

-- name: DeleteDeporte :exec
DELETE FROM deportes WHERE id_Deporte = $1;