package provider

import (
	"aggregationframework/infrastructure/api_connector"
	"aggregationframework/internal/api"
	"aggregationframework/internal/feature/get_user_followees"
	"aggregationframework/internal/feature/get_user_followers"
	"context"
	"net"
	"net/http"
	"time"
)

type Provider struct {
	env string
}

func NewProvider(env string) *Provider {
	return &Provider{
		env: env,
	}
}

func (p *Provider) ProvideHttpClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
		},
	}
}

func (p *Provider) ProvideFollowApiConnector(httpClient *http.Client, context context.Context) *api_connector.FollowApiConnector {
	baseURL := "http://localhost:7777/" + p.env + "/followservice/"
	return api_connector.NewFollowApiConnector(baseURL, httpClient, context)
}

func (p *Provider) ProvideReadmodelsApiConnector(httpClient *http.Client, context context.Context) *api_connector.ReadmodelsApiConnector {
	baseURL := "http://localhost:5555/" + p.env + "/readmodels/"
	return api_connector.NewReadmodelsApiConnector(baseURL, httpClient, context)
}

func (p *Provider) ProvideApiControllers(followConnector *api_connector.FollowApiConnector, readmodelsConnector *api_connector.ReadmodelsApiConnector) []api.Controller {
	return []api.Controller{
		get_user_followers.NewGetUserFollowersController(get_user_followers.NewGetUserFollowersService(get_user_followers.NewGetUserFollowersRepository(followConnector, readmodelsConnector))),
		get_user_followees.NewGetUserFolloweesController(get_user_followees.NewGetUserFolloweesService(get_user_followees.NewGetUserFolloweesRepository(followConnector, readmodelsConnector))),
	}
}
