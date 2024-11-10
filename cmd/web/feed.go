package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Youssef-Shehata/yapocalypse/internal/database"
	"github.com/google/uuid"
)

type FeedParams struct {
	PageNumber int       `json:"page_number"`
	UserId     uuid.UUID `json:"user_id"`
}

func newFeedParams(numberQuery string, userIdQuery string) (FeedParams, error) {
	log.Printf("pageNum : %v , userId : %v ", numberQuery, userIdQuery)
	pageNumber, err := strconv.Atoi(numberQuery)
	if err != nil {
		log.Printf("  ERROR : failed to parse the provided page number : %v", err)
		pageNumber = 0
	}

	userId, err := uuid.Parse(userIdQuery)
	if err != nil {
		log.Printf("  ERROR : failed to parse the provided id : %v", err)
		return FeedParams{}, err
	}

	return FeedParams{PageNumber: pageNumber, UserId: userId}, nil

}

func (f FeedParams) GetCacheKey() string {

	numStr := strconv.Itoa(f.PageNumber)

	key := f.UserId.String() + numStr

	return key

}

func (cfg *apiConfig) getFeed(w http.ResponseWriter, r *http.Request) {

	feedParams, err := newFeedParams(r.URL.Query().Get("page_number"), r.URL.Query().Get("user_id"))
	if err != nil {
		http.Error(w, "failed to parse user id", http.StatusBadRequest)
		return
	}

	key := feedParams.GetCacheKey()

	//CACHE HIT
	cachedYaps, error := cacheGet(cfg, key)
	if error == nil {
		log.Println("cache hit")
		respondWithJSON(w, http.StatusOK, cachedYaps)
        return
	}

	//CACHE MISS
	//request from db
    log.Println("requesting yaps from db")
	yaps, err := cfg.query.GetFeed(cfg.ctx, database.GetFeedParams{
		UserID: feedParams.UserId,
		Offset: int32(feedParams.PageNumber) * 20,
	})

	if err != nil {
		log.Printf("  ERROR : failed to fetch user feed : %v", err)
		http.Error(w, "Failed to fetch feed of user( %v ) : ", http.StatusInternalServerError)
		return
	}

	var resYaps []Yap
	for _, yap := range yaps {
		resYaps = append(resYaps, Yap{
			ID:        yap.ID,
			UpdatedAt: yap.UpdatedAt,
			CreatedAt: yap.CreatedAt,
			Body:      yap.Body,
			UserID:    yap.UserID,
		})
	}

	if len(resYaps) != 0 {
		cacheSet(cfg, key, resYaps)
	}
    log.Println("responding with yaps from db")
	respondWithJSON(w, http.StatusOK, resYaps)
}
