package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	token "github.com/Festivals-App/festivals-identity-server/jwt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func GenerateAccessToken(user *token.User, db *sql.DB, auth *token.AuthService) (string, error) {

	userID := fmt.Sprint(user.ID)
	userRole := user.Role
	userFestivals, err := GetEntitiesForUser(Festival, db, userID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch festivals for user.")
		return "", errors.New("could not generate access token. please try again later")
	}
	userArtists, err := GetEntitiesForUser(Artist, db, userID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch artists for user.")
		return "", errors.New("could not generate access token. please try again later")
	}
	userLocations, err := GetEntitiesForUser(Location, db, userID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch locations for user.")
		return "", errors.New("could not generate access token. please try again later")
	}
	userEvents, err := GetEntitiesForUser(Event, db, userID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch events for user.")
		return "", errors.New("could not generate access token. please try again later")
	}
	userLinks, err := GetEntitiesForUser(Link, db, userID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch links for user.")
		return "", errors.New("could not generate access token. please try again later")
	}
	userImages, err := GetEntitiesForUser(Image, db, userID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch images for user.")
		return "", errors.New("could not generate access token. please try again later")
	}
	userPlaces, err := GetEntitiesForUser(Place, db, userID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch places for user.")
		return "", errors.New("could not generate access token. please try again later")
	}
	userTags, err := GetEntitiesForUser(Tag, db, userID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch tags for user.")
		return "", errors.New("could not generate access token. please try again later")
	}

	claims := token.UserClaims{
		UserID:        userID,
		UserRole:      userRole,
		UserFestivals: userFestivals,
		UserArtists:   userArtists,
		UserLocations: userLocations,
		UserEvents:    userEvents,
		UserLinks:     userLinks,
		UserImages:    userImages,
		UserPlaces:    userPlaces,
		UserTags:      userTags,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(auth.TokenLifetime)),
			Issuer:    auth.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(auth.SigningKey)
}

func RegenerateAccessToken(user *token.User, oldClaims *token.UserClaims, db *sql.DB, auth *token.AuthService) (string, error) {

	userID := fmt.Sprint(user.ID)
	userRole := user.Role
	userFestivals, err := GetEntitiesForUser(Festival, db, userID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch festivals for user.")
		return "", errors.New("could not generate access token. please try again later")
	}
	userArtists, err := GetEntitiesForUser(Artist, db, userID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch artists for user.")
		return "", errors.New("could not generate access token. please try again later")
	}
	userLocations, err := GetEntitiesForUser(Location, db, userID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch locations for user.")
		return "", errors.New("could not generate access token. please try again later")
	}
	userEvents, err := GetEntitiesForUser(Event, db, userID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch events for user.")
		return "", errors.New("could not generate access token. please try again later")
	}
	userLinks, err := GetEntitiesForUser(Link, db, userID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch links for user.")
		return "", errors.New("could not generate access token. please try again later")
	}
	userImages, err := GetEntitiesForUser(Image, db, userID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch images for user.")
		return "", errors.New("could not generate access token. please try again later")
	}
	userPlaces, err := GetEntitiesForUser(Place, db, userID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch places for user.")
		return "", errors.New("could not generate access token. please try again later")
	}
	userTags, err := GetEntitiesForUser(Tag, db, userID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch tags for user.")
		return "", errors.New("could not generate access token. please try again later")
	}

	claims := token.UserClaims{
		UserID:        userID,
		UserRole:      userRole,
		UserFestivals: userFestivals,
		UserArtists:   userArtists,
		UserLocations: userLocations,
		UserEvents:    userEvents,
		UserLinks:     userLinks,
		UserImages:    userImages,
		UserPlaces:    userPlaces,
		UserTags:      userTags,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: oldClaims.ExpiresAt,
			Issuer:    auth.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(auth.SigningKey)
}
