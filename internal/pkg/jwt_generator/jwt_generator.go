package jwt_generator

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type IJWTGenerator interface {
	Generate(id string, expire time.Duration) (string, error)
}

type JWTGenerator struct {
	privateKey *rsa.PrivateKey
	issuer     string
	nowFunc    func() time.Time
}

func NewJWTGenerator(issuer string, privateKey string) *JWTGenerator {
	pk, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		panic("cannot parse private key")
	}

	return &JWTGenerator{
		privateKey: pk,
		issuer:     issuer,
		nowFunc:    time.Now,
	}
}

func (j *JWTGenerator) Generate(id string, expire time.Duration) (string, error) {
	now := j.nowFunc()
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.RegisteredClaims{
		Issuer:    j.issuer,
		Subject:   id,
		ExpiresAt: jwt.NewNumericDate(now.Add(expire)),
		IssuedAt:  jwt.NewNumericDate(now),
	})
	return token.SignedString(j.privateKey)
}
