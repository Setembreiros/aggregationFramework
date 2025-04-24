package api_connector

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	model "aggregationframework/internal/model/domain"

	"github.com/rs/zerolog/log"
)

type ReadmodelsApiConnector struct {
	*ApiConnector
}

type FollowerMetadataContent struct {
	Followers []model.Follower `json:"followers"`
}

type FolloweeMetadataContent struct {
	Followees []model.Followee `json:"followees"`
}

func NewReadmodelsApiConnector(baseURL string, httpClient *http.Client, context context.Context) *ReadmodelsApiConnector {
	return &ReadmodelsApiConnector{
		ApiConnector: &ApiConnector{
			baseURL:    baseURL,
			httpClient: httpClient,
			context:    context,
		},
	}
}

func (c *ReadmodelsApiConnector) GetFollowersMetadata(followerIds []string) ([]model.Follower, error) {
	if len(followerIds) == 0 {
		return []model.Follower{}, nil
	}

	uri := fmt.Sprintf("followers?")
	for _, followerId := range followerIds {
		if followerId != "" {
			uri += fmt.Sprintf("&followerId=%s", followerId)
		}
	}

	result, err := c.SendApiRequest(http.MethodGet, uri)
	if err != nil {
		return []model.Follower{}, err
	}

	followerMetadtaContent, err := deserializeFollowerMetadataContent(result.Content)
	if err != nil {
		return nil, NewContentDeserializationError()
	}

	return followerMetadtaContent.Followers, nil
}

func deserializeFollowerMetadataContent(content any) (*FollowerMetadataContent, error) {
	jsonBytes, err := json.Marshal(content)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Failed to serialize follower metadata content")
		return nil, NewContentDeserializationError()
	}

	var followerMetadtaContent FollowerMetadataContent

	err = json.Unmarshal(jsonBytes, &followerMetadtaContent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Failed to deserialize follower metadata content")
		return nil, NewContentDeserializationError()
	}

	return &followerMetadtaContent, nil
}

func (c *ReadmodelsApiConnector) GetFolloweesMetadata(followeeIds []string) ([]model.Followee, error) {
	if len(followeeIds) == 0 {
		return []model.Followee{}, nil
	}

	uri := fmt.Sprintf("followees?")
	for _, followeeId := range followeeIds {
		if followeeId != "" {
			uri += fmt.Sprintf("&followeeId=%s", followeeId)
		}
	}

	result, err := c.SendApiRequest(http.MethodGet, uri)
	if err != nil {
		return []model.Followee{}, err
	}

	followeeMetadtaContent, err := deserializeFolloweeMetadataContent(result.Content)
	if err != nil {
		return nil, NewContentDeserializationError()
	}

	return followeeMetadtaContent.Followees, nil
}

func deserializeFolloweeMetadataContent(content any) (*FolloweeMetadataContent, error) {
	jsonBytes, err := json.Marshal(content)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Failed to serialize followee metadata content")
		return nil, NewContentDeserializationError()
	}

	var followeeMetadtaContent FolloweeMetadataContent

	err = json.Unmarshal(jsonBytes, &followeeMetadtaContent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Failed to deserialize followee metadata content")
		return nil, NewContentDeserializationError()
	}

	return &followeeMetadtaContent, nil
}
