package token

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtProvider struct {
	cfgToken ConfigToken
}

func NewJWTProvider(cfgToken ConfigToken) ITokenProvider {
	return &jwtProvider{cfgToken: cfgToken}
}

func (j *jwtProvider) GenerateAccessToken(userID uuid.UUID, role string) (string, error) {
	claims := &UserClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfgToken.AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Go-CollabSpace",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.cfgToken.AccessTokenSecret))
}

func (j *jwtProvider) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32) // 32 bytes = 256 bits, use a secure random generator in production
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (j *jwtProvider) ValidateAccessToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, j.keyFunc)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrInvalidToken
}

func (j *jwtProvider) keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(j.cfgToken.AccessTokenSecret), nil
}
