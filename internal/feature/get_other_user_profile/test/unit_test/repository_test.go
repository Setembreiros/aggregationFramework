package unit_test_get_other_user_profile

import (
	"errors"
	"testing"

	"aggregationframework/internal/feature/get_other_user_profile"
	mock_get_other_user_profile "aggregationframework/internal/feature/get_other_user_profile/test/mock"
	model "aggregationframework/internal/model/domain"

	"github.com/stretchr/testify/assert"
)

var repository *get_other_user_profile.GetOtherUserProfileRepository
var FollowConnector *mock_get_other_user_profile.MockfollowConnector
var readmodelsConnector *mock_get_other_user_profile.MockreadmodelsConnector

func setUpRepository(t *testing.T) {
	setUp(t)
	FollowConnector = mock_get_other_user_profile.NewMockfollowConnector(ctrl)
	readmodelsConnector = mock_get_other_user_profile.NewMockreadmodelsConnector(ctrl)
	repository = get_other_user_profile.NewGetOtherUserProfileRepository(FollowConnector, readmodelsConnector)
}

func TestGetOtherUserProfileFromRepository_WhenApiConnectorReturnsSuccess(t *testing.T) {
	setUpRepository(t)
	username := "usernameA"
	currentUsername := "usernameB"
	expectedUserProfile := &model.UserProfile{
		Username:                username,
		IsFollowedByCurrentUser: true,
	}
	FollowConnector.EXPECT().CheckIfRelationshipExists(username, currentUsername).Return(true, nil)
	readmodelsConnector.EXPECT().GetUserProfile(username).Return(expectedUserProfile, nil)

	userProfile, err := repository.GetOtherUserProfile(username, currentUsername)

	assert.Nil(t, err)
	assert.Equal(t, userProfile, expectedUserProfile)
}

func TestErrorOnGetOtherUserProfileFromRepository_WhenFollowConnectorFails(t *testing.T) {
	setUpRepository(t)
	username := "usernameA"
	currentUsername := "usernameB"
	expectedUserProfile := &model.UserProfile{
		Username:                username,
		IsFollowedByCurrentUser: true,
	}
	FollowConnector.EXPECT().CheckIfRelationshipExists(username, currentUsername).Return(false, errors.New("some error"))
	readmodelsConnector.EXPECT().GetUserProfile(username).Return(expectedUserProfile, nil)

	userProfile, err := repository.GetOtherUserProfile(username, currentUsername)

	assert.NotNil(t, err)
	assert.Nil(t, userProfile)
}

func TestErrorOnGetOtherUserProfileFromRepository_WhenReadmodelsConnectorFails(t *testing.T) {
	setUpRepository(t)
	username := "usernameA"
	currentUsername := "usernameB"
	FollowConnector.EXPECT().CheckIfRelationshipExists(username, currentUsername).Return(true, nil)
	readmodelsConnector.EXPECT().GetUserProfile(username).Return(nil, errors.New("some error"))

	userProfile, err := repository.GetOtherUserProfile(username, currentUsername)

	assert.NotNil(t, err)
	assert.Nil(t, userProfile)
}
