package main

import (
	"encoding/json"
	"fmt"
	"github.com/Youssef-Shehata/yapocalypse/cmd/types"
	"github.com/Youssef-Shehata/yapocalypse/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type Yap = types.Yap

func (cfg *apiConfig) getYapById(w http.ResponseWriter, r *http.Request) {
	stringId := r.PathValue("yap_id")

	id, err := uuid.Parse(stringId)
	if err != nil {
		log.Printf("  ERROR: Invalid id: (%v) \n %v \n", stringId, err.Error())
		http.Error(w, fmt.Sprintln("Invalid Id :", err.Error()), http.StatusBadRequest)
		return
	}

	yap, err := cfg.query.GetYapById(cfg.ctx, id)
	if err != nil {
		log.Printf("  ERROR yap not found %v \n%v\n", id, err.Error())
		http.Error(w, "", http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, Yap{
		ID:        yap.ID,
		CreatedAt: yap.CreatedAt,
		UpdatedAt: yap.UpdatedAt,
		Body:      yap.Body,
		UserID:    yap.UserID,
	})
}
func (cfg *apiConfig) getYaps(w http.ResponseWriter, r *http.Request) {
	urlQueries := r.URL.Query()
	authorId := urlQueries.Get("user_id")
	var resYaps []Yap

	userId, err := uuid.Parse(authorId)
	if err != nil {
		log.Printf("  ERROR: Invalid id: (%v) \n %v \n", authorId, err.Error())
		http.Error(w, fmt.Sprintln("Invalid Id :", err.Error()), http.StatusBadRequest)
	}

	yaps, err := cfg.query.GetYapsByUserId(cfg.ctx, userId)
	if err != nil {
		log.Printf("  ERROR couldnt get yaps: %v \n", err)
		http.Error(w, fmt.Sprintln("failed to get yaps:", err.Error()), http.StatusInternalServerError)
		return
	}

	for _, yap := range yaps {
		resYaps = append(resYaps, Yap{
			ID:        yap.ID,
			CreatedAt: yap.CreatedAt,
			UpdatedAt: yap.UpdatedAt,
			Body:      yap.Body,
			UserID:    yap.UserID,
		})
	}
}

func (cfg *apiConfig) CreateYap(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body string `json:"body"`
	}

	userId, ok := r.Context().Value("userid").(uuid.UUID)
	if !ok {
		log.Printf("  ERROR: User isnt Authenticated \n")
		http.Error(w, "Failed to parse token", http.StatusUnauthorized)
	}

	p := params{}
	j := json.NewDecoder(r.Body)

	error := j.Decode(&p)
	if error != nil {
		log.Printf("  ERROR: parsing json in request : %v \n", error)
		http.Error(w, fmt.Sprintln("failed to parse request", error.Error()), http.StatusBadRequest)
		return
	}

	if len(p.Body) > 60 {
		log.Printf("  ERROR: request has more than 60 characters \n")
		http.Error(w, "damn boi, chill; Yap cant exceed 60 characters \n", http.StatusBadRequest)
		return
	}
    //+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
    //DUAL WRITE PROBLEM OVER HERE 
    //+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

    //DB WRITE
	db_yap, err := cfg.query.NewYap(cfg.ctx, database.NewYapParams{
		UserID: userId,
		Body:   p.Body,
	})

	if err != nil {
		log.Printf("  ERROR: couldnt store yap in db : %v \n", err)
		http.Error(w, fmt.Sprintln("failed to store yap:", err.Error()), http.StatusInternalServerError)
		return
	}

	yap := Yap{
		ID:        db_yap.ID,
		UpdatedAt: db_yap.UpdatedAt,
		CreatedAt: db_yap.CreatedAt,
		Body:      db_yap.Body,
		UserID:    db_yap.UserID,
	}

    //KAFKA WRITE

	respondWithJSON(w, 200, yap)
}

func (cfg *apiConfig) DeleteYap(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userid").(uuid.UUID)
	if !ok {
		log.Println("  ERROR: parsing context token ")
		http.Error(w, "", http.StatusInternalServerError)
	}

	stringId := r.PathValue("id")
	byteId := []byte(stringId)
	yapId, err := uuid.ParseBytes(byteId)

	if err != nil {
		log.Printf("  ERROR: Invalid id: (%v) \n %v \n", stringId, err.Error())
		http.Error(w, fmt.Sprintln("Invalid Id :", err.Error()), http.StatusBadRequest)
		return
	}
	yap, yapNotFound := cfg.query.GetYapById(cfg.ctx, yapId)
	if yapNotFound != nil {

		log.Printf("  ERROR: couldnt delete yap%v \n", yapNotFound.Error())
		http.Error(w, fmt.Sprintln("failed to delete yap:", yapNotFound.Error()), http.StatusNotFound)
	}

	error := cfg.query.DeleteYap(cfg.ctx, database.DeleteYapParams{
		ID:     yap.ID,
		UserID: userId,
	})

	if error != nil {
		log.Printf("  ERROR: couldnt delete yap %v \n", error.Error())
		http.Error(w, fmt.Sprintln("failed to delete yap:", error.Error()), http.StatusBadRequest)
		return
	}

	w.WriteHeader(204)
}
