package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Youssef-Shehata/yapocalypse/internal/auth"
	"github.com/Youssef-Shehata/yapocalypse/internal/database"
	"github.com/Youssef-Shehata/yapocalypse/pkg/types"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type User = types.User

func (cfg *apiConfig) UpdateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	userId, ok := r.Context().Value("userid").(uuid.UUID)
	if !ok {
		cfg.logger.Log(ERROR, fmt.Errorf("Parsing Token"))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var p params
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		cfg.logger.Log(ERROR, fmt.Errorf("Parsing json",err))
		http.Error(w, "failed to parse json", http.StatusBadRequest)
		return
	}

	password, err := auth.HashPassword(p.Password)
	if err != nil {
		cfg.logger.Log(ERROR, fmt.Errorf("Hashing Password",err))
		http.Error(w, fmt.Sprintf("couldnt hash password %v", err.Error()), http.StatusBadRequest)
		return
	}

	user, err := cfg.query.UpdateUser(cfg.ctx, database.UpdateUserParams{
		Email:    p.Email,
		Password: password,
		ID:       userId,
	})

	if err != nil {
		cfg.logger.Log(ERROR, fmt.Errorf("Updating User",err))
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
		cfg.logger.Log(ERROR, fmt.Errorf("Parsing Json",err))
		http.Error(w, "failed parse json :", http.StatusBadRequest)
	}

	if p.Event != "user.upgraded" {
		cfg.logger.Log(ERROR, fmt.Errorf("Wrong event"))
		http.Error(w, "", http.StatusNoContent)
		return
	}
	id, err := uuid.Parse(p.Data.User_id)

	if err != nil {
		cfg.logger.Log(ERROR, fmt.Errorf("Invalid Id",err))
		http.Error(w, "failed parse id :", http.StatusBadRequest)
		return
	}

	if error := cfg.query.SubscribeToPremuim(cfg.ctx, id); error != nil {
		cfg.logger.Log(ERROR, fmt.Errorf("failed to subscribe",error))
		http.Error(w, "", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

}
func (cfg *apiConfig) signUp(w http.ResponseWriter, r *http.Request) {

	type params struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		Username         string `json:"username"`
		ExpiresInSeconds int    `json:"expires_in"`
	}
	var p params
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		cfg.logger.Log(ERROR, fmt.Errorf("Parsing Json",err))
		http.Error(w, fmt.Sprint("bad request : ", err.Error()), http.StatusBadRequest)
		return
	}
	if p.Email == "" {
		cfg.logger.Log(ERROR, fmt.Errorf( "Empty Email"))
		http.Error(w, "email can't be empty", http.StatusBadRequest)
		return

	}
	if p.Username == "" {
		cfg.logger.Log(ERROR, fmt.Errorf( "Empty Username"))
		http.Error(w, "username can't be empty", http.StatusBadRequest)
		return

	}

	hashedPass, err := auth.HashPassword(p.Password)
	if err != nil {
		cfg.logger.Log(ERROR, fmt.Errorf("Hashing Password",err))
		http.Error(w, "failure hashing password", http.StatusInternalServerError)
		return
	}

	user, err := cfg.query.CreateUser(cfg.ctx, database.CreateUserParams{Email: p.Email, Password: hashedPass})

	// HOW TO KNOW USERNAME IS TAKEN WITH THIS VAGE ERROR SHIT
	if err != nil {
		time.Sleep(time.Second)
		cfg.logger.Log(ERROR, fmt.Errorf("Creating User",err))
		http.Error(w, "couldnt create new user", http.StatusInternalServerError)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, p.ExpiresInSeconds)
	if err != nil {
		cfg.logger.Log(ERROR, fmt.Errorf("Creating Token",err))
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
		cfg.logger.Log(ERROR, fmt.Errorf("Parsing Json",err))
		http.Error(w, fmt.Sprint("bad request : ", err.Error()), http.StatusBadRequest)
		return
	}

	if p.Email == "" {
		cfg.logger.Log(ERROR, fmt.Errorf( "Empty Email"))
		http.Error(w, "email can't be empty", http.StatusBadRequest)
		return
	}
	user, err := cfg.query.GetUserByEmail(cfg.ctx, p.Email)
	if err != nil {
		cfg.logger.Log(ERROR, fmt.Errorf( "Wrong Email Or Password"))
		http.Error(w, "Wrong Email or Password", http.StatusUnauthorized)
		return
	}

	if error := auth.CheckHashedPassword(p.Password, user.Password); error != nil {
		cfg.logger.Log(ERROR, fmt.Errorf("Wrong Email Or Password",error))
		http.Error(w, "Wrong Email or Password", http.StatusUnauthorized)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, p.ExpiresInSeconds)
	if err != nil {
		cfg.logger.Log(ERROR, fmt.Errorf("Creating Token",err))
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
