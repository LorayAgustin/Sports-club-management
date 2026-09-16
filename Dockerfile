# 1. Imagen base oficial de Go sobre Alpine Linux (liviana y eficiente) 
FROM golang:1.25-alpine 

# 2. Directorio de trabajo interno dentro del contenedor 
WORKDIR /app 

# 3. Copiamos los archivos de módulos y descargamos las dependencias 
COPY go.mod go.sum ./ 
RUN go mod download 

# 4. Copiamos el resto del código fuente (incluyendo la carpeta db/sqlc generada) 
COPY . . 

# 5. Compilamos el binario a partir del main.go 
RUN go build -o /app/servidor main.go 

# 6. Comando para iniciar el ejecutable 
CMD ["/app/servidor"]