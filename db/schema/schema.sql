CREATE TABLE usuarios (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(250) NOT NULL,
    email VARCHAR(250) NOT NULL UNIQUE,
    contrasena VARCHAR(250) NOT NULL,
    fechaNacimiento DATE NOT NULL,
    rol VARCHAR(50) NOT NULL
);

CREATE TABLE socios (
    id_Socio SERIAL PRIMARY KEY,
    id_Usuario INT NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    dni VARCHAR(10) NOT NULL UNIQUE,
    fechaAlta DATE NOT NULL,
    fechaBaja DATE
);

CREATE TABLE deportes (
    id_Deporte SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE predios (
    id_Predio SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE,
    direccion VARCHAR(250) NOT NULL
);

CREATE TABLE categorias (
    id_Categoria SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE,
    id_Deporte INT NOT NULL REFERENCES deportes(id_Deporte) ON DELETE CASCADE
);

CREATE TABLE inscripciones (
    id_Inscripcion SERIAL PRIMARY KEY,
    id_Socio INT NOT NULL REFERENCES socios(id_Socio) ON DELETE CASCADE,
    id_Categoria INT NOT NULL REFERENCES categorias(id_Categoria) ON DELETE CASCADE,
    fechaInscripcion DATE NOT NULL,
    estado VARCHAR(20) NOT NULL,
    fechaBaja DATE
);

CREATE TABLE practicas (
    id_Practica SERIAL PRIMARY KEY,
    id_Categoria INT NOT NULL REFERENCES categorias(id_Categoria) ON DELETE CASCADE,
    id_Predio INT NOT NULL REFERENCES predios(id_Predio) ON DELETE CASCADE,
    diaSemana VARCHAR(20) NOT NULL,
    horaInicio TIMESTAMP NOT NULL,
    horaFin TIMESTAMP NOT NULL
);


