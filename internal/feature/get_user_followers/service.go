package get_user_followers

import (
	model "aggregationframework/internal/model/domain"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	GetUserFollowers(username string, lastPostId string, limit int) ([]model.Follower, string, error)
}

type GetUserFollowersService struct {
	repository Repository
}

func NewGetUserFollowersService(repository Repository) *GetUserFollowersService {
	return &GetUserFollowersService{
		repository: repository,
	}
}

func (s *GetUserFollowersService) GetUserFollowers(username string, lastFollowerId string, limit int) ([]model.Follower, string, error) {
	followers, lastFollowerId, err := s.repository.GetUserFollowers(username, lastFollowerId, limit)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting  %s's followers", username)
		return followers, lastFollowerId, err
	}

	return followers, lastFollowerId, nil
}
