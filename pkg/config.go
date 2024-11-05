package main
import (
	"context"
	"database/sql"
	"github.com/Youssef-Shehata/yapocalypse/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)


type apiConfig struct {
	platform string
	homeHits atomic.Int32
	db       *sql.DB
	query    *database.Queries
	ctx      context.Context
	secret   string
	api_key  string
}

func Init() (*apiConfig, *http.ServeMux) {

	godotenv.Load()
	ctx := context.Background()
	mux := http.NewServeMux()
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
    log.Printf("url : %v" , os.Getenv("DB_URL"))
	query := database.New(db)
	if err != nil {
		log.Fatal("ERROR: connecting to db")
	}

	cfg := &apiConfig{ctx: ctx, platform: os.Getenv("PLATFORM"), homeHits: atomic.Int32{}, db: db, secret: os.Getenv("SECRET"), query: query, api_key: os.Getenv("PREMUIM_API_KEY")}
	return cfg, mux
}

