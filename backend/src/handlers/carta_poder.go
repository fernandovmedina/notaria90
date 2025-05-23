package handlers

import (
	"log"
	"net/http"

	"github.com/fernandovmedina/notaria90/backend/src/database"
)

func CartaPoder(w http.ResponseWriter, r *http.Request) {
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
