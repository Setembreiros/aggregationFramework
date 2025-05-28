package unit_test_get_other_user_profile

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"aggregationframework/internal/feature/get_other_user_profile"
	mock_get_other_user_profile "aggregationframework/internal/feature/get_other_user_profile/test/mock"
	model "aggregationframework/internal/model/domain"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var controllerService *mock_get_other_user_profile.MockService
var controller *get_other_user_profile.GetOtherUserProfileController

func setUpHandler(t *testing.T) {
	setUp(t)
	controllerService = mock_get_other_user_profile.NewMockService(ctrl)
	controller = get_other_user_profile.NewGetOtherUserProfileController(controllerService)
}

func TestGetOtherUserProfileWithController_WhenSuccess(t *testing.T) {
	setUpHandler(t)
	ginContext.Request, _ = http.NewRequest("GET", "/userprofile", nil)
	expectedUsername := "usernameA"
	expectedCurrentUsername := "usernameB"
	ginContext.Params = []gin.Param{{Key: "username", Value: expectedUsername}}
	u := url.Values{}
	u.Add("currentUsername", expectedCurrentUsername)
	ginContext.Request.URL.RawQuery = u.Encode()
	controllerService.EXPECT().GetOtherUserProfile(expectedUsername, expectedCurrentUsername).Return(
		&model.UserProfile{
			Username:                expectedUsername,
			IsFollowedByCurrentUser: true,
		}, nil)
	expectedBodyResponse := `
	{
		"error":false,
		"message":"200OK",
		"content":{
			"userProfile":{
				"username":"usernameA",
				"name":"",
				"bio":"",
				"link":"",
				"followersAmount":0,
				"followeesAmount":0,
				"isFollowedByCurrentUser":true,
				"postsAmount":0
			}
		}
	}`

	controller.GetOtherUserProfile(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerErrorOnGetOtherUserProfileWithController_WhenServiceCallFails(t *testing.T) {
	setUpHandler(t)
	ginContext.Request, _ = http.NewRequest("GET", "/userprofile", nil)
	expectedUsername := "usernameA"
	expectedCurrentUsername := "usernameB"
	ginContext.Params = []gin.Param{{Key: "username", Value: expectedUsername}}
	u := url.Values{}
	u.Add("currentUsername", expectedCurrentUsername)
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedError := errors.New("some error")
	controllerService.EXPECT().GetOtherUserProfile(expectedUsername, expectedCurrentUsername).Return(nil, expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content": null
	}`

	controller.GetOtherUserProfile(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequestErrorOnGetOtherUserProfileWithController_WhenCurrentUserIsEmpty(t *testing.T) {
	setUpHandler(t)
	ginContext.Request, _ = http.NewRequest("GET", "/userprofile", nil)
	expectedUsername := "usernameA"
	ginContext.Params = []gin.Param{{Key: "username", Value: expectedUsername}}
	u := url.Values{}
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedError := "Invalid query parameters, currentUsername has to be not empty"
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError + `",
		"content":null
	}`

	controller.GetOtherUserProfile(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}
