package api_connector

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type FollowApiConnector struct {
	*ApiConnector
}

type FollowerIdsContent struct {
	Followers      []string `json:"followers"`
	LastFollowerID string   `json:"lastFollowerId"`
}

type FolloweeIdsContent struct {
	Followees      []string `json:"followees"`
	LastFolloweeID string   `json:"lastFolloweeId"`
}

func NewFollowApiConnector(baseURL string, httpClient *http.Client, context context.Context) *FollowApiConnector {
	return &FollowApiConnector{
		ApiConnector: &ApiConnector{
			baseURL:    baseURL,
			httpClient: httpClient,
			context:    context,
		},
	}
}

func (c *FollowApiConnector) GetUserFollowerIds(username, lastFollowerId string, limit int) ([]string, string, error) {
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

func (c *FollowApiConnector) GetUserFolloweeIds(username, lastFolloweeId string, limit int) ([]string, string, error) {
	uri := fmt.Sprintf("followees/%s?limit=%d", username, limit)
	if lastFolloweeId != "" {
		uri += fmt.Sprintf("&lastFolloweeId=%s", lastFolloweeId)
	}

	result, err := c.SendApiRequest(http.MethodGet, uri)
	if err != nil {
		return []string{}, "", err
	}

	followeeIdsContent, err := deserializeFolloweeIdsContent(result.Content)
	if err != nil {
		return nil, "", NewContentDeserializationError()
	}

	return followeeIdsContent.Followees, followeeIdsContent.LastFolloweeID, nil
}

func deserializeFolloweeIdsContent(content any) (*FolloweeIdsContent, error) {
	jsonBytes, err := json.Marshal(content)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Failed to serialize followeeIds content")
		return nil, NewContentDeserializationError()
	}

	var followeeIdsContent FolloweeIdsContent

	err = json.Unmarshal(jsonBytes, &followeeIdsContent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Failed to deserialize followeeIds content")
		return nil, NewContentDeserializationError()
	}

	return &followeeIdsContent, nil
}
