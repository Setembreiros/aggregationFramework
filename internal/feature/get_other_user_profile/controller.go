package get_other_user_profile

import (
	"aggregationframework/internal/api"

	model "aggregationframework/internal/model/domain"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=controller.go -destination=test/mock/controller.go

type GetOtherUserProfileController struct {
	service Service
}

type Service interface {
	GetOtherUserProfile(username string, currentUsername string) (*model.UserProfile, error)
}

type GetOtherUserProfileResponse struct {
	UserProfile *model.UserProfile `json:"userProfile"`
}

func NewGetOtherUserProfileController(service Service) *GetOtherUserProfileController {
	return &GetOtherUserProfileController{
		service: service,
	}
}

func (controller *GetOtherUserProfileController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/userprofile/:username", controller.GetOtherUserProfile)
}

func (controller *GetOtherUserProfileController) GetOtherUserProfile(c *gin.Context) {
	log.Info().Msg("Handling Request GET GetOtherUserProfile")
	username, currentUsername := getQueryParameters(c)
	if username == "" {
		return
	}

	userprofile, err := controller.service.GetOtherUserProfile(username, currentUsername)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, &GetOtherUserProfileResponse{
		UserProfile: userprofile,
	})
}

func getQueryParameters(c *gin.Context) (string, string) {
	username := c.Param("username")
	if username == "" {
		api.SendBadRequest(c, "Missing username parameter")
		return "", ""
	}
	currentUsername := c.DefaultQuery("currentUsername", "")
	if currentUsername == "" {
		api.SendBadRequest(c, "Invalid query parameters, currentUsername has to be not empty")
		return "", ""
	}

	return username, currentUsername
}
