package unit_test_get_user_followees

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"aggregationframework/internal/feature/get_user_followees"
	mock_get_user_followees "aggregationframework/internal/feature/get_user_followees/test/mock"
	model "aggregationframework/internal/model/domain"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var controllerService *mock_get_user_followees.MockService
var controller *get_user_followees.GetUserFolloweesController

func setUpHandler(t *testing.T) {
	setUp(t)
	controllerService = mock_get_user_followees.NewMockService(ctrl)
	controller = get_user_followees.NewGetUserFolloweesController(controllerService)
}

func TestGetUserFolloweesWithController_WhenSuccess(t *testing.T) {
	setUpHandler(t)
	ginContext.Request, _ = http.NewRequest("GET", "/followees", nil)
	expectedUsername := "usernameA"
	expectedLastFolloweeId := "followee4"
	expectedLimit := 4
	ginContext.Params = []gin.Param{{Key: "username", Value: expectedUsername}}
	u := url.Values{}
	u.Add("lastFolloweeId", expectedLastFolloweeId)
	u.Add("limit", strconv.Itoa(expectedLimit))
	ginContext.Request.URL.RawQuery = u.Encode()
	controllerService.EXPECT().GetUserFollowees(expectedUsername, expectedLastFolloweeId, expectedLimit).Return(
		[]model.Followee{
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
		}, "followee7", nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"followees":[
				{
					"username": "followee5",
					"fullname": "fullname5"
				},
				{
					"username": "followee6",
					"fullname": "fullname6"
				},
				{
					"username": "followee7",
					"fullname": "fullname7"
				}
			],
			"lastFolloweeId":"followee7"
		}
	}`

	controller.GetUserFollowees(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestGetUserFolloweesWithController_WhenSuccessWithDefaultPaginationParameters(t *testing.T) {
	setUpHandler(t)
	ginContext.Request, _ = http.NewRequest("GET", "/followees", nil)
	expectedUsername := "usernameA"
	ginContext.Params = []gin.Param{{Key: "username", Value: expectedUsername}}
	expectedDefaultLastFolloweeId := ""
	expectedDefaultLimit := 12
	controllerService.EXPECT().GetUserFollowees("usernameA", expectedDefaultLastFolloweeId, expectedDefaultLimit).Return(
		[]model.Followee{
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
		}, "followee7", nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"followees":[
				{
					"username": "followee5",
					"fullname": "fullname5"
				},
				{
					"username": "followee6",
					"fullname": "fullname6"
				},
				{
					"username": "followee7",
					"fullname": "fullname7"
				}
			],
			"lastFolloweeId":"followee7"
		}
	}`

	controller.GetUserFollowees(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerErrorOnGetUserFolloweesWithController_WhenServiceCallFails(t *testing.T) {
	setUpHandler(t)
	ginContext.Request, _ = http.NewRequest("GET", "/followees", nil)
	expectedUsername := "usernameA"
	ginContext.Params = []gin.Param{{Key: "username", Value: expectedUsername}}
	expectedError := errors.New("some error")
	controllerService.EXPECT().GetUserFollowees("usernameA", "", 12).Return([]model.Followee{}, "", expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content": null
	}`

	controller.GetUserFollowees(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequestErrorOnGetUserPostsWithController_WhenLimitSmallerThanOne(t *testing.T) {
	setUpHandler(t)
	ginContext.Request, _ = http.NewRequest("GET", "/followees", nil)
	username := "usernameA"
	lastFolloweeId := "followee4"
	wrongLimit := 0
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	u := url.Values{}
	u.Add("lastFolloweeId", lastFolloweeId)
	u.Add("limit", strconv.Itoa(wrongLimit))
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedError := "Invalid pagination parameters, limit has to be greater than 0"
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError + `",
		"content":null
	}`

	controller.GetUserFollowees(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}
