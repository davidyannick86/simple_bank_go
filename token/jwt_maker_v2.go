package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySizeV2 = 32

type JWTMakerV2 struct {
	secretKey string
}

func NewJWTMakerV2(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySizeV2 {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySizeV2)
	}
	return &JWTMakerV2{secretKey: secretKey}, nil
}

func (maker *JWTMakerV2) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	// Create the JWT token using the payload and secret key
	claims := jwt.MapClaims{
		"username": payload.Username,
		"exp":      payload.ExpiredAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (maker *JWTMakerV2) VerifyToken(tokenString string) (*Payload, error) {
	// Parse and validate the token
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !jwtToken.Valid {
		return nil, ErrInvalidToken
	}

	// Extract claims
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	expClaim, ok := claims["exp"].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}

	expTime := time.Unix(int64(expClaim), 0)
	if time.Now().After(expTime) {
		return nil, ErrExpiredToken
	}

	payload := &Payload{
		Username:  username,
		ExpiredAt: expTime,
	}
	return payload, nil
}
