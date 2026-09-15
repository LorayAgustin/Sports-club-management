-- name: CreateUsuario :one
INSERT INTO usuarios (nombre, email, contrasena, fechaNacimiento, rol)
    VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetUsuario :one
SELECT * FROM usuarios WHERE id = $1;

-- name: ListUsuarios :many
SELECT * FROM usuarios;

-- name: UpdateUsuario :exec
UPDATE usuarios SET 
        nombre = $2, 
        email = $3, 
        contrasena = $4, 
        fechaNacimiento = $5, 
        rol = $6 
    WHERE id = $1; 

-- name: DeleteUsuario :exec
DELETE FROM usuarios WHERE id = $1;