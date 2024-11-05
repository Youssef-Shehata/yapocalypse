package main

import (
	"encoding/json"
	"fmt"
	"github.com/Youssef-Shehata/yapocalypse/internal/auth"
	"github.com/Youssef-Shehata/yapocalypse/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	Username  string    `json:"username"`
	Premuim   bool      `json:"premuim"`
}

func (cfg *apiConfig) UpdateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	userId, ok := r.Context().Value("userid").(uuid.UUID)
	if !ok {
		log.Println("  ERROR: parsing context token")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var p params
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("  ERROR: parsing json in request :%v", err)
		http.Error(w, "failed to parse json", http.StatusBadRequest)
		return
	}

	password, err := auth.HashPassword(p.Password)
	if err != nil {
		log.Printf("  ERROR: hashing password :%v", err.Error())
		http.Error(w, fmt.Sprintf("couldnt hash password %v", err.Error()), http.StatusBadRequest)
		return
	}

	user, err := cfg.query.UpdateUser(cfg.ctx, database.UpdateUserParams{
		Email:    p.Email,
		Password: password,
		ID:       userId,
	})

	if err != nil {
		log.Printf("  ERROR: updating user :%v", err.Error())
		http.Error(w, fmt.Sprintf("couldnt update user %v", err.Error()), http.StatusBadRequest)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        userId,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Username:  user.Username,
		Premuim:   user.Premuim,
	})
}

func (cfg *apiConfig) SubscribeToPremuim(w http.ResponseWriter, r *http.Request) {
	key := auth.GetAPIKey(r.Header)
	if key != cfg.api_key {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	type userId struct {
		User_id string `json:"user_id"`
	}
	type params struct {
		Event string `json:"event"`
		Data  userId `json:"data"`
	}

	var p params
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("  ERROR: couldnt parse request %v \n", err.Error())
		http.Error(w, "failed parse json :", http.StatusBadRequest)
	}

	if p.Event != "user.upgraded" {
		log.Printf("  ERROR: premuim webhook recieved wrong event\n")
		http.Error(w, "", http.StatusNoContent)
		return
	}
	id, err := uuid.Parse(p.Data.User_id)

	if err != nil {
		log.Printf("  ERROR: couldnt parse id %v \n", err.Error())
		http.Error(w, "failed parse id :", http.StatusBadRequest)
		return
	}

	if error := cfg.query.SubscribeToPremuim(cfg.ctx, id); error != nil {
		log.Printf("  ERROR: couldnt subscribe %v \n", error.Error())
		http.Error(w, "", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

}
func (cfg *apiConfig) signUp(w http.ResponseWriter, r *http.Request) {

	type params struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in"`
	}
	var p params
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("  ERROR bad request api/createUser: %v \n", err)
		http.Error(w, fmt.Sprint("bad request : ", err.Error()), http.StatusBadRequest)
		return
	}
	if p.Email == "" {
		log.Printf("  ERROR bad request to api/createUser: empty email field\n")
		http.Error(w, "email can't be empty", http.StatusBadRequest)
		return

	}

	hashedPass, err := auth.HashPassword(p.Password)
	if err != nil {
		log.Printf("  ERROR failed to create new user : %v \n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := cfg.query.CreateUser(cfg.ctx, database.CreateUserParams{Email: p.Email, Password: hashedPass})

	if err != nil {
		time.Sleep(time.Second)
		log.Printf("  ERROR failed to create new user : %v \n", err)
		http.Error(w, "couldnt create new user", http.StatusInternalServerError)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, p.ExpiresInSeconds)
	if err != nil {
		log.Printf("  ERROR making token : %v", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
		Username:  user.Username,
		Premuim:   user.Premuim,
	})
}
func (cfg *apiConfig) logIn(w http.ResponseWriter, r *http.Request) {

	type params struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	var p params
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("  ERROR bad request api/createUser: %v \n", err)
		http.Error(w, fmt.Sprint("bad request : ", err.Error()), http.StatusBadRequest)
		return
	}

	if p.Email == "" {
		log.Printf("  ERROR bad request to api/createUser: empty email field\n")
		http.Error(w, "email can't be empty", http.StatusBadRequest)
		return
	}
	user, err := cfg.query.GetUserByEmail(cfg.ctx, p.Email)
	if err != nil {
		log.Printf("  ERROR Wrong email or password \n %v", err)
		http.Error(w, "Wrong Email or Password", http.StatusUnauthorized)
		return
	}

	if error := auth.CheckHashedPassword(p.Password, user.Password); err != nil {
		log.Printf("  ERROR Wrong email or password \n %v", error)
		http.Error(w, "Wrong Email or Password", http.StatusUnauthorized)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, p.ExpiresInSeconds)
	if err != nil {
		log.Printf("  ERROR making token : %v", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return

	}
	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})
}
