package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Youssef-Shehata/yapocalypse/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var (
	DB_URL          = os.Getenv("DB_URL")
	PLATFORM        = os.Getenv("PLATFORM")
	SECRET          = os.Getenv("SECRET")
	PREMUIM_API_KEY = os.Getenv("PREMUIM_API_KEY")
)

type apiConfig struct {
	platform string
	homeHits atomic.Int32
	db       *sql.DB
	rdb      *redis.Client
	query    *database.Queries
	ctx      context.Context
	secret   string
	api_key  string
}

func Init() (*apiConfig, *http.ServeMux) {

	godotenv.Load()
	ctx := context.Background()
	mux := http.NewServeMux()
	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatal("ERROR: connecting to db")
	}
    query := database.New(db)


	rdb := newRedisClient()
	cfg := &apiConfig{ctx: ctx, platform: PLATFORM, homeHits: atomic.Int32{}, db: db, secret: SECRET, query: query, api_key: PREMUIM_API_KEY, rdb: rdb  }
	return cfg, mux
}