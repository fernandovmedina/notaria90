package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fernandovmedina/notaria90/backend/src/database"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	r.ParseForm()

	nombre := r.FormValue("nombre")
	correo := r.FormValue("correo")
	password := r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error al encriptar la contraseña:", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec("INSERT INTO NOTARIA90.USUARIOS (NOMBRE, CORREO, PASSWORD) VALUES (?, ?, ?)", nombre, correo, string(hashedPassword))
	if err != nil {
		log.Println("Error al registrar usuario:", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	var id int
	err = database.DB.QueryRow("SELECT ID_USUARIO FROM NOTARIA90.USUARIOS WHERE CORREO = ?", correo).Scan(&id)
	if err != nil {
		log.Println("Error al obtener ID de usuario:", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "usuario",
		Value:    strconv.Itoa(id),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false, // Considerar true en producción con HTTPS
	}
	http.SetCookie(w, &cookie)
	// Respondemos con un OK para indicar que el registro fue exitoso
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Registro exitoso")) // Opcional: puedes enviar un mensaje
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println("Error al parsear formulario en Login:", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	correo := r.FormValue("correo")
	password := r.FormValue("password")

	var hashedPassword string
	var id int
	err := database.DB.QueryRow("SELECT PASSWORD, ID_USUARIO FROM NOTARIA90.USUARIOS WHERE CORREO=?", correo).Scan(&hashedPassword, &id)
	if err != nil {
		log.Println("Usuario no encontrado:", err)
		http.Error(w, "Usuario o contraseña incorrectos", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		log.Println("Contraseña incorrecta:", err)
		http.Error(w, "Usuario o contraseña incorrectos", http.StatusUnauthorized)
		return
	}

	cookie := &http.Cookie{
		Name:     "usuario",
		Value:    strconv.Itoa(id),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false, // Asegúrate que coincida con tu entorno
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
	log.Println("Cookie creada para usuario ID:", id)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login exitoso"))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Establecer cabeceras comunes primero
	origin := "http://localhost:4321" // Asegúrate que este sea el origen de tu frontend Astro
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", origin)

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Métodos que permites para esta ruta
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")  // Cabeceras que permites en la solicitud real
		w.WriteHeader(http.StatusOK)
		return
	}

	// Si no es OPTIONS, es la solicitud POST para Logout
	// (Asegúrate que las cabeceras de arriba como Allow-Origin y Allow-Credentials también se apliquen aquí)
	// w.Header().Set("Access-Control-Allow-Origin", origin) // Ya está arriba, pero no hace daño repetirla si hay dudas.
	// w.Header().Set("Access-Control-Allow-Credentials", "true") // Ya está arriba.

	cookie := http.Cookie{
		Name:     "usuario",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,                // Coincidir con la cookie de Login
		SameSite: http.SameSiteLaxMode, // Coincidir con la cookie de Login
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout exitoso"))
	log.Println("Usuario deslogueado")
}
