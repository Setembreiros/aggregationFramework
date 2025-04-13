package api_connector

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type FollowerApiConnector struct {
	*ApiConnector
}

type FollowerIdsContent struct {
	Followers      []string `json:"followers"`
	LastFollowerID string   `json:"lastFollowerId"`
}

func NewFollowerApiConnector(baseURL string, httpClient *http.Client, context context.Context) *FollowerApiConnector {
	return &FollowerApiConnector{
		ApiConnector: &ApiConnector{
			baseURL:    baseURL,
			httpClient: httpClient,
			context:    context,
		},
	}
}

func (c *FollowerApiConnector) GetUserFollowerIds(username, lastFollowerId string, limit int) ([]string, string, error) {
	uri := fmt.Sprintf("followers/%s?limit=%d", username, limit)
	if lastFollowerId != "" {
		uri += fmt.Sprintf("&lastFollowerId=%s", lastFollowerId)
	}

	result, err := c.SendApiRequest(http.MethodGet, uri)
	if err != nil {
		return []string{}, "", err
	}

	followerIdsContent, err := deserializeFollowerIdsContent(result.Content)
	if err != nil {
		return nil, "", NewContentDeserializationError()
	}

	return followerIdsContent.Followers, followerIdsContent.LastFollowerID, nil
}

func deserializeFollowerIdsContent(content any) (*FollowerIdsContent, error) {
	jsonBytes, err := json.Marshal(content)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Failed to serialize followerIds content")
		return nil, NewContentDeserializationError()
	}

	var followerIdsContent FollowerIdsContent

	err = json.Unmarshal(jsonBytes, &followerIdsContent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Failed to deserialize followerIds content")
		return nil, NewContentDeserializationError()
	}

	return &followerIdsContent, nil
}
