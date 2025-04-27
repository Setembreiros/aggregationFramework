package integration_test_get_user_followees

import (
	"aggregationframework/internal/feature/get_user_followees"
	mock_get_user_followees "aggregationframework/internal/feature/get_user_followees/test/mock"
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

var controller *get_user_followees.GetUserFolloweesController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context
var FollowConnector *mock_get_user_followees.MockfollowConnector
var readmodelsConnector *mock_get_user_followees.MockreadmodelsConnector

func setUp(t *testing.T) {
	// Mocks
	ctrl := gomock.NewController(t)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
	FollowConnector = mock_get_user_followees.NewMockfollowConnector(ctrl)
	readmodelsConnector = mock_get_user_followees.NewMockreadmodelsConnector(ctrl)

	// Real infrastructure and services
	repository := get_user_followees.NewGetUserFolloweesRepository(FollowConnector, readmodelsConnector)
	service := get_user_followees.NewGetUserFolloweesService(repository)
	controller = get_user_followees.NewGetUserFolloweesController(service)
}

func TestGetUserFollowees_WhenApiConnectorReturnsSuccess(t *testing.T) {
	setUp(t)
	username := "username1"
	lastFolloweeId := "username2"
	limit := 4
	ginContext.Request, _ = http.NewRequest("GET", "/followees", nil)
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	u := url.Values{}
	u.Add("lastFolloweeId", lastFolloweeId)
	u.Add("limit", strconv.Itoa(limit))
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedFolloweeIds := []string{"followee5", "followee6", "followee7"}
	expectedLastFolloweeId := "followee7"
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
	FollowConnector.EXPECT().GetUserFolloweeIds(username, lastFolloweeId, limit).Return(expectedFolloweeIds, expectedLastFolloweeId, nil)
	readmodelsConnector.EXPECT().GetFolloweesMetadata(expectedFolloweeIds).Return(expectedFollowees, nil)

	controller.GetUserFollowees(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
}
