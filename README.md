# Sport-Club-Management

Nuestra pagina web permite administrar la información del club, los deportes que se ofrecen y los socios que pertenecen al mismo.

Este repositorio contiene la resolucion de la tercera entrega del **TP3: Exponiendo la Logica como API**

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
├── handlers/              # Handlers HTTP de las distintas entidades
│   └── usuarios.go
├── logic/                 # Logica de las distintas entidades
│   └── usuarios.go
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
├── test_api.sh            # Script de ejecución automatizada de pruebas
├── request.sh             # Script contenedor de pruebas a realizar
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

## Documentación de la API REST

A continuación se detallan los endpoints expuestos por el servicio para la gestión de **[Usuarios]**:

| Método     | Endpoint         | Descripción                                 | Body (JSON)                                                                                                | Respuestas HTTP                                  |
| ---------- | ---------------- | ------------------------------------------- | ---------------------------------------------------------------------------------------------------------- | ------------------------------------------------ |
| **POST**   | `/usuarios`      | Crea un nuevo usuario                       | `{ "nombre": "...", "email": "...", "contrasena": "...", "fecha_nacimiento": "YYYY-MM-DD", "rol": "..." }` | `201 Created`<br />`400 Bad Request`               |
| **GET**    | `/usuarios`      | Obtiene el listado completo de usuarios     | *Ninguno*                                                                                                  | `200 OK`<br />`500 Internal Error`                 |
| **GET**    | `/usuarios/{id}` | Busca un usuario por su ID                  | *Ninguno*                                                                                                  | `200 OK`<br />`404 Not Found`                      |
| **PUT**    | `/usuarios/{id}` | Actualiza los datos de un usuario existente | `{ "nombre": "...", "email": "...", "contrasena": "...", "fecha_nacimiento": "YYYY-MM-DD", "rol": "..." }` | `200 OK`<br />`400 Bad Request`<br />`404 Not Found` |
| **DELETE** | `/usuarios/{id}` | Elimina un usuario por su ID                | *Ninguno*                                                                                                  | `204 No Content`<br />`404 Not Found`              |

---

### Capa de Lógica de Negocio (`/logic`)
* **Responsabilidad:** Contiene las funciones puras y validaciones del dominio (ej. `ValidateUsuario`, comprobación de formatos de fecha, reglas de negocio).
* **Desacoplamiento:** Es totalmente independiente de HTTP y de la base de datos; no utiliza `http.Request` ni `http.ResponseWriter`, lo que permite probarla mediante tests unitarios de Go sin necesidad de levantar un servidor.

---

### Capa de Handlers HTTP (`/handlers`)
* **Responsabilidad:** Actúa como controlador de la API REST.
* **Flujo:**
  1. Recibe la petición HTTP y decodifica el cuerpo en formato JSON (`r.Body`).
  2. Invoca las validaciones de la capa `logic`.
  4. Llama a las consultas persistentes autogeneradas por `sqlc`.
  5. Construye la respuesta HTTP estableciendo el encabezado `Content-Type: application/json` y retornando los códigos de estado apropiados (`201 Created`, `200 OK`, `400 Bad Request`, `404 Not Found`, `500 Internal Error`).

---

## Instrucciones de ejecucion:
1. **Clonar el repositorio y acceder al proyecto.**
    * Luego posicionate con cd en la carpeta Sports-club-management.
```bash
git clone https://github.com/LorayAgustin/Sports-club-management.git
cd Sports-club-management
```
    
2. **Cambiar a la rama de tp3.**
    * En caso de que se cree una rama main principal por defecto (ya que no subimos la main), realizar el comando para cambiar a la rama tp3.
```bash
git checkout tp3
```

3. **Ejecucion de pruebas.**
```bash
chmod +x test_api.sh
./test_api.sh
```
--- 

## Explicacion ejecucion test_api.sh
1. Limpieza de contenedores y volumnes viejos
2. Generación de código con sqlc **(sqlc generate)**
3. Levantamiento del contenedor de la Base de Datos
4. Espera activa a que PostgreSQL esté listo
5. Compilación y despliegue del servicio web Go
6. Espera activa a la API HTTP
7. Ejecucion de las Pruebas HTTP de la API REST (`requests.sh`)
8. Limpieza posterior del entorno