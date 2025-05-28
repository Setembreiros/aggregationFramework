package get_other_user_profile

import (
	model "aggregationframework/internal/model/domain"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	GetOtherUserProfile(username string, currentUsername string) (*model.UserProfile, error)
}

type GetOtherUserProfileService struct {
	repository Repository
}

func NewGetOtherUserProfileService(repository Repository) *GetOtherUserProfileService {
	return &GetOtherUserProfileService{
		repository: repository,
	}
}

func (s *GetOtherUserProfileService) GetOtherUserProfile(username string, currentUsername string) (*model.UserProfile, error) {
	userProfile, err := s.repository.GetOtherUserProfile(username, currentUsername)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting  %s's profile", username)
		return userProfile, err
	}

	return userProfile, nil
}
