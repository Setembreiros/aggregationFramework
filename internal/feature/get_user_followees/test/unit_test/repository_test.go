package unit_test_get_user_followees

import (
	"errors"
	"testing"

	database "aggregationframework/internal/db"
	mock_database "aggregationframework/internal/db/test/mock"
	"aggregationframework/internal/feature/get_user_followees"
	mock_get_user_followees "aggregationframework/internal/feature/get_user_followees/test/mock"
	model "aggregationframework/internal/model/domain"

	"github.com/stretchr/testify/assert"
)

var cacheClient *mock_database.MockCacheClient
var repository *get_user_followees.GetUserFolloweesRepository
var followConnector *mock_get_user_followees.MockfollowConnector
var readmodelsConnector *mock_get_user_followees.MockreadmodelsConnector

func setUpRepository(t *testing.T) {
	setUp(t)
	cacheClient = mock_database.NewMockCacheClient(ctrl)
	followConnector = mock_get_user_followees.NewMockfollowConnector(ctrl)
	readmodelsConnector = mock_get_user_followees.NewMockreadmodelsConnector(ctrl)
	repository = get_user_followees.NewGetUserFolloweesRepository(database.NewCache(cacheClient), followConnector, readmodelsConnector)
}

func TestGetUserFolloweesFromRepository_WhenApiConnectorReturnsSuccess(t *testing.T) {
	setUpRepository(t)
	username := "usernameA"
	lastFolloweeId := "followee4"
	limit := 4
	expectedFolloweeIds := []string{"followee5", "followee6", "followee7"}
	expectedLastFolloweeId := "followee4"
	expectedFollowees := []model.Followee{
		{
			Username: "followee5",
			Fullname: "fullname5",
		},
		{
			Username: "followee6",
			Fullname: "fullname6",
		},
		{
			Username: "followee7",
			Fullname: "fullname7",
		},
	}
	cacheClient.EXPECT().GetUserFollowees(username, lastFolloweeId, limit).Return([]model.Followee{}, "", false)
	followConnector.EXPECT().GetUserFolloweeIds(username, lastFolloweeId, limit).Return(expectedFolloweeIds, expectedLastFolloweeId, nil)
	readmodelsConnector.EXPECT().GetFolloweesMetadata(expectedFolloweeIds).Return(expectedFollowees, nil)
	cacheClient.EXPECT().SetUserFollowees(username, lastFolloweeId, limit, expectedFollowees)

	followees, lastFolloweeId, err := repository.GetUserFollowees(username, lastFolloweeId, limit)

	assert.Nil(t, err)
	assert.Equal(t, followees, expectedFollowees)
	assert.Equal(t, lastFolloweeId, expectedLastFolloweeId)
}

func TestGetUserFolloweesFromRepository_WhenCacheReturnsSuccess(t *testing.T) {
	setUpRepository(t)
	username := "usernameA"
	lastFolloweeId := "followee4"
	limit := 4
	expectedFollowees := []model.Followee{
		{
			Username: "followee5",
			Fullname: "fullname5",
		},
		{
			Username: "followee6",
			Fullname: "fullname6",
		},
		{
			Username: "followee7",
			Fullname: "fullname7",
		},
	}
	expectedLastFolloweeId := "followee4"
	cacheClient.EXPECT().GetUserFollowees(username, lastFolloweeId, limit).Return(expectedFollowees, expectedLastFolloweeId, true)

	followees, lastFolloweeId, err := repository.GetUserFollowees(username, lastFolloweeId, limit)

	assert.Nil(t, err)
	assert.Equal(t, followees, expectedFollowees)
	assert.Equal(t, lastFolloweeId, expectedLastFolloweeId)
}

func TestErrorOnGetUserFolloweesFromRepository_WhenFollowConnectorFails(t *testing.T) {
	setUpRepository(t)
	username := "usernameA"
	lastFolloweeId := "followee4"
	limit := 4
	cacheClient.EXPECT().GetUserFollowees(username, lastFolloweeId, limit).Return([]model.Followee{}, "", false)
	followConnector.EXPECT().GetUserFolloweeIds(username, lastFolloweeId, limit).Return([]string{}, "", errors.New("some error"))

	followees, lastFolloweeId, err := repository.GetUserFollowees(username, lastFolloweeId, limit)

	assert.NotNil(t, err)
	assert.Equal(t, followees, []model.Followee{})
	assert.Equal(t, lastFolloweeId, "")
}

func TestErrorOnGetUserFolloweesFromRepository_WhenReadmodelsConnectorFails(t *testing.T) {
	setUpRepository(t)
	username := "usernameA"
	lastFolloweeId := "followee4"
	limit := 4
	expectedFolloweeIds := []string{"followee5", "followee6", "followee7"}
	expectedLastFolloweeId := "followee4"
	cacheClient.EXPECT().GetUserFollowees(username, lastFolloweeId, limit).Return([]model.Followee{}, "", false)
	followConnector.EXPECT().GetUserFolloweeIds(username, lastFolloweeId, limit).Return(expectedFolloweeIds, expectedLastFolloweeId, nil)
	readmodelsConnector.EXPECT().GetFolloweesMetadata(expectedFolloweeIds).Return([]model.Followee{}, errors.New("some error"))

	followees, lastFolloweeId, err := repository.GetUserFollowees(username, lastFolloweeId, limit)

	assert.NotNil(t, err)
	assert.Equal(t, followees, []model.Followee{})
	assert.Equal(t, lastFolloweeId, "")
}
