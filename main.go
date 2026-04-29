package main

import (
	"arc-events/handlers"
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	handlers.ArcEventsHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /arc-events", handlers.GetArcEvents)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handlers.CorsMiddleware(handlers.GzipMiddleware(mux)),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	log.Println("Server is running on localhost:8080")

	<-ctx.Done()
	log.Println("Stutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("Shutdown error: ", err)
	}
	log.Println("Done")
}
