package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)
func main() {
    
	godotenv.Load()
    PORT := os.Getenv("PORT")


	cfg, mux := Init()
    //TODO: ADD AUTH MIDDLWARE
	mux.HandleFunc("GET /admin/health", healthHandler)

	mux.HandleFunc("GET /api/v1/yaps/user/{user_id}", cfg.getYaps)
	mux.HandleFunc("GET /api/v1/yaps/{yap_id}", cfg.getYapById)

	//pulling feed as chunks of with limits
	mux.HandleFunc("GET /api/v1/feed", cfg.getFeed)

	mux.HandleFunc("POST /api/v1/yaps", cfg.authMiddleware(cfg.CreateYap).ServeHTTP)
	mux.HandleFunc("POST /api/v1/signup", cfg.signUp)
	mux.HandleFunc("POST /api/v1/login", cfg.logIn)

    //trying out webhooks 
	mux.HandleFunc("POST /api/v1/premuim/webhook", cfg.SubscribeToPremuim)

	mux.HandleFunc("GET /api/v1/followers/{user_id}", cfg.GetFollowers)
	mux.HandleFunc("GET /api/v1/followees/{user_id}", cfg.GetFollowees)
	mux.HandleFunc("POST /api/v1/followers", cfg.Follow)

	mux.HandleFunc("PUT /api/v1/users/", cfg.authMiddleware(cfg.premuimMiddleware(cfg.UpdateUser).ServeHTTP).ServeHTTP)

	mux.HandleFunc("DELETE /api/v1/yaps/{id}", cfg.authMiddleware(cfg.DeleteYap).ServeHTTP)

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "main page , welcome son \n")
	})

	server := http.Server{Handler: mux, WriteTimeout: 10 * time.Second, ReadTimeout: 10 * time.Second, Addr:fmt.Sprintf("localhost:%v", PORT)}
    log.Printf("server listenning %v",server.Addr)
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
