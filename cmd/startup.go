package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"aggregationframework/cmd/provider"
	"aggregationframework/internal/api"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type App struct {
	Ctx              context.Context
	Cancel           context.CancelFunc
	Env              string
	configuringTasks sync.WaitGroup
	runningTasks     sync.WaitGroup
}

func (app *App) Startup() {
	app.configuringLog()

	log.Info().Msgf("Starting Aggregationframework service in [%s] enviroment...\n", app.Env)

	provider := provider.NewProvider(app.Env)
	httpClient := provider.ProvideHttpClient()
	FollowConnector := provider.ProvideFollowApiConnector(httpClient, app.Ctx)
	readmodelsConnector := provider.ProvideReadmodelsApiConnector(httpClient, app.Ctx)
	apiEnpoint := provider.ProvideApiEndpoint(FollowConnector, readmodelsConnector)

	app.runServerTasks(apiEnpoint)
}

func (app *App) configuringLog() {
	if app.Env == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Logger = log.With().Caller().Logger()
}

func (app *App) runServerTasks(apiEnpoint *api.Api) {
	app.runningTasks.Add(1)
	go app.runApiEndpoint(apiEnpoint)

	blockForever()

	app.shutdown()
}

func (app *App) runApiEndpoint(apiEnpoint *api.Api) {
	defer app.runningTasks.Done()

	err := apiEnpoint.Run(app.Ctx)
	if err != nil {
		log.Panic().Err(err).Msg("Closing AggregationFramework Api failed")
	}
	log.Info().Msg("AggregationFramework Api stopped")
}

func blockForever() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh
}

func (app *App) shutdown() {
	app.Cancel()
	log.Info().Msg("Shutting down AggregationFramework Service...")
	app.runningTasks.Wait()
	log.Info().Msg("AggregationFramework Service stopped")
}
