package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gitlab.com/sample_projects/technonext-assessment/internal/repo"
	"time"
)

type JWTProvider interface {
	GenerateAccess(username string, ttl time.Duration) (signed string, jti string, exp int64, err error)
	Parse(token string) (jwt.MapClaims, error)
}

type AuthService struct {
	repo       *repo.AuthRepo
	jwt        JWTProvider
	redis      *redis.Client
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(ar *repo.AuthRepo, jwt JWTProvider, redis *redis.Client, accessTTL, refreshTTL time.Duration) *AuthService {
	return &AuthService{repo: ar, jwt: jwt, redis: redis, accessTTL: accessTTL, refreshTTL: refreshTTL}
}

func (s *AuthService) Login(ctx context.Context, user, pass string) (token string, exp int64, refresh string, err error) {
	u, err := s.repo.FindUserByUsername(user, pass)
	if err != nil {
		return "", 0, "", errors.New("bad-credentials")
	}
	tok, _, expUnix, err := s.jwt.GenerateAccess(u.Username, s.accessTTL)
	if err != nil {
		return "", 0, "", err
	}
	ref, _, _, _ := s.jwt.GenerateAccess(u.Username, s.refreshTTL)
	return tok, expUnix, ref, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	if s.redis == nil {
		return nil
	}
	claims, err := s.jwt.Parse(token)
	if err != nil {
		return err
	}
	jti, _ := claims["jti"].(string)
	exp, _ := claims["exp"].(float64)
	if jti == "" || exp == 0 {
		return errors.New("invalid token")
	}
	ttl := time.Until(time.Unix(int64(exp), 0))
	if ttl <= 0 {
		return nil
	}
	return s.redis.Set(ctx, "bl:"+jti, "1", ttl).Err()
}

func (s *AuthService) IsRevoked(ctx context.Context, jti string) (bool, error) {
	if s.redis == nil {
		return false, nil
	}
	n, err := s.redis.Exists(ctx, "bl:"+jti).Result()
	return n == 1, err
}
