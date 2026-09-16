-- name: CreateCategoria :one
INSERT INTO categorias (nombre, id_Deporte)   
    VALUES ($1, $2) RETURNING *;

-- name: GetCategoria :one
SELECT * FROM categorias WHERE id_Categoria = $1;

-- name: ListCategorias :many
SELECT * FROM categorias;

-- name: UpdateCategoria :exec
UPDATE categorias SET 
        nombre = $2,
        id_Deporte = $3
    WHERE id_Categoria = $1;

-- name: DeleteCategoria :exec
DELETE FROM categorias WHERE id_Categoria = $1;