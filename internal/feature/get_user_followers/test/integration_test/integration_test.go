package integration_test_get_user_followers

import (
	database "aggregationframework/internal/db"
	"aggregationframework/internal/feature/get_user_followers"
	mock_get_user_followers "aggregationframework/internal/feature/get_user_followers/test/mock"
	model "aggregationframework/internal/model/domain"
	integration_test_arrange "aggregationframework/test/integration_test_common/arrange"
	integration_test_assert "aggregationframework/test/integration_test_common/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

var cache *database.Cache
var controller *get_user_followers.GetUserFollowersController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context
var followerConnector *mock_get_user_followers.MockfollowerConnector
var readmodelsConnector *mock_get_user_followers.MockreadmodelsConnector

func setUp(t *testing.T) {
	// Mocks
	ctrl := gomock.NewController(t)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
	followerConnector = mock_get_user_followers.NewMockfollowerConnector(ctrl)
	readmodelsConnector = mock_get_user_followers.NewMockreadmodelsConnector(ctrl)

	// Real infrastructure and services
	cache = integration_test_arrange.CreateTestCache(t, ginContext)
	repository := get_user_followers.NewGetUserFollowersRepository(cache, followerConnector, readmodelsConnector)
	service := get_user_followers.NewGetUserFollowersService(repository)
	controller = get_user_followers.NewGetUserFollowersController(service)
}

func tearDown() {
	cache.Client.Clean()
}

func TestGetUserFollowers_WhenApiConnectorReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	username := "username1"
	lastFollowerId := "username2"
	limit := 4
	ginContext.Request, _ = http.NewRequest("GET", "/followers", nil)
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	u := url.Values{}
	u.Add("lastFollowerId", lastFollowerId)
	u.Add("limit", strconv.Itoa(limit))
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedFollowerIds := []string{"follower5", "follower6", "follower7"}
	expectedLastFollowerId := "follower7"
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
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"followers":[
				{
					"username": "follower5",
					"fullname": "fullname5"
				},
				{
					"username": "follower6",
					"fullname": "fullname6"
				},
				{
					"username": "follower7",
					"fullname": "fullname7"
				}
			],
			"lastFollowerId":"follower7"
		}
	}`
	followerConnector.EXPECT().GetUserFollowerIds(username, lastFollowerId, limit).Return(expectedFollowerIds, expectedLastFollowerId, nil)
	readmodelsConnector.EXPECT().GetFollowersMetadata(expectedFollowerIds).Return(expectedFollowers, nil)

	controller.GetUserFollowers(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
	integration_test_assert.AssertCachedUserFollowersExists(t, cache, username, lastFollowerId, limit, expectedFollowers)
}

func TestGetUserFollowers_WhenCacheReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	username := "username1"
	lastFollowerId := "username2"
	limit := 4
	populateCache(t, username, lastFollowerId, limit)
	ginContext.Request, _ = http.NewRequest("GET", "/followers", nil)
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	u := url.Values{}
	u.Add("lastFollowerId", lastFollowerId)
	u.Add("limit", strconv.Itoa(limit))
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"followers":[
				{
					"username": "follower5",
					"fullname": "fullname5"
				},
				{
					"username": "follower6",
					"fullname": "fullname6"
				},
				{
					"username": "follower7",
					"fullname": "fullname7"
				}
			],
			"lastFollowerId":"follower7"
		}
	}`

	controller.GetUserFollowers(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
}

func populateCache(t *testing.T, followeeId, lastFollowerId string, limit int) {
	integration_test_arrange.AddCachedFollowersToCache(t, cache, followeeId, lastFollowerId, limit, []model.Follower{
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
	})
}
