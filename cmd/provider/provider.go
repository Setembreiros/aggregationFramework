package provider

import (
	"aggregationframework/internal/api"
)

type Provider struct {
	env string
}

func NewProvider(env string) *Provider {
	return &Provider{
		env: env,
	}
}

func (p *Provider) ProvideApiEndpoint() *api.Api {
	return api.NewApiEndpoint(p.env, p.ProvideApiControllers())
}

func (p *Provider) ProvideApiControllers() []api.Controller {
	return []api.Controller{}
}
