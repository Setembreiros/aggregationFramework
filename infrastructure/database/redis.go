package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	model "aggregationframework/internal/model/domain"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisCacheClient struct {
	ctx    context.Context
	client *redis.Client
}

type FollowersData struct {
	Followers      []model.Follower `json:"followers"`
	LastFollowerId string           `json:"lastFollowerId"`
}

func NewRedisClient(cacheUri, cachePassword string, ctx context.Context) *RedisCacheClient {
	redisConfig := &redis.Options{
		Addr:     cacheUri,
		Password: cachePassword,
		DB:       0, // Use default DB
		Protocol: 2, // Connection protocol
	}
	client := &RedisCacheClient{
		ctx:    ctx,
		client: redis.NewClient(redisConfig),
	}

	client.verifyConnection()

	return client
}

func (c *RedisCacheClient) verifyConnection() {
	err := c.client.Set(c.ctx, "foo", "bar", 10*time.Second).Err()
	if err != nil {
		log.Error().Stack().Err(err).Msg("Conection to Redis not stablished")
		panic(err)
	}

	_, err = c.client.Get(c.ctx, "foo").Result()
	if err != nil {
		log.Error().Stack().Err(err).Msg("Conection to Redis not stablished")
		panic(err)
	}
	log.Info().Msgf("Connection to Redis established.")
}

func (c *RedisCacheClient) Clean() {
	err := c.client.FlushDB(c.ctx).Err()
	if err != nil {
		log.Warn().Stack().Err(err).Msg("Failed to clean entire Redis cache")
		return
	}

	log.Info().Msg("Entire Redis cache cleaned successfully")
}

func (c *RedisCacheClient) SetUserFollowers(username string, lastFollowerId string, limit int, followers []model.Follower) {
	cacheKey := c.generateCacheKey(username, lastFollowerId, limit)

	newLastFollowerId := ""
	if len(followers) > 0 {
		newLastFollowerId = followers[len(followers)-1].Username
	}

	data := FollowersData{
		Followers:      followers,
		LastFollowerId: newLastFollowerId,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Warn().Stack().Err(err).Msg("Failed to serialize followers data")
		return
	}

	err = c.client.Set(c.ctx, cacheKey, jsonData, 5*time.Minute).Err()
	if err != nil {
		log.Warn().Stack().Err(err).Msg("Failed to set followers in cache")
	}
}

func (c *RedisCacheClient) GetUserFollowers(username string, lastFollowerId string, limit int) ([]model.Follower, string, bool) {
	cacheKey := c.generateCacheKey(username, lastFollowerId, limit)

	jsonData, err := c.client.Get(c.ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return []model.Follower{}, "", false
		}

		log.Warn().Stack().Err(err).Msg("Failed to retrieve followers from cache")
		return []model.Follower{}, "", false
	}

	var data FollowersData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Warn().Stack().Err(err).Msg("Failed to deserialize followers data")
		return []model.Follower{}, "", false
	}

	log.Info().Msgf("Data retrieve from cache for key %s", cacheKey)

	return data.Followers, data.LastFollowerId, true
}

func (c *RedisCacheClient) generateCacheKey(username string, lastFollowerId string, limit int) string {
	return fmt.Sprintf("followers:%s:%s:%d", username, lastFollowerId, limit)
}
