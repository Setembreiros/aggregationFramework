package unit_test_get_user_followers

import (
	"errors"
	"testing"

	"aggregationframework/internal/feature/get_user_followers"
	mock_get_user_followers "aggregationframework/internal/feature/get_user_followers/test/mock"
	model "aggregationframework/internal/model/domain"

	"github.com/stretchr/testify/assert"
)

var serviceRepository *mock_get_user_followers.MockRepository
var service *get_user_followers.GetUserFollowersService

func setUpService(t *testing.T) {
	setUp(t)
	serviceRepository = mock_get_user_followers.NewMockRepository(ctrl)
	service = get_user_followers.NewGetUserFollowersService(serviceRepository)
}

func TestGetUserFollowersWithService_WhenSuccess(t *testing.T) {
	setUpService(t)
	username := "usernameA"
	lastFollowerId := "follower4"
	limit := 4
	expectedFollowers := []model.Follower{
		{
			Username: "follower5",
			Fullname: "fullname5",
		},
		{
			Username: "follower6",
			Fullname: "fullname6",
		},
		{
			Username: "follower7",
			Fullname: "fullname7",
		},
	}
	expectedLastFollowerId := "follower7"
	serviceRepository.EXPECT().GetUserFollowers(username, lastFollowerId, limit).Return(expectedFollowers, expectedLastFollowerId, nil)

	followers, lastFollowerId, err := service.GetUserFollowers(username, lastFollowerId, limit)

	assert.Nil(t, err)
	assert.Equal(t, followers, expectedFollowers)
	assert.Equal(t, lastFollowerId, expectedLastFollowerId)
}

func TestErrorOnGetUserFollowersWithService_WhenGetUserFollowersFails(t *testing.T) {
	setUpService(t)
	username := "usernameA"
	lastFollowerId := "follower4"
	limit := 4
	expectedFollowers := []model.Follower{}
	expectedLastFollowerId := ""
	serviceRepository.EXPECT().GetUserFollowers(username, lastFollowerId, limit).Return(expectedFollowers, expectedLastFollowerId, errors.New("some error"))

	followers, lastFollowerId, err := service.GetUserFollowers(username, lastFollowerId, limit)

	assert.NotNil(t, err)
	assert.Equal(t, followers, expectedFollowers)
	assert.Equal(t, lastFollowerId, expectedLastFollowerId)
	assert.Contains(t, loggerOutput.String(), "Error getting  "+username+"'s followers")
}
