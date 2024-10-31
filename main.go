package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Youssef-Shehata/http-server/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("  ERROR : couldn't parse json\n")
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(code)
	w.Write(data)
}
func respondWithError(w http.ResponseWriter, code int, msg string) {
	type error struct {
		Msg string `json:"msg"`
	}
	w.WriteHeader(code)
	log.Printf("  ERROR : %v \n", msg)
	if code == 200 {
		respondWithJSON(w, code, error{Msg: msg})
	}
	w.WriteHeader(code)
}
func healthHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "server is all good\n")

}

type apiConfig struct {
	platform string
	homeHits atomic.Int32
	db       *sql.DB
	ctx      context.Context
}

func (cfg *apiConfig) middlewareMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.homeHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(403)
		return

	}
	w.WriteHeader(http.StatusOK)
	cfg.homeHits.Store(0)
	dbQueries := database.New(cfg.db)
    dbQueries.ResetUser(cfg.ctx)
    dbQueries.ResetTweets(cfg.ctx)


	fmt.Fprintf(w, "hits been reset to : 0")

}

func (cfg *apiConfig) metrics(w http.ResponseWriter, r *http.Request) {

	type page struct {
		Hits int32
	}
	tmpl, err := template.ParseFiles("./metrics.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, page{Hits: cfg.homeHits.Load()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type Tweet struct {
	Id         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	User_id    uuid.UUID `json:"user_id"`
	Body       string    `json:"body"`
}

func (cfg *apiConfig) tweet(w http.ResponseWriter, r *http.Request) {
	tweet := Tweet{}
	j := json.NewDecoder(r.Body)
	err := j.Decode(&tweet)

	if err != nil {
		respondWithError(w, 500, "parsing json in post request : api/validate")
		return
	}

	if len(tweet.Body) > 140 {
		respondWithError(w, 200, "cant exceed 140 characters in request body\n")
		return
	}

	if strings.ContainsAny(strings.ToLower(tweet.Body), "fuck") {
		tweet.Body = strings.ReplaceAll(tweet.Body, "fuck", "****")
	}
    
    db:= database.New(cfg.db)
    t ,err:= db.CreateTweet(cfg.ctx, database.CreateTweetParams{
    	UserID: uuid.NullUUID{UUID: tweet.User_id,Valid: true},
    	Body:  tweet.Body, 
    })
    if err != nil{
        respondWithError(w,http.StatusInternalServerError,"couldnt upload tweet")
        return
    }

	respondWithJSON(w, 200, t)
}
func (c *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {

	type params struct {
		Email string `json:"email"`
	}
	j := json.NewDecoder(r.Body)
	p := params{}
	err := j.Decode(&p)
	if err != nil || p.Email == "" {
		respondWithError(w, http.StatusBadRequest, "couldnt decode json in request ")
		return
	}
	dbQueries := database.New(c.db)
	newUser, err := dbQueries.CreateUser(c.ctx, p.Email)
	sh7tt := User(newUser)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt create new user")
		return
	}

	respondWithJSON(w, 200, sh7tt)
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func main() {

	godotenv.Load()
	ctx := context.Background()
	mux := http.NewServeMux()
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))

	if err != nil {
		log.Fatal("ERROR: connecting to db")
	}

	cfg := &apiConfig{ctx: ctx, platform: os.Getenv("PLATFORM"), homeHits: atomic.Int32{}, db: db}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app", cfg.middlewareMetrics(handler))
	mux.HandleFunc("GET /admin/health", healthHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.metrics)
	mux.HandleFunc("POST /admin/reset", cfg.reset)
	mux.HandleFunc("POST /api/tweet", cfg.tweet)
	mux.HandleFunc("POST /api/create_user", cfg.createUser)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "main page , welcome son \n")
	})

	//mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	if r.URL.Path != "/" {
	//		w.WriteHeader(404)
	//		fmt.Fprintf(w, "404 Page Not Found\n")
	//		return
	//	}
	//	fmt.Fprintf(w, "welcome to the home page\n")

	//})

	server := http.Server{Handler: mux, Addr: "localhost:8080"}
	e := server.ListenAndServe()

	if e != nil {
		log.Print("  ERROR : listeninig to request\n", err)

	}

}
