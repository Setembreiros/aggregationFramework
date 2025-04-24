package database

import model "aggregationframework/internal/model/domain"

//go:generate mockgen -source=cache.go -destination=test/mock/cache.go

type Cache struct {
	Client CacheClient
}

type CacheClient interface {
	Clean()
	SetUserFollowers(username string, lastFollowerId string, limit int, followers []model.Follower)
	GetUserFollowers(username string, lastFollowerId string, limit int) ([]model.Follower, string, bool)
	SetUserFollowees(username string, lastFolloweeId string, limit int, followees []model.Followee)
	GetUserFollowees(username string, lastFolloweeId string, limit int) ([]model.Followee, string, bool)
}

func NewCache(client CacheClient) *Cache {
	return &Cache{
		Client: client,
	}
}
