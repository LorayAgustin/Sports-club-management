#!/bin/bash
set -e

echo "=== 1. Tareas previas: Limpiando contenedores y volúmenes viejos ===" 
docker compose down -v --remove-orphans 

echo "=== 2. Levantando contenedor de Base de Datos (PostgreSQL) ===" 
docker compose up -d db 

echo "=== 3. Esperando que PostgreSQL esté listo y acepte conexiones ===" 
until docker exec postgres-tpEspecial pg\_isready -U admin\_tpEspecial -d tpEspecial\_db; do 
  echo "Esperando a la base de datos..." 
  sleep 2 
done 

echo "=== 4. Generando código de sqlc (si aplica) ===" 
if command -v sqlc >/dev/null 2>&1; then 
  sqlc generate 
fi

BASE_DIR="http://localhost:8080/"

echo "Prueba 1: Crear un usuario VALIDO (POST /usuarios) - Se espera: respuesta 201"
curl -i -X POST "$BASE_DIR/usuarios/" \
    -H "Content-Type: application/json" \
    -d '{
        "nombre":"Agustin Loray",
        "email":"agustin@mail.com",
        "contrasena":"Contrasena123",
        "fechaNacimiento":"2000-01-01",
        "rol":"USUARIO"
    }'

echo -e
echo "Prueba 2: Crear un usuario INVALIDO (POST /usuarios) - Se espera: Error 400"
echo "Contaseña invalida"


echo "Prueba 3: Listar todos los usuarios (GET /usuarios) - Se espera: respuesta 200"
curl -i -X GET "$BASE_DIR/usuarios"

echo "=== 6. Tareas posteriores: Limpiando contenedores y volúmenes ===" 
docker compose down -v 

echo "=== ¡Pruebas completadas exitosamente! ==="