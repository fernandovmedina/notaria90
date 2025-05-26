package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/fernandovmedina/notaria90/backend/src/database"
)

func GetTramites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	cookie, err := r.Cookie("usuario")
	if err != nil {
		log.Println("Error al obtener la cookie:", err.Error())
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	var jsonData sql.NullString
	err = database.DB.QueryRow("SELECT obtener_documentos_usuario(?)", cookie.Value).Scan(&jsonData)
	if err != nil {
		log.Println("Error al consultar los trámites:", err.Error())
		http.Error(w, "Error al obtener los trámites", http.StatusInternalServerError)
		return
	}

	if !jsonData.Valid {
		jsonData.String = `{"cartas_poder":[],"escrituras_publicas":[]}`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"body":   json.RawMessage(jsonData.String),
	})
}
