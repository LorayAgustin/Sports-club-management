# Sport-Club-Management

Aplicación web para administrar la información de un club deportivo: usuarios, socios, deportes, predios, categorías, inscripciones y prácticas.

Este repositorio contiene la resolución de la entrega final **TP5: Aplicación web con vistas renderizadas desde el servidor (Templ)**.

> **Restricción de la entrega:** el proyecto **no contiene archivos JavaScript (`.js`)**. Toda la interfaz se genera en el servidor con plantillas **Templ** y los estilos se resuelven con hojas CSS estáticas. Los archivos `*.js` están excluidos en el `.gitignore`.

---

## Stack tecnológico

| Componente | Tecnología | Rol en el proyecto |
| :--- | :--- | :--- |
| **Lenguaje** | Go | Servidor HTTP, lógica de negocio y controladores |
| **Base de datos** | PostgreSQL 16 (Alpine) | Persistencia de datos |
| **Acceso a datos** | `sqlc` + `pgx/v5` | Generación de código Go tipado a partir de consultas SQL |
| **Vistas** | `templ` | Plantillas HTML tipadas, compiladas a código Go |
| **Estilos** | CSS | Hojas de estilo estáticas en `/static` |
| **Infraestructura** | Docker / Docker Compose | Orquestación de la aplicación y la base de datos |
| **Automatización** | Bash (`run-app.sh`) | Build, ejecución y limpieza automatizada del entorno |

---

## Arbol de directorios

```
Sports-club-management/
├── db/
│   ├── queries/               # Consultas SQL anotadas para sqlc (CRUD)
│   │   ├── categorias.sql
│   │   ├── deportes.sql
│   │   ├── inscripciones.sql
│   │   ├── practicas.sql
│   │   ├── predios.sql
│   │   ├── socios.sql
│   │   └── usuarios.sql
│   ├── schema/                # Esquema DDL de la base de datos
│   │   └── schema.sql
│   └── sqlc/                  # Código Go generado por `sqlc generate` (no se versiona)
│       ├── db.go
│       ├── models.go
│       ├── categorias.sql.go
│       ├── deportes.sql.go
│       ├── inscripciones.sql.go
│       ├── practicas.sql.go
│       ├── predios.sql.go
│       ├── socios.sql.go
│       └── usuarios.sql.go
├── handlers/                  # Controladores web (HTTP + render de vistas)
│   ├── categorias.go
│   ├── deportes.go
│   ├── home.go
│   ├── inscripciones.go
│   ├── practicas.go
│   ├── predios.go
│   ├── socios.go
│   └── usuarios.go
├── logic/                     # Validaciones y reglas de negocio de cada entidad
│   ├── categorias.go
│   ├── deportes.go
│   ├── inscripciones.go
│   ├── practicas.go
│   ├── predios.go
│   ├── socios.go
│   └── usuarios.go
├── views/                     # Vistas y componentes Templ
│   ├── layout.templ           # Layout base (<head>, hoja de estilos y <body>)
│   ├── index.templ            # Menú principal (acceso a los 7 módulos)
│   ├── usuarios.templ         # Página de usuarios
│   ├── usuarios_form.templ    # Componente: formulario de alta
│   ├── usuarios_list.templ    # Componente: listado
│   ├── socios.templ
│   ├── socios_form.templ
│   ├── socios_list.templ
│   ├── deportes.templ
│   ├── deporte_form.templ
│   ├── deportes_list.templ
│   ├── predios.templ
│   ├── predios_form.templ
│   ├── predios_list.templ
│   ├── categorias.templ
│   ├── categorias_form.templ
│   ├── categorias_list.templ
│   ├── inscripciones.templ
│   ├── inscripciones_form.templ
│   ├── inscripciones_list.templ
│   ├── practicas.templ
│   ├── practicas_form.templ
│   ├── practicas_list.templ
│   └── *_templ.go             # Código Go generado por `templ generate` (uno por cada .templ, no se versiona)
├── static/                    # Archivos estáticos servidos en /static/ (solo CSS)
│   ├── styles.css
│   ├── styles-dark.css        # Tema utilizado por el layout base
│   ├── styles-minimal.css
│   └── styles-sport.css
├── Dockerfile                 # Compila la app: sqlc generate, templ generate y go build
├── docker-compose.yml         # Orquestador multi-contenedor (App Go + PostgreSQL)
├── run-app.sh                 # Script de ejecución automatizada (build, run y limpieza)
├── go.mod                     # Módulo de Go y definición de dependencias
├── go.sum                     # Sumas de verificación de dependencias
├── main.go                    # Punto de entrada: conexión a la BD, rutas y servidor HTTP
├── sqlc.yaml                  # Configuración y mapeos para sqlc
├── requests.sh                # Heredado de la entrega TP3 (pruebas de la API REST con JSON, no aplica a TP5)
├── .env.example               # Plantilla de variables de entorno
└── README.md                  # Documentación del proyecto
```

