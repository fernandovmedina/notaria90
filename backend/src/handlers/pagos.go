package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fernandovmedina/notaria90/backend/src/database"
)

type Pago struct {
	Id                 int    `json:"id_pago_escritura"`
	Nombre             string `json:"nombre"`
	Tarjeta            string `json:"tarjeta"`
	FechaDeVencimiento string `json:"fecha_de_vencimiento"`
	CVV                int    `json:"cvv"`
	Monto              string `json:"monto"`
	IdEscritura        *int   `json:"id_escritura_publica,omitempty"`
	IdCartaPoder       *int   `json:"id_carta_poder,omitempty"`
}

func RegistrarPago(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var pago Pago
	err := json.NewDecoder(r.Body).Decode(&pago)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	// Insertar el pago
	query := `INSERT INTO NOTARIA90.PAGOS (NOMBRE, TARJETA, FECHA_DE_VENCIMIENTO, CVV, MONTO) VALUES (?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, pago.Nombre, pago.Tarjeta, pago.FechaDeVencimiento, pago.CVV, pago.Monto)
	if err != nil {
		http.Error(w, "Error al insertar el pago", http.StatusInternalServerError)
		return
	}

	// Obtener el ID del nuevo pago
	idPago, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error al obtener ID del pago", http.StatusInternalServerError)
		return
	}

	// Asociar el pago con una escritura o una carta poder
	if pago.IdEscritura != nil {
		_, err = database.DB.Exec(
			`UPDATE NOTARIA90.ESCRITURA_PUBLICA SET ID_PAGO = ? WHERE ID_ESCRITURA_PUBLICA = ?`,
			idPago, *pago.IdEscritura,
		)
		if err != nil {
			http.Error(w, "Error al asociar pago con escritura", http.StatusInternalServerError)
			return
		}
	} else if pago.IdCartaPoder != nil {
		_, err = database.DB.Exec(
			`UPDATE NOTARIA90.CARTA_PODER SET ID_PAGO = ? WHERE ID_CARTA_PODER = ?`,
			idPago, *pago.IdCartaPoder,
		)
		if err != nil {
			http.Error(w, "Error al asociar pago con carta poder", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Debe proporcionar id_escritura_publica o id_carta_poder", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Pago registrado y asociado correctamente.")
}
