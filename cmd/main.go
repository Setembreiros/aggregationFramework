package main

import (
	"context"
	"os"
	"strings"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	env := strings.TrimSpace(os.Getenv("ENVIRONMENT"))

	app := &App{
		Ctx:    ctx,
		Cancel: cancel,
		Env:    env,
	}

	app.Startup()
}
