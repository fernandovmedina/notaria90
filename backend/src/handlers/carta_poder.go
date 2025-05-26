package handlers

import (
	"log"
	"net/http"

	"github.com/fernandovmedina/notaria90/backend/src/database"
)

func CartaPoder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "https://notaria90.vercel.app")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var err error

	var poderante = r.FormValue("poderante")
	var apoderado = r.FormValue("apoderado")
	var domicilio = r.FormValue("domicilio")
	var telefono = r.FormValue("telefono")
	var rfc = r.FormValue("rfc")

	cookie, err := r.Cookie("usuario")

	if err != nil {
		log.Println(err.Error())
	}

	_, err = database.DB.Exec("INSERT INTO CARTA_PODER(ID_USUARIO,NOMBRE_PODERANTE,NOMBRE_APODERADO,DOMICILIO,TELEFONO,RFC)VALUES(?,?,?,?,?,?)", cookie.Value, poderante, apoderado, domicilio, telefono, rfc)

	if err != nil {
		log.Println(err.Error())
	}
}
