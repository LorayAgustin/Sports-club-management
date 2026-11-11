// URLs de la API REST en Go
const API_DEPORTES = "http://localhost:8080/deportes"
const API_PREDIOS = "http://localhost:8080/predios"

// Referencias DOM - Deportes
const listaDeportes = document.getElementById("lista-deportes");
const formDeporte = document.getElementById("form-deporte");
const inputNombre = document.getElementById("nombre-deporte");

// Referencias DOM - Predios
const listaPredios = document.getElementById("lista-predios");
const formPredio = document.getElementById("form-predio");
const inputNombrePredio = document.getElementById("nombre-predio");
const inputDireccionPredio = document.getElementById("direccion-predio");

// CARGA INICIAL
document.addEventListener("DOMContentLoaded", () => {
    obtenerDeportes();
    obtenerPredios();
});

// MODULO 1: DEPORTES

function obtenerDeportes() {
    fetch(API_DEPORTES)
        .then(response => {
            if (!response.ok) {
                throw new Error(`Error al obtener el listado de deportes: ${response.status}`);
            }
            return response.json();
        })
        .then(deportes => {
            renderizarLista(deportes); // Manda el arreglo a renderizar en el DOM
        })
        .catch(error => {
            console.error("Error en la peticion GET:", error);
            listaDeportes.innerHTML = "<li style='color: red;'>Error al cargar los deportes</li>";
        });
}

function renderizarLista(deportes) {
    listaDeportes.innerHTML = ""; // Limpiar la lista anterior
    if (!deportes || deportes.length === 0) {
        listaDeportes.innerHTML = "<li>No hay deportes disponibles</li>";
        return;
    }

    deportes.forEach(deporte => {
        const deporteId = deporte.id_deporte || deporte.id
        const li = document.createElement("li");
        li.style.marginBottom = "10px";
        li.innerHTML = `
            <span><strong>ID ${deporteId}:</strong> ${deporte.nombre}</span>
            <button class="btn-eliminar" data-id="${deporteId}" style="margin-left: 10px; color: white;"> 
                Eliminar 
            </button>
        `;
        listaDeportes.appendChild(li);
    });
}

// Evento POST al enviar al formulario - Deportes
formDeporte.addEventListener("submit", (event) => {
    event.preventDefault();
    const nombreValor = inputNombre.value.trim();
    if (nombreValor === "") {
        alert("Ingrese un nombre válido.");
        return
    }

    fetch(API_DEPORTES, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({ nombre: nombreValor })
    })
    .then(response => { 
        if (!response.ok) { 
            return response.text().then(msg => { throw new Error(msg || "Error al crear"); }); 
        } 
        return response.json(); 
    })
    .then(nuevoDeporte => { 
        console.log("Deporte creado:", nuevoDeporte); 
        inputNombre.value = ""; // Limpiar el input 
        obtenerDeportes(); // Volver a pedir la lista actualizada 
    })
    .catch(error => { 
        console.error("Error en POST:", error); 
        alert("Error al guardar: " + error.message); 
    });
});

// Delegación de eventos para eliminar (DELETE) en la lista <ul></ul>
listaDeportes.addEventListener("click", (event) => {
    if (event.target.classList.contains("btn-eliminar")) {
        const id = event.target.dataset.id;
        if (confirm(`¿Estás seguro de eliminar el deporte con ID ${id}?`)) {
            eliminarDeporte(id);
        }
    }
});

// Funcion para enviar la peticion DELETE
function eliminarDeporte(id){
    fetch(`${API_DEPORTES}/${id}`, {
        method: "DELETE"
    })
    .then(response => {
        if (!response.ok) {
            throw new Error("No se pudo eliminar el elemento.");
        }
        obtenerDeportes();
    })
    .catch(error => {
        console.error("Error en DELETE:", error);
        alert("Error al borrar: " + error.message);
    });
}

//MODULO 2: PREDIOS

function obtenerPredios() {
    fetch(API_PREDIOS)
        .then(res => res.ok ? res.json() : Promise.reject(res.status))
        .then(predios => renderizarPredios(predios))
        .catch(err => {
            console.error("Error GET predios:", err);
            listaPredios.innerHTML = "<li style='color: red;'>Error al cargar predios</li>";
        });
}

function renderizarPredios(predios) {
    listaPredios.innerHTML = "";

    if (!predios || predios.length === 0) {
        listaPredios.innerHTML = "<li>No hay predios registrados</li>";
        return;
    }

    predios.forEach(predio => {
        // Obtenemos id_predio devuelto por sqlc/PostgreSQL
        const id = predio.id_predio || predio.id;

        const li = document.createElement("li");
        li.innerHTML = `
            <span><strong>ID ${id}:</strong> ${predio.nombre} (${predio.direccion})</span>
            <button class="btn-eliminar btn-eliminar-predio" data-id="${id}">
                Eliminar
            </button>
        `;
        listaPredios.appendChild(li);
    });
}
formPredio.addEventListener("submit", (e) => {
    e.preventDefault();

    const nombre = inputNombrePredio.value.trim();
    const direccion = inputDireccionPredio.value.trim();

    if (!nombre || !direccion) {
        alert("Por favor completa el nombre y la dirección del predio.");
        return;
    }

    fetch(API_PREDIOS, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ nombre, direccion })
    })
        .then(res => {
            if (!res.ok) {
                return res.text().then(msg => { throw new Error(msg); });
            }
            return res.json();
        })
        .then(() => {
            inputNombrePredio.value = "";
            inputDireccionPredio.value = "";
            obtenerPredios();
        })
        .catch(err => alert("Error al guardar predio: " + err.message));
});

listaPredios.addEventListener("click", (e) => {
    if (e.target.classList.contains("btn-eliminar-predio")) {
        const id = e.target.dataset.id;

        if (confirm(`¿Eliminar predio ID ${id}?`)) {
            fetch(`${API_PREDIOS}/${id}`, {
                method: "DELETE"
            })
                .then(res => res.ok ? obtenerPredios() : Promise.reject("Error al borrar"))
                .catch(err => alert(err));
        }
    }
});