package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Función para registrar los datos en la bitácora
func registrarEnBitacora(nombre, correo, mensaje string) error {
	// Abrimos el archivo en modo Append (añadir al final) o lo creamos si no existe
	archivo, err := os.OpenFile("bitacora.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer archivo.Close()

	// Formateamos la línea con la fecha actual y los datos recibidos
	fecha := time.Now().Format("2006-01-02 15:04:05")
	lineaLog := fmt.Sprintf("[%s] Nombre: %s | Correo: %s | Mensaje: %s\n", fecha, nombre, correo, mensaje)

	// Escribimos en el archivo
	_, err = archivo.WriteString(lineaLog)
	return err
}

// Handler para procesar el formulario de contacto
func contactoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parseamos los datos que vienen del formulario HTML
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	nombre := r.FormValue("nombre")
	correo := r.FormValue("correo")
	mensaje := r.FormValue("mensaje")

	// Guardamos los datos en la bitácora
	err = registrarEnBitacora(nombre, correo, mensaje)
	if err != nil {
		log.Println("Error al escribir en la bitácora:", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// Respuesta al usuario en el navegador
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h2>¡Solicitud recibida con éxito! Hemos registrado tus datos en la bitácora.</h2><a href='/'>Volver al inicio</a>")
}

func main() {
	// Servir los archivos estáticos de la carpeta "web" (nuestro index.html)
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	// Ruta para procesar el formulario
	http.HandleFunc("/enviar-contacto", contactoHandler)

	port := ":8080"
	fmt.Printf("Servidor corriendo en http://localhost%s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
	}
}
