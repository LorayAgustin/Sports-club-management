-- name: CreateSocio :one
INSERT INTO socios (id_Usuario, dni, fechaAlta, fechaBaja)
    VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetSocio :one
SELECT * FROM socios WHERE id_Socio = $1;

-- name: ListSocios :many
SELECT * FROM socios;   

-- name: UpdateSocio :exec
UPDATE socios SET 
        id_Usuario = $2, 
        dni = $3, 
        fechaAlta = $4, 
        fechaBaja = $5
    WHERE id_Socio = $1;

-- name: DeleteSocio :exec
DELETE FROM socios WHERE id_Socio = $1;