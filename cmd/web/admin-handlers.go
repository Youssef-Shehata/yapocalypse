package main

import (
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "server is all good\n")
}
