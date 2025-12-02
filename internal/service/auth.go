package service

import (
	"context"

	"connectrpc.com/connect"
	v1 "github.com/coreyvan/kid-dictionary/gen/kiddictionary/v1"
	"github.com/coreyvan/kid-dictionary/gen/kiddictionary/v1/kiddictionaryv1connect"
)

var _ kiddictionaryv1connect.AuthServiceHandler = (*AuthService)(nil)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Register(ctx context.Context, req *connect.Request[v1.RegisterRequest]) (*connect.Response[v1.RegisterResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}

func (s *AuthService) Login(ctx context.Context, req *connect.Request[v1.LoginRequest]) (*connect.Response[v1.LoginResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}

func (s *AuthService) RefreshToken(ctx context.Context, req *connect.Request[v1.RefreshTokenRequest]) (*connect.Response[v1.RefreshTokenResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}