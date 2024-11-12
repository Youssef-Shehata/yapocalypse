package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func newRedisClient() *redis.Client {
	redisHost := os.Getenv("REDIS_URL")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}

func cacheSet(cfg *apiConfig, key string, yaps []Yap) {

	jsonYaps, err := json.Marshal(yaps)
	if err != nil {
        cfg.logger.Log(ERROR , fmt.Errorf("Marshaling Json",err))
		return
	}
	if err := cfg.rdb.Set(cfg.ctx, key, jsonYaps, 10*time.Hour).Err(); err != nil {
        cfg.logger.Log(ERROR , fmt.Errorf("Setting to Redis",err))
		return
	}

}

func cacheGet(cfg *apiConfig, key string) ([]Yap, error) {

	cachedTweet, err := cfg.rdb.Get(cfg.ctx, key).Result()
	if err == redis.Nil {
        cfg.logger.Log(INFO, fmt.Errorf( "Redis Cache Miss"))
        return nil , err
	} else if err != nil {
        cfg.logger.Log(ERROR , fmt.Errorf("Redis Read",err))
        return nil , err
	}

	var yaps []Yap

	if err := json.NewDecoder(strings.NewReader(cachedTweet)).Decode(&yaps); err != nil {
        cfg.logger.Log(ERROR , fmt.Errorf("Redis Response",err))
		return nil, err
	}

	return yaps, nil
}
