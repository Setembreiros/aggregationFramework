package get_other_user_profile

import (
	model "aggregationframework/internal/model/domain"
)

//go:generate mockgen -source=repository.go -destination=test/mock/repository.go

type followConnector interface {
	CheckIfRelationshipExists(followeeId, followerId string) (bool, error)
}

type readmodelsConnector interface {
	GetUserProfile(username string) (*model.UserProfile, error)
}

type GetOtherUserProfileRepository struct {
	followConnector     followConnector
	readmodelsConnector readmodelsConnector
}

func NewGetOtherUserProfileRepository(followConnector followConnector, readmodelsConnector readmodelsConnector) *GetOtherUserProfileRepository {
	return &GetOtherUserProfileRepository{
		followConnector:     followConnector,
		readmodelsConnector: readmodelsConnector,
	}
}

func (r *GetOtherUserProfileRepository) GetOtherUserProfile(username string, currentUsername string) (*model.UserProfile, error) {
	isFollowedChan := make(chan bool, 1)
	isFollowedErrChan := make(chan error, 1)
	userProfileChan := make(chan *model.UserProfile, 1)
	userProfileErrChan := make(chan error, 1)

	go r.checkIfRelationshipExistsGoroutine(username, currentUsername, isFollowedChan, isFollowedErrChan)
	go r.getUserProfileGoroutine(username, userProfileChan, userProfileErrChan)

	isFollowedByCurrentUser := <-isFollowedChan
	isFollowedErr := <-isFollowedErrChan
	userProfile := <-userProfileChan
	userProfileErr := <-userProfileErrChan

	if isFollowedErr != nil {
		return nil, isFollowedErr
	}

	if userProfileErr != nil {
		return nil, userProfileErr
	}

	userProfile.IsFollowedByCurrentUser = isFollowedByCurrentUser
	return userProfile, nil
}

func (r *GetOtherUserProfileRepository) checkIfRelationshipExistsGoroutine(username string, currentUsername string, resultChan chan<- bool, errChan chan<- error) {
	isFollowed, err := r.followConnector.CheckIfRelationshipExists(username, currentUsername)
	resultChan <- isFollowed
	errChan <- err
}

func (r *GetOtherUserProfileRepository) getUserProfileGoroutine(username string, resultChan chan<- *model.UserProfile, errChan chan<- error) {
	userProfile, err := r.readmodelsConnector.GetUserProfile(username)
	resultChan <- userProfile
	errChan <- err
}
