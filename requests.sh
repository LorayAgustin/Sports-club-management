#!/bin/bash
BASE_DIR="http://localhost:8080"

echo "Prueba 1: Crear un usuario VALIDO (POST /usuarios) - Se espera: respuesta 201"
curl -i -X POST "$BASE_DIR/usuarios/" \
    -H "Content-Type: application/json" \
    -d '{
        "nombre":"Agustin Loray",
        "email":"agustin@mail.com",
        "contrasena":"Contrasena123",
        "fechaNacimiento":"2000-01-01",
        "rol":"SOCIO"
    }'

echo -e
echo "Prueba 2: Crear un usuario NO VALIDO (POST /usuarios) - Se espera: Error 400"
echo "Contaseña invalida"


echo "Prueba 3: Listar todos los usuarios (GET /usuarios) - Se espera: respuesta 200"
curl -i -X GET "$BASE_DIR/usuarios"