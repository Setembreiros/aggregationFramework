package unit_test_get_user_followers

import (
	"errors"
	"testing"

	database "aggregationframework/internal/db"
	mock_database "aggregationframework/internal/db/test/mock"
	"aggregationframework/internal/feature/get_user_followers"
	mock_get_user_followers "aggregationframework/internal/feature/get_user_followers/test/mock"
	model "aggregationframework/internal/model/domain"

	"github.com/stretchr/testify/assert"
)

var cacheClient *mock_database.MockCacheClient
var repository *get_user_followers.GetUserFollowersRepository
var FollowConnector *mock_get_user_followers.MockFollowConnector
var readmodelsConnector *mock_get_user_followers.MockreadmodelsConnector

func setUpRepository(t *testing.T) {
	setUp(t)
	cacheClient = mock_database.NewMockCacheClient(ctrl)
	FollowConnector = mock_get_user_followers.NewMockFollowConnector(ctrl)
	readmodelsConnector = mock_get_user_followers.NewMockreadmodelsConnector(ctrl)
	repository = get_user_followers.NewGetUserFollowersRepository(database.NewCache(cacheClient), FollowConnector, readmodelsConnector)
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
	cacheClient.EXPECT().GetUserFollowers(username, lastFollowerId, limit).Return([]model.Follower{}, "", false)
	FollowConnector.EXPECT().GetUserFollowerIds(username, lastFollowerId, limit).Return(expectedFollowerIds, expectedLastFollowerId, nil)
	readmodelsConnector.EXPECT().GetFollowersMetadata(expectedFollowerIds).Return(expectedFollowers, nil)
	cacheClient.EXPECT().SetUserFollowers(username, lastFollowerId, limit, expectedFollowers)

	followers, lastFollowerId, err := repository.GetUserFollowers(username, lastFollowerId, limit)

	assert.Nil(t, err)
	assert.Equal(t, followers, expectedFollowers)
	assert.Equal(t, lastFollowerId, expectedLastFollowerId)
}

func TestGetUserFollowersFromRepository_WhenCacheReturnsSuccess(t *testing.T) {
	setUpRepository(t)
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
	expectedLastFollowerId := "follower4"
	cacheClient.EXPECT().GetUserFollowers(username, lastFollowerId, limit).Return(expectedFollowers, expectedLastFollowerId, true)

	followers, lastFollowerId, err := repository.GetUserFollowers(username, lastFollowerId, limit)

	assert.Nil(t, err)
	assert.Equal(t, followers, expectedFollowers)
	assert.Equal(t, lastFollowerId, expectedLastFollowerId)
}

func TestErrorOnGetUserFollowersFromRepository_WhenFollowConnectorFails(t *testing.T) {
	setUpRepository(t)
	username := "usernameA"
	lastFollowerId := "follower4"
	limit := 4
	cacheClient.EXPECT().GetUserFollowers(username, lastFollowerId, limit).Return([]model.Follower{}, "", false)
	FollowConnector.EXPECT().GetUserFollowerIds(username, lastFollowerId, limit).Return([]string{}, "", errors.New("some error"))

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
	cacheClient.EXPECT().GetUserFollowers(username, lastFollowerId, limit).Return([]model.Follower{}, "", false)
	FollowConnector.EXPECT().GetUserFollowerIds(username, lastFollowerId, limit).Return(expectedFollowerIds, expectedLastFollowerId, nil)
	readmodelsConnector.EXPECT().GetFollowersMetadata(expectedFollowerIds).Return([]model.Follower{}, errors.New("some error"))

	followers, lastFollowerId, err := repository.GetUserFollowers(username, lastFollowerId, limit)

	assert.NotNil(t, err)
	assert.Equal(t, followers, []model.Follower{})
	assert.Equal(t, lastFollowerId, "")
}
