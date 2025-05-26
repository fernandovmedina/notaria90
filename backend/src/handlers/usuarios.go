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
	w.Header().Set("Access-Control-Allow-Origin", "https://notaria90.vercel.app")
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
		Secure:   true,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Registro exitoso"))
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "https://notaria90.vercel.app")
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
		Secure:   true,
	}

	http.SetCookie(w, cookie)
	log.Println("Cookie creada para usuario ID:", id)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login exitoso"))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "https://notaria90.vercel.app" { // solo permitimos este origen explícitamente
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	cookie := http.Cookie{
		Name:     "usuario",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
	log.Println("Usuario deslogueado")
}
