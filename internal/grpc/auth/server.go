package auth

import (
	"context"
	authV1 "github.com/carinfinin/auth_proto/gen/go/auth"
	"google.golang.org/grpc"
)

type serverAPI struct {
	authV1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	authV1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, in *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	panic("implement me Login")
}

//func (s *serverAPI) Register(ctx context.Context, in *authV1.LoginRequest) (*authV1.LoginResponse, error) {
//	panic("implement me Register")
//}
//
//func (s *serverAPI) IsAdmin(ctx context.Context, in *authV1.LoginRequest) (*authV1.LoginResponse, error) {
//	panic("implement me IsAdmin")
//}
//
//func (s *serverAPI) mustEmbedUnimplementedAuthServer(ctx context.Context, in *authV1.LoginRequest) (*authV1.LoginResponse, error) {
//	panic("implement me IsAdmin")
//}
