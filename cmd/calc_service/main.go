package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/scriptoxin/yandex-liceum-go-calc/internal/handlers"
	"github.com/scriptoxin/yandex-liceum-go-calc/pkg/db"
)

func main() {
	// Инициализируем БД (файл calc.db рядом с бинарником)
	if err := db.Init("calc.db"); err != nil {
		log.Fatalf("DB init failed: %v", err)
	}

	r := mux.NewRouter()
	// публичные эндпойнты
	r.HandleFunc("/api/v1/register", handlers.Register).Methods("POST")
	r.HandleFunc("/api/v1/login", handlers.Login).Methods("POST")

	// защищённая часть
	auth := r.PathPrefix("/api/v1").Subrouter()
	auth.Use(handlers.AuthMiddleware)
	auth.HandleFunc("/calculate", handlers.Calculate).Methods("POST")
	auth.HandleFunc("/expressions", handlers.GetExpressions).Methods("GET")
	auth.HandleFunc("/expressions/{id}", handlers.GetExpression).Methods("GET")

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