---

## Arquitectura

La aplicación renderiza el HTML **en el servidor**: el navegador envía formularios HTML estándar y el servidor responde con páginas generadas con **Templ**. Ya no se expone una API REST con JSON.

```
Navegador ──HTTP (GET / POST de formularios)──▶ Handlers ──▶ sqlc ──▶ PostgreSQL
    ▲                                              │  │
    │                                              │  └──▶ Logic (validaciones)
    └──────────────── HTML (Templ) ◀───────────────┘
```

### Capa de Lógica de Negocio (`/logic`)
* **Responsabilidad:** Contiene las validaciones y reglas del dominio de cada entidad (ej. `ValidateUsuario`, `ValidateSocio`, `ValidatePredio`: formatos de email, fechas, roles, estados, etc.).
* **Desacoplamiento:** No utiliza `http.Request` ni `http.ResponseWriter`, por lo que puede probarse con tests unitarios de Go sin levantar un servidor. Cuando una validación necesita consultar la base de datos (ej. unicidad del nombre de usuario) lo hace a través de una interfaz (`UserDB`), no de una dependencia concreta.

### Capa de Handlers HTTP (`/handlers`)
* **Responsabilidad:** Actúan como **controladores web**: reciben la petición, consultan o modifican la base de datos y devuelven HTML renderizado con Templ. Cada handler recibe en su constructor el cliente `*sqlc.Queries`.
* **Flujo de un `GET`** (ej. `/usuarios`):
  1. Valida que la ruta sea la esperada.
  2. Consulta los datos necesarios con las funciones generadas por `sqlc` (`ListUsuarios`, y en las entidades con relaciones también `ListSocios`, `ListCategorias`, `ListPredios`, `ListDeportes`).
  3. Renderiza la página completa con Templ directamente sobre la respuesta (`.Render(r.Context(), w)`).
* **Flujo de un `POST`** (ej. alta de un usuario):
  1. Lee el formulario (`r.ParseForm()` / `r.FormValue()`) y convierte los tipos (fechas, enteros, horas).
  2. Valida con la capa `logic`.
  3. Persiste con la consulta `Create...` generada por `sqlc`.
  4. Responde con una redirección **`303 See Other`** a la página de la entidad (patrón Post/Redirect/Get), que evita reenvíos del formulario al refrescar el navegador.
* **Códigos de estado:** `400 Bad Request` (formulario inválido o error de validación), `405 Method Not Allowed` (método distinto de `GET`/`POST`), `404 Not Found` (ruta inexistente) y `500 Internal Server Error` (fallo de base de datos o de renderizado). Los errores se devuelven como texto plano.

### Capa de Vistas (`/views`: Templ)
* **Templ:** cada archivo `.templ` define componentes HTML tipados en Go. `templ generate` los compila a archivos `*_templ.go`, por lo que los errores de plantilla se detectan en tiempo de compilación.
* **Composición:** `layout.templ` define la estructura base (`<head>`, hoja de estilos `styles-dark.css` y `<body>`), y cada página lo envuelve con su contenido. Los formularios (`*_form.templ`) y los listados (`*_list.templ`) son componentes independientes que cada página compone.
* **Escape automático:** Templ escapa los valores interpolados, lo que previene inyección de HTML/XSS.
* **Sin JavaScript:** no hay archivos `.js` en el proyecto; el directorio `static/` contiene únicamente hojas de estilo.

---

## Rutas de la aplicación

Todas las rutas de entidades aceptan `GET` (formulario de alta + listado) y `POST` (alta).

