#!/bin/bash
BASE_DIR="http://localhost:8080"

echo "=================================================="
echo "=== 1. PRUEBAS DE CREACIÓN (POST /usuarios/) ==="
echo "=================================================="

echo "Prueba 1: Crear un usuario VÁLIDO - Se espera: 201 Created"
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

echo "Prueba 2: Intentar crear usuario DUPLICADO - Se espera: 400 Bad Request"
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

echo "Prueba 3: Crear usuario con datos INVÁLIDOS - Se espera: 400 Bad Request"
curl -i -X POST "$BASE_DIR/usuarios/" \
  -H "Content-Type: application/json" \
  -d '{
    "nombre":"Genaro Stampone",
    "email":"email-invalido",
    "contrasena":"123",
    "fecha_nacimiento":"2000-01-01T00:00:00Z",
    "rol":"usuario"
  }'
echo -e "\n"

echo "Prueba 4: Enviar JSON malformado/roto - Se espera: 400 Bad Request"
curl -i -X POST "$BASE_DIR/usuarios/" \
  -H "Content-Type: application/json" \
  -d '{ "nombre": "Incompleto"'
echo -e "\n"

echo "=================================================="
echo "=== 2. PRUEBAS DE LECTURA (GET /usuarios/) ==="
echo "=================================================="

echo "Prueba 5: Listar todos los usuarios - Se espera: 200 OK"
curl -i -X GET "$BASE_DIR/usuarios/"
echo -e "\n"

echo "Prueba 6: Obtener un usuario existente por ID (ID 1) - Se espera: 200 OK"
curl -i -X GET "$BASE_DIR/usuarios/1"
echo -e "\n"

echo "Prueba 7: Consultar con ID no numérico (ID abc) - Se espera: 400 Bad Request"
curl -i -X GET "$BASE_DIR/usuarios/abc"
echo -e "\n"

echo "Prueba 8: Consultar usuario inexistente (ID 999) - Se espera: 404 Not Found"
curl -i -X GET "$BASE_DIR/usuarios/999"
echo -e "\n"

echo "=================================================="
echo "=== 3. PRUEBAS DE ACTUALIZACIÓN (PUT /usuarios/) ==="
echo "=================================================="

echo "Prueba 9: Actualizar usuario existente (ID 1) - Se espera: 200 OK"
curl -i -X PUT "$BASE_DIR/usuarios/1" \
  -H "Content-Type: application/json" \
  -d '{
    "nombre":"Agustin Loray",
    "email":"agustin.actualizado@email.com",
    "contrasena":"Contrasena123$",
    "fecha_nacimiento":"2000-01-01T00:00:00Z",
    "rol":"administrador"
  }'
echo -e "\n"

echo "Prueba 10: Actualizar usuario inexistente (ID 999) - Se espera: 404 Not Found"
curl -i -X PUT "$BASE_DIR/usuarios/999" \
  -H "Content-Type: application/json" \
  -d '{
    "nombre":"Usuario Inexistente",
    "email":"inexistente@email.com",
    "contrasena":"Contrasena123$",
    "fecha_nacimiento":"2000-01-01T00:00:00Z",
    "rol":"usuario"
  }'
echo -e "\n"

echo "=================================================="
echo "=== 4. PRUEBAS DE ELIMINACIÓN (DELETE /usuarios/) ==="
echo "=================================================="

echo "Prueba 11: Eliminar usuario existente (ID 1) - Se espera: 200 OK o 204 No Content"
curl -i -X DELETE "$BASE_DIR/usuarios/1"
echo -e "\n"

echo "Prueba 12: Confirmar eliminación (GET al ID 1 borrado) - Se espera: 404 Not Found"
curl -i -X GET "$BASE_DIR/usuarios/1"
echo -e "\n"

echo "Prueba 13: Eliminar usuario inexistente (ID 999) - Se espera: 404 Not Found"
curl -i -X DELETE "$BASE_DIR/usuarios/999"
echo -e "\n"