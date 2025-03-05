package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/scriptoxin/yandex-liceum-go-calc/internal/handlers"
)

func main() {
	// Получаем порт из переменной окружения PORT, по умолчанию 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()
	// API для работы с выражениями
	r.HandleFunc("/api/v1/calculate", handlers.HandleCalculate).Methods("POST")
	r.HandleFunc("/api/v1/expressions", handlers.HandleExpressions).Methods("GET")
	r.HandleFunc("/api/v1/expressions/{id}", handlers.HandleExpressionByID).Methods("GET")

	// Endpoints для работы с задачами вычислений
	r.HandleFunc("/internal/task", handlers.HandleGetTask).Methods("GET")
	r.HandleFunc("/internal/task", handlers.HandlePostTask).Methods("POST")

	log.Printf("Orchestrator is running on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
