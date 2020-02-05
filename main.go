package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sparrc/go-ping"
)

var internalDNS = os.Getenv("internal_dns")

func main() {
	http.HandleFunc("/ping", pingHeander)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
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
	pinger.Run()
	pinger.Statistics()
	fmt.Fprintf(w, "Message: %+v", pinger.Statistics())
}
