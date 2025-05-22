package handlers

import (
	"log"
	"net/http"

	"github.com/fernandovmedina/notaria90/backend/src/database"
)

func RegisterApointment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var nombre = r.FormValue("nombre")
	var apellidos = r.FormValue("apellidos")
	var dia = r.FormValue("dia")

	cookie, err := r.Cookie("usuario")
	if err != nil {
		log.Println("Error obteniendo cookie:", err.Error())
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	_, err = database.DB.Exec("INSERT INTO NOTARIA90.CITAS(ID_USUARIO, NOMBRE, APELLIDO, DIA) VALUES (?, ?, ?, ?)",
		cookie.Value, nombre, apellidos, dia)

	if err != nil {
		log.Println("Error insertando cita:", err.Error())
		http.Error(w, "Error al registrar cita", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