| Método | Ruta | Descripción | Respuesta |
| :--- | :--- | :--- | :--- |
| **GET** | `/` | Menú principal de navegación | Página HTML |
| **GET / POST** | `/usuarios` | Gestión de usuarios | `GET`: página HTML · `POST`: `303` a `/usuarios` |
| **GET / POST** | `/socios` | Gestión de socios | `GET`: página HTML · `POST`: `303` a `/socios` |
| **GET / POST** | `/deportes` | Gestión de deportes | `GET`: página HTML · `POST`: `303` a `/deportes` |
| **GET / POST** | `/predios` | Gestión de predios / sedes | `GET`: página HTML · `POST`: `303` a `/predios` |
| **GET / POST** | `/categorias` | Gestión de categorías | `GET`: página HTML · `POST`: `303` a `/categorias` |
| **GET / POST** | `/inscripciones` | Gestión de inscripciones | `GET`: página HTML · `POST`: `303` a `/inscripciones` |
| **GET / POST** | `/practicas` | Gestión de prácticas | `GET`: página HTML · `POST`: `303` a `/practicas` |
| **GET** | `/static/*` | Archivos estáticos | Archivo CSS |

Si el `POST` falla por validación o por un error de la base de datos, la respuesta es un `400` o `500` con el mensaje en texto plano.

### Campos del formulario de alta (`POST`, `application/x-www-form-urlencoded`)

| Ruta | Campos |
| :--- | :--- |
| `/usuarios` | `nombre`, `email`, `contrasena`, `fechaNacimiento` (`YYYY-MM-DD`), `rol` |
| `/socios` | `id_usuario`, `dni`, `fechaAlta` (`YYYY-MM-DD`, opcional: por defecto la fecha actual) |
| `/deportes` | `nombre` |
| `/predios` | `nombre`, `direccion` |
| `/categorias` | `nombre`, `id_deporte` |
| `/inscripciones` | `id_socio`, `id_categoria`, `fecha_inscripcion` (`YYYY-MM-DD`), `estado` |
| `/practicas` | `id_categoria`, `id_predio`, `dia_semana`, `hora_inicio` (`HH:MM`), `hora_fin` (`HH:MM`) |

**Reglas de validación de usuarios:** nombre de al menos 3 caracteres y único; email con formato válido; contraseña de al menos 8 caracteres, sin espacios, con un número, una mayúscula y un carácter especial; fecha de nacimiento no posterior a hoy; rol `administrador`, `usuario`, `socio` o `empleado`.

---

## Datos de Conexión y Persistencia (PostgreSQL)

La base de datos se ejecuta en un contenedor Docker. Los valores se definen en el archivo `.env` (ver `.env.example`):

| Parámetro | Valor |
| :--- | :--- |
| **Motor** | PostgreSQL 16 (Alpine) |
| **Contenedor** | `postgres-tpEspecial` |
| **Host (Localhost)** | `localhost` |
| **Host (Red Docker)** | `db` / `postgres-tpEspecial` |
| **Puerto** | `5432` |
| **Base de Datos** | Variable `POSTGRES_DB` |
| **Usuario** | Variable `POSTGRES_USER` |
| **Contraseña** | Variable `POSTGRES_PASSWORD` |
| **Cadena de conexión de la app** | Variable `DB_URL` |

El esquema (`db/schema/schema.sql`) se carga automáticamente al inicializar el contenedor. Los datos se guardan en el volumen `postgres_data`.

> **Importante:** `run-app.sh` elimina los volúmenes al finalizar (`docker compose down -v`), por lo que los datos cargados durante una ejecución **no se conservan** entre ejecuciones del script.

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

Las consultas SQL parametrizadas siguen el ciclo **CRUD** (Create, Read, Update, Delete). Mediante las anotaciones de **sqlc** se indica el tipo de retorno esperado, y la herramienta genera código Go seguro y tipado:

* **CreateTabla** **(** **CREATE** **):** Inserta un nuevo registro (`INSERT INTO ...`) con la anotación `:one`, retornando la entidad recién creada con su ID generado.
* **GetTabla** **(** **READ** **/** **GET** **):** Obtiene un registro a partir de su clave primaria con la anotación `:one`.
* **ListTabla** **(** **READ** **/** **LIST** **):** Recupera la lista completa de registros con la anotación `:many`.
* **UpdateTabla** **(** **UPDATE** **):** Modifica los campos de un registro existente con la anotación `:exec`.
* **DeleteTabla** **(** **DELETE** **):** Elimina un registro por su ID con la anotación `:exec`.

