package server

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"gitlab.com/sample_projects/technonext-assessment/internal/service"
)

type Authz struct {
	jwt  service.JWTProvider
	auth *service.AuthService
}

func NewAuthz(jwt service.JWTProvider, auth *service.AuthService) *Authz {
	return &Authz{jwt: jwt, auth: auth}
}

func (a *Authz) Require(ctx context.Context) (jwt.MapClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, unauthorized()
	}
	vs := md.Get("authorization")
	if len(vs) == 0 {
		return nil, unauthorized()
	}
	parts := strings.SplitN(vs[0], " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, unauthorized()
	}

	claims, err := a.jwt.Parse(parts[1])
	if err != nil {
		return nil, unauthorized()
	}

	if a.auth != nil {
		if jti, _ := claims["jti"].(string); jti != "" {
			revoked, _ := a.auth.IsRevoked(ctx, jti)
			if revoked {
				return nil, unauthorized()
			}
		}
	}
	return claims, nil
}

func unauthorized() error {
	return status.Error(codes.Unauthenticated, `{"message":"Unauthorized","type":"error","code":401}`)
}
