package auth

import (
	"context"
	"github.com/carinfinin/auth_proto/gen/go/auth"
	"google.golang.org/grpc"
)

type serverAPI struct {
	auth_golang.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	auth_golang.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, in *auth_golang.LoginRequest) (*auth_golang.LoginResponse, error) {
	panic("implement me Login")
}

func (s *serverAPI) Register(ctx context.Context, in *auth_golang.LoginRequest) (*auth_golang.LoginResponse, error) {
	panic("implement me Register")
}

func (s *serverAPI) IsAdmin(ctx context.Context, in *auth_golang.LoginRequest) (*auth_golang.LoginResponse, error) {
	panic("implement me IsAdmin")
}
