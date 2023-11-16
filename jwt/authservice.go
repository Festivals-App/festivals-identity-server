package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/Festivals-App/festivals-identity-server/server/config"
	"github.com/Festivals-App/festivals-identity-server/server/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

// Authentication interface lists the methods that our authentication service should implement
type Authentication interface {
	Authenticate(reqUser *model.User, user *model.User) bool
	GenerateAccessToken(user *model.User) (string, error)
	GenerateRefreshToken(user *model.User) (string, error)
	GenerateCustomKey(userID string, password string) string
	ValidateAccessToken(token string) (string, error)
	ValidateRefreshToken(token string) (string, string, error)
}

type AuthService struct {
	config *config.Config
}

type AccessTokenCustomClaims struct {
	UserID  string
	KeyType string
	jwt.StandardClaims
}

// GenerateAccessToken generates a new access token for the given user
func (auth *AuthService) GenerateAccessToken(user *model.User) (string, error) {

	userID := string(user.ID)
	tokenType := "access"

	claims := AccessTokenCustomClaims{
		userID,
		tokenType,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(auth.config.JwtExpiration)).Unix(),
			Issuer:    "festivals.identity.server",
		},
	}

	signBytes, err := os.ReadFile(auth.config.AccessTokenPrivateKeyPath)
	if err != nil {
		log.Error().Err(err).Msg("unable to read private key for access token")
		return "", errors.New("could not generate access token. please try again later")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse private key for access token")
		return "", errors.New("could not generate access token. please try again later")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(signKey)
}

// RefreshTokenCustomClaims specifies the claims for refresh token
type RefreshTokenCustomClaims struct {
	UserID    string
	CustomKey string
	KeyType   string
	jwt.StandardClaims
}

// GenerateRefreshToken generates a new refresh token for the given user
func (auth *AuthService) GenerateRefreshToken(user *model.User) (string, error) {

	cusKey := auth.GenerateCustomKey(string(user.ID), user.TokenHash)
	tokenType := "refresh"

	claims := RefreshTokenCustomClaims{
		string(user.ID),
		cusKey,
		tokenType,
		jwt.StandardClaims{
			Issuer: "festivals.identity.server",
		},
	}

	signBytes, err := os.ReadFile(auth.config.RefreshTokenPrivateKeyPath)
	if err != nil {
		log.Error().Err(err).Msg("unable to read private key for refresh token")
		return "", errors.New("could not generate refresh token. please try again later")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse private key for refresh token")
		return "", errors.New("could not generate refresh token. please try again later")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(signKey)
}

// GenerateCustomKey creates a new key for our jwt payload
// the key is a hashed combination of the userID and user tokenhash
func (auth *AuthService) GenerateCustomKey(userID string, tokenHash string) string {

	// data := userID + tokenHash
	h := hmac.New(sha256.New, []byte(tokenHash))
	h.Write([]byte(userID))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
