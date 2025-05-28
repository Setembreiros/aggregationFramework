package unit_test_get_other_user_profile

import (
	"errors"
	"testing"

	"aggregationframework/internal/feature/get_other_user_profile"
	mock_get_other_user_profile "aggregationframework/internal/feature/get_other_user_profile/test/mock"
	model "aggregationframework/internal/model/domain"

	"github.com/stretchr/testify/assert"
)

var serviceRepository *mock_get_other_user_profile.MockRepository
var service *get_other_user_profile.GetOtherUserProfileService

func setUpService(t *testing.T) {
	setUp(t)
	serviceRepository = mock_get_other_user_profile.NewMockRepository(ctrl)
	service = get_other_user_profile.NewGetOtherUserProfileService(serviceRepository)
}

func TestGetOtherUserProfileWithService_WhenSuccess(t *testing.T) {
	setUpService(t)
	expectedUsername := "usernameA"
	expectedCurrentUsername := "usernameB"
	expectedUserProfile := &model.UserProfile{
		Username:                expectedUsername,
		IsFollowedByCurrentUser: true,
	}
	serviceRepository.EXPECT().GetOtherUserProfile(expectedUsername, expectedCurrentUsername).Return(expectedUserProfile, nil)

	userProfile, err := service.GetOtherUserProfile(expectedUsername, expectedCurrentUsername)

	assert.Nil(t, err)
	assert.Equal(t, userProfile, expectedUserProfile)
}

func TestErrorOnGetOtherUserProfileWithService_WhenRepositoryFails(t *testing.T) {
	setUpService(t)
	expectedUsername := "usernameA"
	expectedCurrentUsername := "usernameB"
	serviceRepository.EXPECT().GetOtherUserProfile(expectedUsername, expectedCurrentUsername).Return(nil, errors.New("some error"))

	userProfile, err := service.GetOtherUserProfile(expectedUsername, expectedCurrentUsername)

	assert.NotNil(t, err)
	assert.Nil(t, userProfile)
	assert.Contains(t, loggerOutput.String(), "Error getting  "+expectedUsername+"'s profile")
}
