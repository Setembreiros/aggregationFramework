package unit_test_get_user_followers

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"aggregationframework/internal/feature/get_user_followers"
	mock_get_user_followers "aggregationframework/internal/feature/get_user_followers/test/mock"
	model "aggregationframework/internal/model/domain"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var controllerService *mock_get_user_followers.MockService
var controller *get_user_followers.GetUserFollowersController

func setUpHandler(t *testing.T) {
	setUp(t)
	controllerService = mock_get_user_followers.NewMockService(ctrl)
	controller = get_user_followers.NewGetUserFollowersController(controllerService)
}

func TestGetUserFollowersWithController_WhenSuccess(t *testing.T) {
	setUpHandler(t)
	ginContext.Request, _ = http.NewRequest("GET", "/followers", nil)
	expectedUsername := "usernameA"
	expectedLastFollowerId := "follower4"
	expectedLimit := 4
	ginContext.Params = []gin.Param{{Key: "username", Value: expectedUsername}}
	u := url.Values{}
	u.Add("lastFollowerId", expectedLastFollowerId)
	u.Add("limit", strconv.Itoa(expectedLimit))
	ginContext.Request.URL.RawQuery = u.Encode()
	controllerService.EXPECT().GetUserFollowers(expectedUsername, expectedLastFollowerId, expectedLimit).Return(
		[]model.Follower{
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
		}, "follower7", nil)
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

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestGetUserFollowersWithController_WhenSuccessWithDefaultPaginationParameters(t *testing.T) {
	setUpHandler(t)
	ginContext.Request, _ = http.NewRequest("GET", "/followers", nil)
	expectedUsername := "usernameA"
	ginContext.Params = []gin.Param{{Key: "username", Value: expectedUsername}}
	expectedDefaultLastFollowerId := ""
	expectedDefaultLimit := 12
	controllerService.EXPECT().GetUserFollowers("usernameA", expectedDefaultLastFollowerId, expectedDefaultLimit).Return(
		[]model.Follower{
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
		}, "follower7", nil)
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

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerErrorOnGetUserFollowersWithController_WhenServiceCallFails(t *testing.T) {
	setUpHandler(t)
	ginContext.Request, _ = http.NewRequest("GET", "/followers", nil)
	expectedUsername := "usernameA"
	ginContext.Params = []gin.Param{{Key: "username", Value: expectedUsername}}
	expectedError := errors.New("some error")
	controllerService.EXPECT().GetUserFollowers("usernameA", "", 12).Return([]model.Follower{}, "", expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content": null
	}`

	controller.GetUserFollowers(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequestErrorOnGetUserPostsWithController_WhenLimitSmallerThanOne(t *testing.T) {
	setUpHandler(t)
	ginContext.Request, _ = http.NewRequest("GET", "/followers", nil)
	username := "usernameA"
	lastFollowerId := "follower4"
	wrongLimit := 0
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	u := url.Values{}
	u.Add("lastFollowerId", lastFollowerId)
	u.Add("limit", strconv.Itoa(wrongLimit))
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedError := "Invalid pagination parameters, limit has to be greater than 0"
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError + `",
		"content":null
	}`

	controller.GetUserFollowers(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}
