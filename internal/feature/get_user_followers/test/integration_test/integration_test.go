package integration_test_get_user_followers

import (
	"aggregationframework/internal/feature/get_user_followers"
	mock_get_user_followers "aggregationframework/internal/feature/get_user_followers/test/mock"
	model "aggregationframework/internal/model/domain"
	integration_test_assert "aggregationframework/test/integration_test_common/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

var controller *get_user_followers.GetUserFollowersController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context
var FollowConnector *mock_get_user_followers.MockFollowConnector
var readmodelsConnector *mock_get_user_followers.MockreadmodelsConnector

func setUp(t *testing.T) {
	// Mocks
	ctrl := gomock.NewController(t)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
	FollowConnector = mock_get_user_followers.NewMockFollowConnector(ctrl)
	readmodelsConnector = mock_get_user_followers.NewMockreadmodelsConnector(ctrl)

	// Real infrastructure and services
	repository := get_user_followers.NewGetUserFollowersRepository(FollowConnector, readmodelsConnector)
	service := get_user_followers.NewGetUserFollowersService(repository)
	controller = get_user_followers.NewGetUserFollowersController(service)
}
func TestGetUserFollowers_WhenApiConnectorReturnsSuccess(t *testing.T) {
	setUp(t)
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
	FollowConnector.EXPECT().GetUserFollowerIds(username, lastFollowerId, limit).Return(expectedFollowerIds, expectedLastFollowerId, nil)
	readmodelsConnector.EXPECT().GetFollowersMetadata(expectedFollowerIds).Return(expectedFollowers, nil)

	controller.GetUserFollowers(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
}
