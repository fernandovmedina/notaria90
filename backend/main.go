package main

import (
	"log"
	"net/http"
	"time"

	"github.com/fernandovmedina/notaria90/backend/src/database"
	"github.com/fernandovmedina/notaria90/backend/src/handlers"
)

func main() {
	var err error

	if _, err = database.ConnectDB(); err != nil {
		log.Println(err.Error())
	}

	var mux = http.NewServeMux()

	mux.HandleFunc("GET /api", handlers.Welcome)
	mux.HandleFunc("POST /api/user/register", handlers.Register)
	mux.HandleFunc("POST /api/user/login", handlers.Login)
	mux.HandleFunc("/api/user/logout", handlers.Logout)
	mux.HandleFunc("POST /api/user/register-appointment", handlers.RegisterApointment)
	mux.HandleFunc("POST /api/user/escritura-publica", handlers.EscrituraPublica)
	mux.HandleFunc("POST /api/user/carta-poder", handlers.CartaPoder)

	handler := withCORS(mux)

	var server = http.Server{
		Addr:           ":8080",
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		Handler:        handler,
	}

	log.Println("Server running on http://127.0.0.1:8080/api")

	if err = server.ListenAndServe(); err != nil {
		log.Println(err.Error())
	}
}

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
