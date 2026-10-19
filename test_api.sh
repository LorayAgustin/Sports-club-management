#!/bin/bash
set -e

echo "=== 1. Tareas previas: Limpiando contenedores y volúmenes viejos ===" 
docker compose down -v --remove-orphans 

echo "=== 2. Generando código de sqlc (si aplica) ===" 
if command -v sqlc >/dev/null 2>&1; then 
  sqlc generate 
  echo "Codigo de sqlc generado correctamente."
fi

echo "=== 3. Levantando contenedor de Base de Datos (PostgreSQL) ===" 
docker compose up -d db 

echo "=== 4. Esperando que PostgreSQL esté listo y acepte conexiones ===" 
until docker exec postgres-tpEspecial pg_isready -U admin_tpEspecial -d tpEspecial_db; do 
  echo "Esperando a la base de datos..." 
  sleep 2 
done 

sleep 2

echo "=== 5. Ejecutando Pruebas HTTP de la API REST ==="
#docker compose run --build --rm app go test -v ./...

echo "=== 6. Tareas posteriores: Limpiando contenedores y volúmenes ===" 
docker compose down -v 

echo "=== ¡Pruebas completadas exitosamente! ==="