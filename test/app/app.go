package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	hostname string
	port     string
	err      error
)

func logFatalIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func HandleHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%s | %s", hostname, time.Now())
}

func main() {
	hostname, err = os.Hostname()
	logFatalIf(err)

	http.HandleFunc("/", HandleHTTP)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.ListenAndServe(":"+port, nil)
}
