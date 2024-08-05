package auth

import (
	"auth/internal/services/auth"
	"auth/internal/storage"
	"context"
	"errors"
	authV1 "github.com/carinfinin/auth_proto/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, app int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (id int64, err error)
	IsAdmin(ctx context.Context, id int64) (bool, error)
}

type serverAPI struct {
	authV1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	authV1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, in *authV1.LoginRequest) (*authV1.LoginResponse, error) {

	if in.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if in.GetApp() == 0 {
		return nil, status.Error(codes.InvalidArgument, "app is required")
	}

	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword(), int(in.GetApp()))

	if err != nil {
		if errors.Is(err, auth.ErrorInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "Invalid Argument error")
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &authV1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, in *authV1.RegisterRequest) (*authV1.RegisterResponse, error) {
	if in.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	id, err := s.auth.RegisterNewUser(ctx, in.GetEmail(), in.GetPassword())

	if err != nil {
		if errors.Is(err, auth.ErrorUserExists) {
			return nil, status.Error(codes.AlreadyExists, "User already Exists ")
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &authV1.RegisterResponse{Id: id}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, in *authV1.IsAdminRequest) (*authV1.IsAdminResponse, error) {
	if in.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, in.GetId())

	if err != nil {
		if errors.Is(err, storage.ErrorUserNotFound) {
			return nil, status.Error(codes.AlreadyExists, "User not found ")
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &authV1.IsAdminResponse{IsAdmin: isAdmin}, nil

}
func (s *serverAPI) mustEmbedUnimplementedAuthServer() {}
