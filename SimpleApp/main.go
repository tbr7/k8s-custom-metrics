package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request")
	fmt.Fprintf(w, "Hello, the time is %s", time.Now().Format(time.RFC1123))
}
