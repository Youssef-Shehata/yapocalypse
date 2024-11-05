package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/Youssef-Shehata/yapocalypse/internal/database"
	_ "github.com/lib/pq"
)

func (cfg *apiConfig) metrics(w http.ResponseWriter, r *http.Request) {
	type page struct {
		Hits int32
	}

	tmpl, err := template.ParseFiles("./metrics.html")
	if err != nil {
		log.Printf("  ERROR: /metrics failed to respond %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, page{Hits: cfg.homeHits.Load()})
	if err != nil {
		log.Printf("  ERROR: /metrics failed to respond %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "server is all good\n")
}
func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(403)
		return
	}

	cfg.homeHits.Store(0)
	dbQueries := database.New(cfg.db)
	dbQueries.ResetUser(cfg.ctx)
	dbQueries.ResetYaps(cfg.ctx)
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Home hits has been reset to : 0")
}

