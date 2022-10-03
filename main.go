package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-ping/ping"
)

var internalDNS = os.Getenv("INTERNAL_DNS")

func main() {
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	router := http.NewServeMux()
	router.HandleFunc("/ping", pingHeander)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()
	go gracefulShutdown(srv, quit, done)
	<-done
	log.Println("Server stopped")
}

func gracefulShutdown(server *http.Server, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	close(done)
}

func pingHeander(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ping" {
		http.Error(w, "404", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "405", http.StatusMethodNotAllowed)
		return
	}

	pinger, err := ping.NewPinger(internalDNS)
	if err != nil {
		http.Error(w, "502", http.StatusBadGateway)
		return
	}
	pinger.Count = 3
	err = pinger.Run()
	if err != nil {
		log.Printf("pinger error: %v\n", err)
		http.Error(w, "502", http.StatusBadGateway)
		return
	}
	pinger.Statistics()
	fmt.Fprintf(w, "Message: %+v", pinger.Statistics())
}
