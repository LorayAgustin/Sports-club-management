# 1. Imagen base oficial de Go sobre Alpine Linux (liviana y eficiente) 
FROM golang:1.26-alpine 

# 2. Directorio de trabajo interno dentro del contenedor 
WORKDIR /app 

# 3. Instalar los ejecutables CLI de sqlc y templ
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install github.com/a-h/templ/cmd/templ@latest

# 4. Copiamos los archivos de módulos y descargamos las dependencias 
COPY go.mod go.sum ./ 
RUN go mod download

# 5. Copiamos el resto del código fuente (incluyendo la carpeta db/sqlc generada) 
COPY . . 

# 6. Generar el código de Go de la base de datos (sqlc) y de las vistas (templ)
RUN sqlc generate 
RUN templ generate

# 7. Compilamos el binario a partir del main.go 
RUN go build -o /app/servidor main.go 

# 8. Comando para iniciar el ejecutable 
CMD ["/app/servidor"]