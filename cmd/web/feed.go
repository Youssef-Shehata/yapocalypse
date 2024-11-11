package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Youssef-Shehata/yapocalypse/internal/database"
	"github.com/Youssef-Shehata/yapocalypse/pkg/logger"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	ERROR = logger.ERROR
	INFO  = logger.INFO
)

type FeedParams struct {
	PageNumber int       `json:"page_number"`
	UserId     uuid.UUID `json:"user_id"`
}

func newFeedParams(numberQuery string, userIdQuery string) (FeedParams, error) {
	pageNumber, err := strconv.Atoi(numberQuery)
	if err != nil {
		pageNumber = 0
	}

	userId, err := uuid.Parse(userIdQuery)
	if err != nil {
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
		cfg.logger.Log(ERROR, fmt.Errorf("parsing user_id"))
		http.Error(w, "failed to parse user id", http.StatusBadRequest)
		return
	}

	key := feedParams.GetCacheKey()

	//CACHE HIT
	cachedYaps, error := cacheGet(cfg, key)
	if error == nil {
		cfg.logger.Log(INFO, fmt.Errorf("Cache Hit"))
		respondWithJSON(w, http.StatusOK, cachedYaps)
		return
	}

	//CACHE MISS
	//request from db
	cfg.logger.Log(ERROR, errors.Wrap(err, "Requesting yaps from db"))
	yaps, err := cfg.query.GetFeed(cfg.ctx, database.GetFeedParams{
		UserID: feedParams.UserId,
		Offset: int32(feedParams.PageNumber) * 20,
	})

	if err != nil {
		cfg.logger.Log(ERROR, errors.Wrap(err, "Fetching user feed"))
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
	cfg.logger.Log(INFO, fmt.Errorf("responding with yaps from db"))
	respondWithJSON(w, http.StatusOK, resYaps)
}
