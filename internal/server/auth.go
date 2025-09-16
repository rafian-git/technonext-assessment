package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
	"time"

	ordersv1 "gitlab.com/sample_projects/technonext-assessment/gen/orders/v1"
	"gitlab.com/sample_projects/technonext-assessment/internal/service"
)

type AuthServer struct {
	ordersv1.UnimplementedAuthServiceServer
	svc *service.AuthService
	jwt service.JWTProvider
}

func NewAuthServer(svc *service.AuthService, jwt service.JWTProvider) *AuthServer {
	return &AuthServer{svc: svc, jwt: jwt}
}

func (s *AuthServer) Login(ctx context.Context, req *ordersv1.LoginRequest) (*ordersv1.LoginResponse, error) {
	tok, exp, refresh, err := s.svc.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, statusInvalidCredentials()
	}
	return &ordersv1.LoginResponse{
		TokenType:    "Bearer",
		ExpiresIn:    exp - time.Now().Unix(),
		AccessToken:  tok,
		RefreshToken: refresh,
	}, nil
}

func (s *AuthServer) Logout(ctx context.Context, _ *ordersv1.LogoutRequest) (*ordersv1.LogoutResponse, error) {

	md, _ := metadata.FromIncomingContext(ctx)
	var tok string
	if vs := md.Get("authorization"); len(vs) > 0 {
		parts := strings.SplitN(vs[0], " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			tok = parts[1]
		}
	}
	if tok != "" {
		err := s.svc.Logout(ctx, tok)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, status.Error(codes.Unauthenticated, `{"message":"Unauthorized","type":"error","code":401}`)
	}
	return &ordersv1.LogoutResponse{Message: "Successfully logged out", Type: "success", Code: 200}, nil
}

func statusInvalidCredentials() error {
	return status.Error(codes.InvalidArgument, `{"message":"The user credentials were incorrect.","type":"error","code":400}`)
}
