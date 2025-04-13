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

func NewReadmodelsApiConnector(baseURL string, httpClient *http.Client, context context.Context) *ReadmodelsApiConnector {
	return &ReadmodelsApiConnector{
		ApiConnector: &ApiConnector{
			baseURL:    baseURL,
			httpClient: httpClient,
			context:    context,
		},
	}
}

func (c *ReadmodelsApiConnector) GetFollowerMetadatas(followerIds []string) ([]model.Follower, error) {
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
