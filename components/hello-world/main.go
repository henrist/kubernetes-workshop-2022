package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	greeting := os.Getenv("GREETING")
	if len(greeting) == 0 {
		greeting = "Hello"
	}

	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	log.Printf("Served request from %s\n", r.RemoteAddr)
	fmt.Fprintf(w, "%s from %s\n", greeting, host)
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Running application")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
