package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Youssef-Shehata/yapocalypse/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) GetFollowers(w http.ResponseWriter, r *http.Request) {

	urlQueries := r.URL.Query()
	id := urlQueries.Get("id")

	var resUsers []User

	userId, err := uuid.Parse(id)
	if err != nil {
		log.Printf("  ERROR: Invalid id: (%v) \n %v \n", userId, err.Error())
		http.Error(w, fmt.Sprintln("Invalid Id :", err.Error()), http.StatusBadRequest)
	}

	users, err := cfg.query.GetFollowersOf(cfg.ctx, userId)
	if err != nil {
		log.Printf("  ERROR couldnt get yaps: %v \n", err)
		http.Error(w, fmt.Sprintln("failed to get yaps:", err.Error()), http.StatusInternalServerError)
		return
	}

	for _, user := range users {
		resUsers = append(resUsers, User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Username:  user.Username,
			Premuim:   user.Premuim,
		})
	}

    respondWithJSON(w,http.StatusOK , resUsers)
}

func (cfg *apiConfig) GetFollowees(w http.ResponseWriter, r *http.Request) {

	urlQueries := r.URL.Query()
	id := urlQueries.Get("id")

	var resUsers[]User

	userId, err := uuid.Parse(id)
	if err != nil {
		log.Printf("  ERROR: Invalid id: (%v) \n %v \n", userId, err.Error())
		http.Error(w, fmt.Sprintln("Invalid Id :", err.Error()), http.StatusBadRequest)
	}

	users, err := cfg.query.GetFolloweesOf(cfg.ctx, userId)
	if err != nil {
		log.Printf("  ERROR couldnt get yaps: %v \n", err)
		http.Error(w, fmt.Sprintln("failed to get yaps:", err.Error()), http.StatusInternalServerError)
		return
	}

	for _, user := range users {
		resUsers= append(resUsers,User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Username:  user.Username,
			Premuim:   user.Premuim,
		})
	}

    respondWithJSON(w,http.StatusOK , resUsers)
}

func (cfg *apiConfig) Follow(w http.ResponseWriter, r *http.Request) {

    type params struct{
        Follower_id string `json:"follower_id"`
        Followee_id string `json:"followee_id"`
    }

    var p params
    if err:= json.NewDecoder(r.Body).Decode(&p); err!=nil{
		log.Printf("  ERROR couldnt parse json %v \n", err)
		http.Error(w, fmt.Sprintln("failed to parse json :", err.Error()), http.StatusInternalServerError)
        return
    }

	followee_id, err := uuid.Parse(p.Followee_id)
	if err != nil {
		log.Printf("  ERROR: Invalid id: (%v) \n %v \n", followee_id, err.Error())
		http.Error(w, fmt.Sprintln("Invalid Id :", err.Error()), http.StatusBadRequest)
        return
	}

	followerId, err := uuid.Parse(p.Follower_id)
	if err != nil {
		log.Printf("  ERROR: Invalid id: (%v) \n %v \n", followerId, err.Error())
		http.Error(w, fmt.Sprintln("Invalid Id :", err.Error()), http.StatusBadRequest)
        return
	}

    if err := cfg.query.AddFollower(cfg.ctx , database.AddFollowerParams{
    	FollowerID: followerId,
    	FolloweeID: followee_id,
    }); err != nil {
		log.Printf("  ERROR couldn't follow %v : %v \n",followee_id, err)
		http.Error(w, fmt.Sprintln("failed to follow :",followee_id, err.Error()), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)

}

