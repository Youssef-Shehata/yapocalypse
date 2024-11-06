package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getFeedInitialFeed(w http.ResponseWriter , r * http.Request){

    id := r.URL.Query().Get("id")
    userId,err := uuid.Parse(id)
    if err != nil{
        log.Printf("  ERROR : failed to parse the provided id : %v" , err)
        http.Error(w, "Failed to parse the provided id  "  , http.StatusInternalServerError)
        return
    }

    tweets , err := cfg.query.GetInitialFeed(cfg.ctx , userId)

    if err != nil{
        log.Printf("  ERROR : failed to fetch user feed : %v" , err)
        http.Error(w, "Failed to fetch feed of user( %v ) : "  , http.StatusInternalServerError)
        return
    }

    respondWithJSON(w , http.StatusOK , tweets)
}


func (cfg *apiConfig ) feedTheFeed(w http.ResponseWriter , r *http.Request){


}
