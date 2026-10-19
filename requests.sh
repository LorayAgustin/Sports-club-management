#!/bin/bash
BASE_DIR="http://localhost:8080"

echo "Prueba 1: Crear un usuario VALIDO (POST /usuarios) - Se espera: respuesta 201"
curl -i -X POST "$BASE_DIR/usuarios/" \
    -H "Content-Type: application/json" \
    -d '{
        "nombre":"Agustin Loray",
        "email":"agustin@email.com",
        "contrasena":"Contrasena123$",
        "fecha_nacimiento":"2000-01-01T00:00:00Z",
        "rol":"administrador"
    }'

echo -e "\n"
echo "Prueba 2: Crear un usuario NO VALIDO (POST /usuarios) - Se espera: Error 400"
echo "Contaseña invalida"
curl -i -X POST "$BASE_DIR/usuarios/" \
    -H "Content-Type: application/json" \
    -d '{
        "nombre":"Genaro Stampone",
        "email":"genaro@email.com",
        "contrasena":"invalida",
        "fecha_nacimiento":"2000-01-01T00:00:00Z",
        "rol":"usuario"
    }'

echo -e "\n"
echo "Prueba 3: Listar todos los usuarios (GET /usuarios) - Se espera: respuesta 200"
curl -i -X GET "$BASE_DIR/usuarios/"

echo -e "\n"
echo "Prueba 4: Obtener un usuario existente por ID (GET /usuarios/{id}) - Se espera: respuesta 200"
curl -i -X GET "$BASE_DIR/usuarios/1"

echo -e "\n"
echo "Prueba 5: Obtener un usuario inexistente por ID (GET /usuarios/{id}) - Se espera: Error 404"
curl -i -X GET "$BASE_DIR/usuarios/999"

echo -e "\n"
echo "Prueba 6: Actualizar un usuario existente (PUT /usuarios/{id}) - Se espera: respuesta 200"
curl -i -X PUT "$BASE_DIR/usuarios/1" \
    -H "Content-Type: application/json" \
    -d '{
        "nombre":"Agustin Loray",
        "email":"luca.actualizado@email.com",
        "contrasena":"Contrasena123$",
        "fecha_nacimiento":"2001-01-01T00:00:00Z",
        "rol":"administrador"
    }'

echo -e "\n"
echo "Prueba 7: Actualizar un usuario inexistente (PUT /usuarios/{id}) - Se espera: Error 404"
curl -i -X PUT "$BASE_DIR/usuarios/999" \
    -H "Content-Type: application/json" \
    -d '{
        "nombre":"Usuario Inexistente",
        "email":"usuario.inexistente@email.com",
        "contrasena":"Contrasena123$",
        "fecha_nacimiento":"2000-01-01T00:00:00Z",
        "rol":"usuario"
    }'

echo -e "\n"
echo "Prueba 8: Eliminar un usuario existente (DELETE /usuarios/{id}) - Se espera: respuesta 204"
curl -i -X DELETE "$BASE_DIR/usuarios/1"

echo -e "\n"
echo "Prueba 9: Eliminar un usuario inexistente (DELETE /usuarios/{id}) - Se espera: Error 404"
curl -i -X DELETE "$BASE_DIR/usuarios/999"