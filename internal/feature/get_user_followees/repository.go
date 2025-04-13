package get_user_followees

import (
	database "aggregationframework/internal/db"
	model "aggregationframework/internal/model/domain"
)

//go:generate mockgen -source=repository.go -destination=test/mock/repository.go

type followeeConnector interface {
	GetUserFolloweeIds(username, lastFolloweeId string, limit int) ([]string, string, error)
}

type readmodelsConnector interface {
	GetFolloweesMetadata(username []string) ([]model.Followee, error)
}

type GetUserFolloweesRepository struct {
	cache               *database.Cache
	followeeConnector   followeeConnector
	readmodelsConnector readmodelsConnector
}

func NewGetUserFolloweesRepository(cache *database.Cache, followeeConnector followeeConnector, readmodelsConnector readmodelsConnector) *GetUserFolloweesRepository {
	return &GetUserFolloweesRepository{
		cache:               cache,
		followeeConnector:   followeeConnector,
		readmodelsConnector: readmodelsConnector,
	}
}

func (r *GetUserFolloweesRepository) GetUserFollowees(username string, lastFolloweeId string, limit int) ([]model.Followee, string, error) {
	followees, newLastFolloweeId, found := r.cache.Client.GetUserFollowees(username, lastFolloweeId, limit)
	if found {
		return followees, newLastFolloweeId, nil
	}

	followeeIds, newLastFolloweeId, err := r.followeeConnector.GetUserFolloweeIds(username, lastFolloweeId, limit)
	if err != nil {
		return []model.Followee{}, "", err
	}

	followees, err = r.readmodelsConnector.GetFolloweesMetadata(followeeIds)
	if err != nil {
		return []model.Followee{}, "", err
	}

	r.cache.Client.SetUserFollowees(username, lastFolloweeId, limit, followees)

	return followees, newLastFolloweeId, nil
}
