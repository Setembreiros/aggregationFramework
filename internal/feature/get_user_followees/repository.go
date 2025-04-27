package get_user_followees

import (
	model "aggregationframework/internal/model/domain"
)

//go:generate mockgen -source=repository.go -destination=test/mock/repository.go

type followConnector interface {
	GetUserFolloweeIds(username, lastFolloweeId string, limit int) ([]string, string, error)
}

type readmodelsConnector interface {
	GetFolloweesMetadata(username []string) ([]model.Followee, error)
}

type GetUserFolloweesRepository struct {
	followConnector     followConnector
	readmodelsConnector readmodelsConnector
}

func NewGetUserFolloweesRepository(followConnector followConnector, readmodelsConnector readmodelsConnector) *GetUserFolloweesRepository {
	return &GetUserFolloweesRepository{
		followConnector:     followConnector,
		readmodelsConnector: readmodelsConnector,
	}
}

func (r *GetUserFolloweesRepository) GetUserFollowees(username string, lastFolloweeId string, limit int) ([]model.Followee, string, error) {
	followeeIds, newLastFolloweeId, err := r.followConnector.GetUserFolloweeIds(username, lastFolloweeId, limit)
	if err != nil {
		return []model.Followee{}, "", err
	}

	followees, err := r.readmodelsConnector.GetFolloweesMetadata(followeeIds)
	if err != nil {
		return []model.Followee{}, "", err
	}

	return followees, newLastFolloweeId, nil
}
