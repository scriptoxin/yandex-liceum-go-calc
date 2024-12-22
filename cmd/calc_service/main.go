package main

import (
	"fmt"
	"net/http"

	"github.com/scriptoxin/yandex-liceum-go-calc/internal/handlers"
)

func main() {
	port := "8080" // Конфигурируемый порт

	http.HandleFunc("/api/v1/calculate", handlers.CalculateHandler)

	fmt.Printf("Server is running on port %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
