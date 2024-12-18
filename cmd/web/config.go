package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Youssef-Shehata/yapocalypse/internal/database"
	"github.com/Youssef-Shehata/yapocalypse/pkg/logger"
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
	rdb      *redis.Client
	query    *database.Queries
	ctx      context.Context
	secret   string
	api_key  string
    logger   *logger.Logger
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

    logger ,err:= logger.NewLogger("./server.log")
    if err != nil{
        log.Printf("couldnt open/create log file %v" , err)
    }

	rdb := newRedisClient()
	cfg := &apiConfig{ctx: ctx, platform: PLATFORM,  secret: SECRET, query: query, api_key: PREMUIM_API_KEY, rdb: rdb, logger: logger}
	return cfg, mux

}
