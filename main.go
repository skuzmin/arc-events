package main

import (
	"arc-events/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /arc-events", handlers.GetArcEvents)
	fmt.Println("Server is running on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CorsMiddleware(handlers.GzipMiddleware(mux))))
}
