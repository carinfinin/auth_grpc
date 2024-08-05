package app

import (
	"auth/cmd/app/grpcapp"
	"auth/internal/services/auth"
	"auth/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// init storage

	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	// init service auth
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
