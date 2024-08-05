package grpcapp

import (
	"auth/internal/grpc/auth"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, app int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (id int64, err error)
	IsAdmin(ctx context.Context, id int64) (bool, error)
}

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, authService Auth, port int) *App {

	gRPCServer := grpc.NewServer()

	auth.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {

	const op = "grpcapp.Run"

	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	log.Info("stopping grpc server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()

}
