# Sport-Club-Management

Nuestra pagina web permite administrar la información del club, los deportes que se ofrecen y los socios que pertenecen al mismo.

Este repositorio contiene la resolucion de la segunda entrega del **TP2: Persistiendo el Dominio**

---

## Arbol de directorios

```
Sports-club-management/
├── db/
│   ├── queries/           # Consultas SQL anotadas para sqlc (CRUD)
│   │   ├── categorias.sql
│   │   ├── deportes.sql
│   │   ├── inscripciones.sql
│   │   ├── practicas.sql
│   │   ├── predios.sql
│   │   ├── socios.sql
│   │   └── usuarios.sql
│   ├── schema/            # Esquema DDL de la base de datos
│   │   └── schema.sql
│   └── sqlc/              # Código Go autogenerado por sqlc
│       ├── db.go
│       ├── models.go
│       ├── categorias.sql.go
│       ├── deportes.sql.go
│       ├── inscripciones.sql.go
│       ├── practicas.sql.go
│       ├── predios.sql.go
│       ├── socios.sql.go
│       └── usuarios.sql.go
├── test/                  # Suite de pruebas unitarias CRUD
│   ├── categorias_test.go
│   ├── deportes_test.go
│   ├── inscripciones_test.go
│   ├── practicas_test.go
│   ├── predios_test.go
│   ├── socios_test.go
│   └── usuarios_test.go
├── Dockerfile             # Receta para compilar la aplicación en Go (Alpine)
├── docker-compose.yml     # Orquestador multi-contenedor (App Go + PostgreSQL)
├── go.mod                 # Módulo de Go y definición de dependencias
├── go.sum                 # Sumas de verificación de dependencias
├── main.go                # Punto de entrada del servidor
├── sqlc.yaml              # Configuración y mapeos para sqlc
├── test.sh                # Script de ejecución automatizada de pruebas
└── README.md              # Documentación del proyecto

```

---

## Datos de Conexión y Persistencia (PostgreSQL)

La base de datos se ejecuta en un contenedor Docker con la siguiente configuración por defecto:

| Parámetro | Valor |
| :--- | :--- |
| **Motor** | PostgreSQL 16 (Alpine) |
| **Host (Localhost)** | `localhost` |
| **Host (Red Docker)** | `db` / `postgres-tpEspecial` |
| **Puerto** | `5432` |
| **Base de Datos** | `tpEspecial_db` |
| **Usuario** | `admin_tpEspecial` |
| **Contraseña** | `admin_contra` |

---

## Informacion sobre el esquema de la BD (db/schema/schema.sql)

| Tabla | Atributo | Descripción |
| :--- | :--- | :--- |
| **`usuarios`** | `id` | Identificador único del usuario |
| | `nombre` | Nombre completo del usuario |
| | `email` | Correo electrónico de acceso |
| | `contrasena` | Contraseña cifrada |
| | `fechaNacimiento` | Fecha de nacimiento |
| | `rol` | Rol en el sistema (ej: 'admin', 'socio') |
| **`socios`** | `id_Socio` | Identificador único del socio |
| | `id_Usuario` | Usuario asociado |
| | `dni` | Documento de identidad |
| | `fechaAlta` | Fecha de ingreso al club |
| | `fechaBaja` | Fecha de egreso/baja |
| **`deportes`** | `id_Deporte` | Identificador del deporte |
| | `nombre` | Nombre de la disciplina deportiva |
| **`predios`** | `id_Predio` | Identificador del predio/sede |
| | `nombre` | Nombre de la sede |
| | `direccion` | Dirección física de la sede |
| **`categorias`** | `id_Categoria` | Identificador de la categoría |
| | `nombre` | Nombre de la división deportiva |
| | `id_Deporte` | Deporte al que pertenece |
| **`inscripciones`** | `id_Inscripcion` | Identificador de la inscripción |
| | `id_Socio` | Socio que se inscribe |
| | `id_Categoria` | Categoría seleccionada |
| | `fechaInscripcion` | Fecha de inscripción |
| | `estado` | Estado de la inscripción ('activa', 'inactiva') |
| | `fechaBaja` | Fecha de cancelación/baja |
| **`practicas`** | `id_Practica` | Identificador de la práctica |
| | `id_Categoria` | Categoría que entrena |
| | `id_Predio` | Predio/sede de entrenamiento |
| | `diaSemana` | Día de la semana de la práctica |
| | `horaInicio` | Horario de inicio del entrenamiento |
| | `horaFin` | Horario de finalización |

---

## Consultas y Operaciones CRUD (db/queries/....sql)

Para gestionar el acceso a los datos, definimos las consultas SQL parametrizadas siguiendo el ciclo **CRUD** (Create, Read, Update, Delete). Mediante las anotaciones de **sqlc**, le indicamos a la herramienta el tipo de retorno esperado para que genere automáticamente el código Go seguro y tipado:

* **CreateTabla** **(** **CREATE** **):** Inserta un nuevo registro en la base de datos (`INSERT INTO ...`) utilizando la anotación `:one` para retornar la entidad recién creada con su ID generado.
* **GetTabla** **(** **READ** **/** **GET** **):** Obtiene un registro específico a partir de su clave primaria (`id`) utilizando la anotación `:one`.
* **ListTabla** **(** **READ** **/** **LIST** **):** Recupera la lista completa de registros almacenados utilizando la anotación `:many` para retornar una colección de datos.
* **UpdateTabla** **(** **UPDATE** **):** Modifica los campos de un registro existente identificado por su ID, utilizando la anotación `:exec` (ejecuta la sentencia sin retornar filas).
* **DeleteTabla** **(** **DELETE** **):** Elimina un registro de la base de datos por su ID, utilizando la anotación `:exec`

--- 

## Requisitos:

* PostgreSQL: 16-alpine
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
chmod +x test.sh
./test.sh
```
--- 

## Explicacion ejecucion test.sh
1. Limpieza de contenedores y volumnes viejos
2. Levantamiento del contenedor de la Base de Datos
3. Espera activa a que PostgreSQL esté listo
4. Generación de código con sqlc **(sqlc generate)**
5. Construcción de imagen y ejecución de pruebas en Docker
    1. Aislamiento de Go (Dockerfile): asegura que el usuario no necesite tener Go instalado en su sistema operativo local
    2. Ejecución de tests
    3. Limpieza del contenedor al terminar las pruebas
6. Baja del contenedor (limpieza de contenedores y volumenes)