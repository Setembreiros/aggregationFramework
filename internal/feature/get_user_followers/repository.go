package get_user_followers

import model "aggregationframework/internal/model/domain"

//go:generate mockgen -source=repository.go -destination=test/mock/repository.go

type followerConnector interface {
	GetUserFollowerIds(username, lastFollowerId string, limit int) ([]string, string, error)
}

type readmodelsConnector interface {
	GetFollowersMetadata(username []string) ([]model.Follower, error)
}

type GetUserFollowersRepository struct {
	followerConnector   followerConnector
	readmodelsConnector readmodelsConnector
}

func NewGetUserFollowersRepository(followerConnector followerConnector, readmodelsConnector readmodelsConnector) *GetUserFollowersRepository {
	return &GetUserFollowersRepository{
		followerConnector:   followerConnector,
		readmodelsConnector: readmodelsConnector,
	}
}

func (r *GetUserFollowersRepository) GetUserFollowers(username string, lastFollowerId string, limit int) ([]model.Follower, string, error) {
	followerIds, newLastFollowerId, err := r.followerConnector.GetUserFollowerIds(username, lastFollowerId, limit)
	if err != nil {
		return []model.Follower{}, "", err
	}

	followers, err := r.readmodelsConnector.GetFollowersMetadata(followerIds)
	if err != nil {
		return []model.Follower{}, "", err
	}

	return followers, newLastFollowerId, nil
}
