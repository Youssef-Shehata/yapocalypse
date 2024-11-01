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
	"sync/atomic"
	"time"

	"github.com/Youssef-Shehata/http-server/internal/auth"
	"github.com/Youssef-Shehata/http-server/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, err := json.Marshal(payload)
	if err != nil {
		log.Printf(fmt.Sprintf("  ERROR : couldn't parse json : %v\n", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(res)
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

type Tweet struct {
	ID        uuid.UUID `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) getTweetById(w http.ResponseWriter, r *http.Request) {
	log.Printf("the id in request was : %v \n", r.PathValue("id"))
	stringId := r.PathValue("id")
	byteId := []byte(stringId)
	id, err := uuid.ParseBytes(byteId)

	if err != nil {
		log.Printf("  ERROR Invalid id: id(%v) \n %v \n", stringId, err.Error())
		http.Error(w, fmt.Sprintln("Invalid Id :", err.Error()), http.StatusBadRequest)
		return
	}
	db := database.New(cfg.db)
	tweet, err := db.GetTweetById(cfg.ctx, id)
	if err != nil {
		log.Printf("  ERROR tweet not found %v \n%v\n", id, err.Error())
		http.Error(w, "", http.StatusNotFound)
		return
	}
	respondWithJSON(w, http.StatusOK, Tweet{
		ID:        tweet.ID,
		CreatedAt: tweet.CreatedAt,
		UpdatedAt: tweet.UpdatedAt,
		Body:      tweet.Body,
		UserID:    tweet.UserID,
	})
}
func (cfg *apiConfig) getTweets(w http.ResponseWriter, r *http.Request) {
	type params struct {
		User_id uuid.UUID `json:"user_id"`
	}
	user := params{}
	j := json.NewDecoder(r.Body)

	err := j.Decode(&user)
	log.Printf("user requesting his tweets : %v", user)
	if err != nil {
		log.Printf("  ERROR parsing json in request api/tweet : %v \n", err)
		http.Error(w, fmt.Sprintln("failed to parse request", err.Error()), http.StatusBadRequest)
		return
	}

	db := database.New(cfg.db)

	tweets, err := db.GetTweets(cfg.ctx, user.User_id)
	if err != nil {
		log.Printf("  ERROR couldnt get tweets: %v \n", err)
		http.Error(w, fmt.Sprintln("failed to get tweets :", err.Error()), http.StatusInternalServerError)
		return
	}

	if len(tweets) == 0 {
		respondWithJSON(w, http.StatusOK, "user has no tweets")
		return
	}

	var resTweets []Tweet
	for _, tweet := range tweets {

		resTweets = append(resTweets, Tweet{
			ID:        tweet.ID,
			CreatedAt: tweet.CreatedAt,
			UpdatedAt: tweet.UpdatedAt,
			Body:      tweet.Body,
			UserID:    tweet.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, resTweets)
}
func (cfg *apiConfig) tweet(w http.ResponseWriter, r *http.Request) {
	tweet := Tweet{}
	j := json.NewDecoder(r.Body)

	err := j.Decode(&tweet)

	if err != nil {
		log.Printf("  ERROR parsing json in request api/tweet : %v \n", err)
		http.Error(w, fmt.Sprintln("failed to parse request", err.Error()), http.StatusBadRequest)
		return
	}

	log.Printf("tweeting : %+v \n", tweet)
	if len(tweet.Body) > 140 {
		log.Printf("  ERROR request has more than 140 characters \n")
		http.Error(w, "Tweet cant exceed 140 characters in request body \n", http.StatusBadRequest)
		return
	}

	db := database.New(cfg.db)
	t, err := db.CreateTweet(cfg.ctx, database.CreateTweetParams{
		UserID: tweet.UserID,
		Body:   tweet.Body,
	})
	if err != nil {
		log.Printf("  ERROR couldnt store tweet in db : %v \n", err)
		http.Error(w, fmt.Sprintln("failed to store tweet :", err.Error()), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, 200, t)
}
func (c *apiConfig) signUp(w http.ResponseWriter, r *http.Request) {

	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	j := json.NewDecoder(r.Body)
	p := params{}
	err := j.Decode(&p)
	if err != nil || p.Email == "" {
		log.Printf("  ERROR bad request api/createUser: %v \n", err)
		http.Error(w, fmt.Sprint("bad request : ", err.Error()), http.StatusBadRequest)
		return
	}
	dbQueries := database.New(c.db)
	hashedPass, err := auth.HashPassword(p.Password)
	if err != nil {
		log.Printf("  ERROR failed to create new user : %v \n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	user, err := dbQueries.CreateUser(c.ctx, database.CreateUserParams{Email: p.Email, Password: hashedPass})
	sh7tt := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	if err != nil {
		log.Printf("  ERROR failed to create new user : %v \n", err)
		http.Error(w, "couldnt create new user", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, 200, sh7tt)
}
func (cfg *apiConfig) logIn(w http.ResponseWriter, r *http.Request) {

	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	j := json.NewDecoder(r.Body)
	p := params{}
	err := j.Decode(&p)

	if err != nil || p.Email == "" {
		log.Printf("  ERROR bad request api/createUser: %v \n", err)
		http.Error(w, fmt.Sprint("bad request : ", err.Error()), http.StatusBadRequest)
		return
	}
	db := database.New(cfg.db)
	user, err := db.GetUserByEmail(cfg.ctx, p.Email)
	if err != nil {
		log.Printf("  ERROR Wrong email or password \n %v", err)
		http.Error(w, "Wrong Email or Password", http.StatusUnauthorized)
		return

	}

	error := auth.CheckHashedPassword(p.Password, user.Password)
	if error != nil {
		log.Printf("  ERROR Wrong email or password \n %v", error)
		http.Error(w, "Wrong Email or Password", http.StatusUnauthorized)
		return
	}
    respondWithJSON(w, http.StatusOK ,User{
    	ID:        user.ID,
    	CreatedAt: user.CreatedAt,
    	UpdatedAt: user.UpdatedAt,
    	Email:     user.Email,
    })
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

	mux.HandleFunc("POST /api/v1/tweets", cfg.tweet)
	mux.HandleFunc("GET /api/v1/tweets", cfg.getTweets)
	mux.HandleFunc("GET /api/v1/tweets/{id}", cfg.getTweetById)
	mux.HandleFunc("POST /api/v1/signup", cfg.signUp)
	mux.HandleFunc("POST /api/v1/login", cfg.logIn)

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
