package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Youssef-Shehata/yapocalypse/internal/auth"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)
func (cfg *apiConfig) premuimMiddleware(next http.HandlerFunc) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, ok := r.Context().Value("id").(uuid.UUID)
		if !ok {
			log.Println("  ERROR: parsing context token ")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		user, err := cfg.query.GetUserById(cfg.ctx, userId)

		if err != nil {
			log.Printf("  ERROR: gettin user : %v\n", err.Error())
			http.Error(w, "", http.StatusNotFound)
			return
		}

		if !user.Premuim {
			http.Error(w, "Premuim required to proceed", http.StatusPaymentRequired)
			return
		}

		next(w, r)
	})

}
func (cfg *apiConfig) authMiddleware(next http.HandlerFunc) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := auth.GetBearerToken(r.Header)

		userId, err := auth.ValidateJWT(token, cfg.secret)
		if err != nil {
			log.Printf("  ERROR: auth token : %v\n", err.Error())
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userid", userId)

		next(w, r.WithContext(ctx))
	})

}
func (cfg *apiConfig) middlewareMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.homeHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

