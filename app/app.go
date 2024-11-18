package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

func New() *App {
	return &App{
		router: loadRoutes(),
		rdb:    redis.NewClient(&redis.Options{}),
	}
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":8080",
		Handler: a.router,
	}

	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("erro ao conectar o redis: %w", err)
	}

	fmt.Println("Iniciando server")

	err = server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("erro iniciando server: %w", err)
	}

	return nil
}
