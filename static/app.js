// URL de la API para obtener el listado de deportes
const API_URL = "http://localhost:8080/deportes"

// Obtener referencias a los elementos del DOM

const listaDeportes = document.getElementById("lista-deportes");
const formDeporte = document.getElementById("form-deporte");
const inputNombre = document.getElementById("nombre-deporte");

// Escuchar la carga del documento HTML para iniciar la peticion
document.addEventListener("DOMContentLoaded", () => {
    obtenerDeportes();
});

// Función para obtener el listado de deportes desde la API
function obtenerDeportes() {
    fetch(API_URL)
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

// Función para renderizar la lista de deportes en el DOM
function renderizarLista(deportes) {
    listaDeportes.innerHTML = ""; // Limpiar la lista anterior

    if (!deportes || deportes.length === 0) {
        listaDeportes.innerHTML = "<li>No hay deportes disponibles</li>";
        return;
    }

    deportes.forEach(deporte => {
        const li = document.createElement("li");
        li.style.marginBottom = "10px";

        li.innerHTML = `
            <span><strong>ID ${deporte.id}:</strong> ${deporte.nombre}</span>
            <button class="btn-eliminar" data-id="${deporte.id}" style="margin-left: 10px; color: red;"> 
                Eliminar 
            </button>
        `;
        listaDeportes.appendChild(li);
    });
}

// Evento POST al enviar al formulario
formDeporte.addEventListener("submit", (event) => {

    event.preventDefault();

    const nombreValor = inputNombre.value.trim();
    if (nombreValor === "") {
        alert("Ingrese un nombre válido.");
        return
    }

    fetch(API_URL, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
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
    fetch(`${API_URL}/${id}`, {
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
