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

echo "=== 5. Ejecutando Tests de Go con paquete 'testing' ===" 
go test -v ./... 

echo "=== 6. Tareas posteriores: Limpiando contenedores y volúmenes ===" 
docker compose down -v 

echo "=== ¡Pruebas completadas exitosamente! ==="