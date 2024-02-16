package token

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

// ValidateAccessToken parses and validates the given access token
// returns the custom claim present in the token payload
func (auth *AuthService) ValidateAccessToken(tokenString string) (*UserClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			log.Error().Msg("Unexpected signing method in auth token")
			return nil, errors.New("Unexpected signing method in auth token")
		}
		return auth.ValidationKey, nil
	})

	if err != nil {
		log.Error().Err(err).Msg("Unable to parse claims")
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid || claims.UserID == "" {
		return nil, errors.New("invalid token: authentication failed")
	}
	return claims, nil
}