En esta entrega, los handlers utilizan las operaciones de **alta** (`Create...`) y **listado** (`List...`), además de `ExistsUserByName` para validar la unicidad del nombre de usuario.

---

## Requisitos previos

* **Git**
* **Docker** y **Docker Compose** (v2, comando `docker compose`), con el servicio de Docker en ejecución
* **Bash** (Linux, macOS, o WSL / Git Bash en Windows) para ejecutar `run-app.sh`
* **curl** (lo usa `run-app.sh` para verificar que el servidor responde)
* Puertos **8080** (aplicación) y **5432** (PostgreSQL) libres en la máquina local
* Conexión a internet en el primer build (descarga de imágenes y dependencias)
* **No** hace falta tener instalados Go, `sqlc` ni `templ` localmente: el `Dockerfile` los instala y ejecuta dentro del contenedor

---

## Instrucciones de ejecucion

1. **Clonar el repositorio y acceder al proyecto.**
```bash
git clone https://github.com/LorayAgustin/Sports-club-management.git
cd Sports-club-management
```

2. **Cambiar a la rama `tp5`.**
    * Si el repositorio queda posicionado en otra rama por defecto, cambiar a `tp5`:
```bash
git checkout tp5
```

3. **Crear el archivo `.env`.**
    * Copiar la plantilla para generar el archivo `.env` local requerido por Docker Compose:
```bash
cp .env.example .env
```

4. **Ejecutar la aplicación con el script automatizado.**
    * Dar permisos de ejecución al script (solo hace falta la primera vez) y ejecutarlo:
```bash
chmod +x run-app.sh
./run-app.sh
```

5. **Esperar el mensaje de confirmación.** Cuando el script imprima `¡Aplicación corriendo en http://localhost:8080!`, abrir la aplicación en el navegador: [http://localhost:8080](http://localhost:8080)

6. **Detener la aplicación.** Presionar `Ctrl+C` en la terminal donde corre el script. El entorno se limpia automáticamente (ver sección siguiente).

---

## Explicacion ejecucion run-app.sh

El script está escrito con `set -e` (se detiene ante cualquier error) y registra al inicio la función de limpieza con `trap cleanup EXIT`, por lo que la limpieza se ejecuta siempre al salir.

### 1. Tareas previas
1. **Limpieza de contenedores y volúmenes previos:** `docker compose down -v --remove-orphans`. Garantiza partir de un entorno limpio y de una base de datos vacía.
2. **Levantamiento de PostgreSQL:** `docker compose up -d db`.
3. **Espera activa de la base de datos:** el script ejecuta `pg_isready` dentro del contenedor `postgres-tpEspecial`, reintentando cada 2 segundos hasta que PostgreSQL acepta conexiones.
4. **Compilación y levantamiento de la aplicación Go:** `docker compose up -d --build app`. Durante el build, el `Dockerfile` ejecuta internamente:
    * `sqlc generate`: genera el código de acceso a datos a partir de `db/queries` y `db/schema`.
    * `templ generate`: compila los archivos `.templ` a código Go (`*_templ.go`).
    * `go build`: compila el binario del servidor (`/app/servidor`).

### 2. Ejecución
1. **Verificación de disponibilidad:** el script consulta con `curl` la URL `http://localhost:8080/usuarios`, reintentando cada 2 segundos hasta que el servidor web responde.
2. **Aviso al usuario:** imprime `¡Aplicación corriendo en http://localhost:8080!` junto con la indicación de presionar `Ctrl+C` para finalizar.
3. **Logs en tiempo real:** ejecuta `docker compose logs -f app`, que mantiene el script en primer plano mostrando los logs de la aplicación hasta que se lo interrumpa.

### 3. Tareas posteriores
Al presionar `Ctrl+C` (o ante cualquier salida del script, incluso por un error en los pasos anteriores), el `trap` ejecuta automáticamente la función `cleanup`:
* `docker compose down -v --remove-orphans`: detiene y elimina los contenedores, la red y los volúmenes (incluida la base de datos), sin dejar recursos huérfanos.

### Alternativa manual (sin el script)
```bash
docker compose up -d --build          # Compila y levanta la base de datos y la aplicación
docker compose logs -f app            # (Opcional) Logs de la aplicación
docker compose down -v --remove-orphans   # Detiene y limpia el entorno
```

---