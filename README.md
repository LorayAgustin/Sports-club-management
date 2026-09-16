# Sport-Club-Management

Nuestra pagina web permite administrar la información del club, los deportes que se ofrecen y los socios que pertenecen al mismo.

Este repositorio contiene la resolucion de la segunda entrega del **TP2: Persistiendo el Dominio**

---

## Requisitos:

* PostgreSQL: 15-alpine
* Docker Compose: 3.8
* Driver de go: github.com/jackc/pgx/v5
* No se necesita tener Go ni sqlc instalados localmente gracias al Dockerfile

---

## Instrucciones de ejecucion:
1. **Clonar el repositorio y acceder al proyecto.**
    * Luego posicionate con cd en la carpeta Sports-club-management.
```bash
git clone https://github.com/LorayAgustin/Sports-club-management.git
cd Sports-club-management
```
    
2. **Cambiar a la rama de tp2.**
    * En caso de que se cree una rama main principal por defecto (ya que no subimos la main), realizar el comando para cambiar a la rama tp2.
```bash
git checkout tp2
```

3. **Correr los tests.**
```bash
./test.sh
```
---

## Explicacion ejecucion test.sh
1. Limpieza de **contenedores y volumnes viejos**
2. Levantamiento del contenedor de la **Base de Datos**
3. Espera activa a que **PostgreSQL** esté listo
4. Generación de código con sqlc **(sqlc generate)**
5. Construcción de imagen y ejecución de pruebas en Docker
    1. Aislamiento de Go (Dockerfile): asegura que el usuario no necesite tener Go instalado en su sistema operativo local
    2. Ejecución de tests
    3. Limpieza del contenedor al terminar las pruebas
6. Baja del contenedor (limpieza de contenedores y volumenes)