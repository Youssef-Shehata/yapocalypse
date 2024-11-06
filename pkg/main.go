package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"

	"log"
	"net/http"
)

func main() {
	cfg, mux := Init()

	mux.Handle("/app", cfg.middlewareMetrics(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /admin/health", healthHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.metrics)

	mux.HandleFunc("GET /api/v1/yaps", cfg.getYaps)
	mux.HandleFunc("GET /api/v1/yaps/{id}", cfg.getYapById)

    //pulling feed as chunks of with limits 
    mux.HandleFunc("GET /api/v1/feed/{id}", cfg.getFeed)

	mux.HandleFunc("POST /admin/reset", cfg.reset)
	mux.HandleFunc("POST /api/v1/yaps", cfg.authMiddleware(cfg.CreateYap).ServeHTTP)
	mux.HandleFunc("POST /api/v1/signup", cfg.signUp)
	mux.HandleFunc("POST /api/v1/login", cfg.logIn)
	mux.HandleFunc("POST /api/v1/premuim/webhook", cfg.SubscribeToPremuim)

	mux.HandleFunc("GET /api/v1/followers/{id}", cfg.GetFollowers)
	mux.HandleFunc("GET /api/v1/followees/{id}", cfg.GetFollowees)
	mux.HandleFunc("POST /api/v1/followers"    , cfg.Follow)

    mux.HandleFunc("PUT /api/v1/users/", cfg.authMiddleware(cfg.premuimMiddleware(cfg.UpdateUser).ServeHTTP).ServeHTTP)

	mux.HandleFunc("DELETE /api/v1/yaps/{id}", cfg.authMiddleware(cfg.DeleteYap).ServeHTTP)




	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "main page , welcome son \n")
	})

	server := http.Server{Handler: mux, Addr: "localhost:8080"}
	err := server.ListenAndServe()

	if err != nil {
		log.Print("  ERROR: starting server:", err)

	}

}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	res, err := json.Marshal(payload)
	if err != nil {
		log.Printf("  ERROR: couldn't parse json : %v\n", err)
		http.Error(w, "failed to marshal json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)
}
