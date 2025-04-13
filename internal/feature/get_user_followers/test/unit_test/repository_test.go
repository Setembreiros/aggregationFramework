package unit_test_get_user_followers

import (
	"errors"
	"testing"

	"aggregationframework/internal/feature/get_user_followers"
	mock_get_user_followers "aggregationframework/internal/feature/get_user_followers/test/mock"
	model "aggregationframework/internal/model/domain"

	"github.com/stretchr/testify/assert"
)

var repository *get_user_followers.GetUserFollowersRepository
var followerConnector *mock_get_user_followers.MockfollowerConnector
var readmodelsConnector *mock_get_user_followers.MockreadmodelsConnector

func setUpRepository(t *testing.T) {
	setUp(t)
	followerConnector = mock_get_user_followers.NewMockfollowerConnector(ctrl)
	readmodelsConnector = mock_get_user_followers.NewMockreadmodelsConnector(ctrl)
	repository = get_user_followers.NewGetUserFollowersRepository(followerConnector, readmodelsConnector)
}

func TestGetUserFollowersFromRepository_WhenApiConnectorReturnsSuccess(t *testing.T) {
	setUpRepository(t)
	username := "usernameA"
	lastFollowerId := "follower4"
	limit := 4
	expectedFollowerIds := []string{"follower5", "follower6", "follower7"}
	expectedLastFollowerId := "follower4"
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
	followerConnector.EXPECT().GetUserFollowerIds(username, lastFollowerId, limit).Return(expectedFollowerIds, expectedLastFollowerId, nil)
	readmodelsConnector.EXPECT().GetFollowersMetadata(expectedFollowerIds).Return(expectedFollowers, nil)

	followers, lastFollowerId, err := repository.GetUserFollowers(username, lastFollowerId, limit)

	assert.Nil(t, err)
	assert.Equal(t, followers, expectedFollowers)
	assert.Equal(t, lastFollowerId, expectedLastFollowerId)
}

func TestErrorOnGetUserFollowersFromRepository_WhenFollowerConnectorFails(t *testing.T) {
	setUpRepository(t)
	username := "usernameA"
	lastFollowerId := "follower4"
	limit := 4
	followerConnector.EXPECT().GetUserFollowerIds(username, lastFollowerId, limit).Return([]string{}, "", errors.New("some error"))

	followers, lastFollowerId, err := repository.GetUserFollowers(username, lastFollowerId, limit)

	assert.NotNil(t, err)
	assert.Equal(t, followers, []model.Follower{})
	assert.Equal(t, lastFollowerId, "")
}

func TestErrorOnGetUserFollowersFromRepository_WhenReadmodelsConnectorFails(t *testing.T) {
	setUpRepository(t)
	username := "usernameA"
	lastFollowerId := "follower4"
	limit := 4
	expectedFollowerIds := []string{"follower5", "follower6", "follower7"}
	expectedLastFollowerId := "follower4"
	followerConnector.EXPECT().GetUserFollowerIds(username, lastFollowerId, limit).Return(expectedFollowerIds, expectedLastFollowerId, nil)
	readmodelsConnector.EXPECT().GetFollowersMetadata(expectedFollowerIds).Return([]model.Follower{}, errors.New("some error"))

	followers, lastFollowerId, err := repository.GetUserFollowers(username, lastFollowerId, limit)

	assert.NotNil(t, err)
	assert.Equal(t, followers, []model.Follower{})
	assert.Equal(t, lastFollowerId, "")
}
