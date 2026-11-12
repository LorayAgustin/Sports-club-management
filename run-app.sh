#!/bin/bash
set -e

# Limpieza automática al presionar Ctrl+C o salir del script
cleanup() {
  echo ""
  echo "=== 3. Tareas posteriores: Limpiando contenedores y volúmenes ==="
  docker compose down -v --remove-orphans
  echo "=== Limpieza finalizada. ==="
}
trap cleanup EXIT

echo "=== 1. Tareas previas: Limpiando contenedores y volúmenes viejos ==="
docker compose down -v --remove-orphans

echo "=== Levantando el contenedor de Base de Datos (PostgreSQL) ==="
docker compose up -d db

echo "=== Esperando que PostgreSQL esté listo y acepte conexiones ==="
until docker exec postgres-tpEspecial pg_isready -U admin_tpEspecial -d tpEspecial_db > /dev/null 2>&1; do
  echo "Esperando a la base de datos..."
  sleep 2
done

echo "=== 4b. Levantando y compilando la Aplicación Go ==="
docker compose up -d --build app

echo "=== 4c. Esperando a que la API de Go responda ==="
until curl -s http://localhost:8080/usuarios > /dev/null; do
  echo "Esperando a que la API de Go responda..."
  sleep 2
done

echo ""
echo "=== 2. ¡Aplicación corriendo en http://localhost:8080! ==="
echo "Presione Ctrl+C cuando termine de probar para detener y limpiar todo."

# Mantiene el script abierto mostrando los logs de la app hasta que el profesor presione Ctrl+C
docker compose logs -f app
