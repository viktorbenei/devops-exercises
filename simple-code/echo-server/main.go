package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

var port = "8182"

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message
	w.Write([]byte(message))
}

func main() {
	http.HandleFunc("/", sayHello)

	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	log.Printf("Starting (on port: %s) ...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
