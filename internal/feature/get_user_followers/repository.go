package get_user_followers

import (
	database "aggregationframework/internal/db"
	model "aggregationframework/internal/model/domain"
)

//go:generate mockgen -source=repository.go -destination=test/mock/repository.go

type FollowConnector interface {
	GetUserFollowerIds(username, lastFollowerId string, limit int) ([]string, string, error)
}

type readmodelsConnector interface {
	GetFollowersMetadata(username []string) ([]model.Follower, error)
}

type GetUserFollowersRepository struct {
	cache               *database.Cache
	FollowConnector     FollowConnector
	readmodelsConnector readmodelsConnector
}

func NewGetUserFollowersRepository(cache *database.Cache, FollowConnector FollowConnector, readmodelsConnector readmodelsConnector) *GetUserFollowersRepository {
	return &GetUserFollowersRepository{
		cache:               cache,
		FollowConnector:     FollowConnector,
		readmodelsConnector: readmodelsConnector,
	}
}

func (r *GetUserFollowersRepository) GetUserFollowers(username string, lastFollowerId string, limit int) ([]model.Follower, string, error) {
	followers, newLastFollowerId, found := r.cache.Client.GetUserFollowers(username, lastFollowerId, limit)
	if found {
		return followers, newLastFollowerId, nil
	}

	followerIds, newLastFollowerId, err := r.FollowConnector.GetUserFollowerIds(username, lastFollowerId, limit)
	if err != nil {
		return []model.Follower{}, "", err
	}

	followers, err = r.readmodelsConnector.GetFollowersMetadata(followerIds)
	if err != nil {
		return []model.Follower{}, "", err
	}

	r.cache.Client.SetUserFollowers(username, lastFollowerId, limit, followers)

	return followers, newLastFollowerId, nil
}
