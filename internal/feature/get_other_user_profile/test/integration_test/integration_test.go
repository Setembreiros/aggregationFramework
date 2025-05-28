package integration_test_other_user_profile

import (
	"aggregationframework/internal/feature/get_other_user_profile"
	mock_get_other_user_profile "aggregationframework/internal/feature/get_other_user_profile/test/mock"
	model "aggregationframework/internal/model/domain"
	integration_test_assert "aggregationframework/test/integration_test_common/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

var controller *get_other_user_profile.GetOtherUserProfileController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context
var FollowConnector *mock_get_other_user_profile.MockfollowConnector
var readmodelsConnector *mock_get_other_user_profile.MockreadmodelsConnector

func setUp(t *testing.T) {
	// Mocks
	ctrl := gomock.NewController(t)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
	FollowConnector = mock_get_other_user_profile.NewMockfollowConnector(ctrl)
	readmodelsConnector = mock_get_other_user_profile.NewMockreadmodelsConnector(ctrl)

	// Real infrastructure and services
	repository := get_other_user_profile.NewGetOtherUserProfileRepository(FollowConnector, readmodelsConnector)
	service := get_other_user_profile.NewGetOtherUserProfileService(repository)
	controller = get_other_user_profile.NewGetOtherUserProfileController(service)
}

func TestGetOtherUserProfile_WhenApiConnectorReturnsSuccess(t *testing.T) {
	setUp(t)
	expectedUsername := "usernameA"
	expectedCurrentUsername := "usernameB"
	ginContext.Request, _ = http.NewRequest("GET", "/userprofile", nil)
	ginContext.Params = []gin.Param{{Key: "username", Value: expectedUsername}}
	u := url.Values{}
	u.Add("currentUsername", expectedCurrentUsername)
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedUserProfile := &model.UserProfile{
		Username: expectedUsername,
	}
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
	FollowConnector.EXPECT().CheckIfRelationshipExists(expectedUsername, expectedCurrentUsername).Return(true, nil)
	readmodelsConnector.EXPECT().GetUserProfile(expectedUsername).Return(expectedUserProfile, nil)

	controller.GetOtherUserProfile(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
}
