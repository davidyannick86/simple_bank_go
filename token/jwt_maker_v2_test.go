// filepath: token/jwt_maker_v2_test.go
package token

import (
	"testing"
	"time"

	"github.com/davidyannick86/simple-bank/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTMakerV2(t *testing.T) {
	maker, err := NewJWTMakerV2(utils.RandomString(minSecretKeySizeV2))
	require.NoError(t, err)

	username := utils.RandomOwner()
	duration := time.Minute

	tokenString, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	payload, err := maker.VerifyToken(tokenString)
	require.NoError(t, err)
	require.NotNil(t, payload)

	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, time.Now().Add(duration), payload.ExpiredAt, time.Second)
}

func TestExpiredJWTTokenV2(t *testing.T) {
	maker, err := NewJWTMakerV2(utils.RandomString(minSecretKeySizeV2))
	require.NoError(t, err)

	tokenString, err := maker.CreateToken(utils.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	payload, err := maker.VerifyToken(tokenString)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNoneV2(t *testing.T) {
	username := utils.RandomOwner()
	expiration := time.Now().Add(time.Minute)
	claims := jwt.MapClaims{
		"username": username,
		"exp":      expiration.Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMakerV2(utils.RandomString(minSecretKeySizeV2))
	require.NoError(t, err)

	payload, err := maker.VerifyToken(tokenString)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
