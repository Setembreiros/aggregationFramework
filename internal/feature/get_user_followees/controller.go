package get_user_followees

import (
	"aggregationframework/internal/api"
	"strconv"

	model "aggregationframework/internal/model/domain"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=controller.go -destination=test/mock/controller.go

type GetUserFolloweesController struct {
	service Service
}

type Service interface {
	GetUserFollowees(username string, lastPostId string, limit int) ([]model.Followee, string, error)
}

type GetUserFolloweesResponse struct {
	Followees      []model.Followee `json:"followees"`
	LastFolloweeId string           `json:"lastFolloweeId"`
}

func NewGetUserFolloweesController(service Service) *GetUserFolloweesController {
	return &GetUserFolloweesController{
		service: service,
	}
}

func (controller *GetUserFolloweesController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/followees/:username", controller.GetUserFollowees)
}

func (controller *GetUserFolloweesController) GetUserFollowees(c *gin.Context) {
	log.Info().Msg("Handling Request GET GetUserFollowees")
	username, lastFolloweeId, limit := getQueryParameters(c)
	if username == "" {
		return
	}

	followees, lastFolloweeId, err := controller.service.GetUserFollowees(username, lastFolloweeId, limit)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, &GetUserFolloweesResponse{
		Followees:      followees,
		LastFolloweeId: lastFolloweeId,
	})
}

func getQueryParameters(c *gin.Context) (string, string, int) {
	username := c.Param("username")
	if username == "" {
		api.SendBadRequest(c, "Missing username parameter")
		return "", "", 0
	}
	lastFolloweeId := c.DefaultQuery("lastFolloweeId", "")
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "12"))
	if err != nil || limit <= 0 {
		api.SendBadRequest(c, "Invalid pagination parameters, limit has to be greater than 0")
		return "", "", 0
	}

	return username, lastFolloweeId, limit
}
