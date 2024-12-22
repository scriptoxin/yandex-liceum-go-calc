package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/scriptoxin/yandex-liceum-go-calc/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/api/v1/calculate", handlers.CalculateHandler)

	fmt.Printf("Server is running on port %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
