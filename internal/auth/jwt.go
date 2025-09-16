package auth

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte(getenv("JWT_SECRET", "technonext_secret"))

type Provider struct{}

func NewProvider() *Provider { return &Provider{} }

func (p *Provider) GenerateAccess(username string, ttl time.Duration) (signed string, jti string, exp int64, err error) {
	jti = randJTI()
	exp = time.Now().Add(ttl).Unix()
	claims := jwt.MapClaims{
		"sub": username,
		"exp": exp,
		"iat": time.Now().Unix(),
		"jti": jti,
		"typ": "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err = token.SignedString(secret)
	return
}

func (p *Provider) Parse(tok string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(tok, func(t *jwt.Token) (interface{}, error) { return secret, nil })
	if err != nil || !t.Valid {
		return nil, err
	}
	if c, ok := t.Claims.(jwt.MapClaims); ok {
		return c, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

func randJTI() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
