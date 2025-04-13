package get_user_followees

import (
	model "aggregationframework/internal/model/domain"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	GetUserFollowees(username string, lastPostId string, limit int) ([]model.Followee, string, error)
}

type GetUserFolloweesService struct {
	repository Repository
}

func NewGetUserFolloweesService(repository Repository) *GetUserFolloweesService {
	return &GetUserFolloweesService{
		repository: repository,
	}
}

func (s *GetUserFolloweesService) GetUserFollowees(username string, lastFolloweeId string, limit int) ([]model.Followee, string, error) {
	followees, lastFolloweeId, err := s.repository.GetUserFollowees(username, lastFolloweeId, limit)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting  %s's followees", username)
		return followees, lastFolloweeId, err
	}

	return followees, lastFolloweeId, nil
}
