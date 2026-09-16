# Sport-Club-Management

Nuestra pagina web permite administrar la información del club, los deportes que se ofrecen y los socios que pertenecen al mismo.

Este repositorio contiene la resolucion de la segunda entrega del **TP2: Persistiendo el Dominio**

---

## Requisitos:
    * Tener instalado Go (version 1.20 o superior)
    * sqlc: v2
    * PostgreSQL: 15-alpine
    * Docker Compose: 3.8
    * Driver de go: github.com/lib/pq

Instruccion de ejecucion:
    1. Clonar el repositorio y acceder al proyecto.
        *Luego posicionate con cd en la carpeta aplicacion_web.
```bash
git clone https://github.com/LorayAgustin/Sports-club-management.git
cd Sports-club-management
git checkout tp2
./test.sh
``/test.sh
```
    
    2. Cambiar a la rama de tp2
        *En caso de que se cree una rama main principal por defecto (ya que no subimos la main), realizar el comando 
        `git checkout tp2` para cambiar a la misma.
```bash
git checkout tp2
./test.sh
```

    3. Correr los tests.
```bash
./test.sh
```