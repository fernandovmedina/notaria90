package handlers

import (
	"log"
	"net/http"

	"github.com/fernandovmedina/notaria90/backend/src/database"
)

func EscrituraPublica(w http.ResponseWriter, r *http.Request) {
	var err error

	var nombre = r.FormValue("nombre")
	var nacionalidad = r.FormValue("nacionalidad")
	var estadoCivil = r.FormValue("estado_civil")
	var ocupacion = r.FormValue("ocupacion")
	var curp = r.FormValue("curp")

	cookie, err := r.Cookie("usuario")

	if err != nil {
		log.Println(err.Error())
	}

	_, err = database.DB.Exec("INSERT INTO NOTARIA90.ESCRITURA_PUBLICA(ID_USUARIO,NOMBRE,NACIONALIDAD,ID_ESTADO_CIVIL,OCUPACION,CURP)VALUES(?,?,?,(SELECT ID_ESTADO_CIVIL FROM ESTADOS_CIVILES WHERE NOMBRE=?),?,?)", cookie.Value, nombre, nacionalidad, estadoCivil, ocupacion, curp)

	if err != nil {
		log.Println(err.Error())
	}
}
