package api_connector

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ApiConnector struct {
	baseURL    string
	httpClient *http.Client
	context    context.Context
}

type BaseResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Content any    `json:"content"`
}

func (c ApiConnector) SendApiRequest(method, uri string) (BaseResponse, error) {
	url := c.baseURL + uri
	req, err := http.NewRequestWithContext(c.context, string(method), url, nil)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error creating %s %s request", method, url)
		return BaseResponse{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("HTTP %s %s call failed", method, url)
		return BaseResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Stack().Err(err).Msgf("%s %s returned a StatusCode: %d ", method, url, resp.StatusCode)
		return BaseResponse{}, NewBadStatusCodeResponseError(resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error reading %s %s body response", method, url)
		return BaseResponse{}, err
	}

	var result BaseResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Error().Stack().Err(err).Msgf("Failed to deserialize %s %s base response data", method, url)
		return BaseResponse{}, err
	}

	if result.Error {
		log.Warn().Stack().Msgf("Call to %s %s responsed with an error", method, url)
		return BaseResponse{}, fmt.Errorf(result.Message)
	}

	return result, nil
}
