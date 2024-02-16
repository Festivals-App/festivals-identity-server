package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	ADMIN   int = 42
	CREATOR int = 1
)

type User struct {
	ID           int       `json:"user_id" sql:"user_id"`
	Email        string    `json:"user_email" sql:"user_email"`
	PasswordHash string    `json:"user_password" sql:"user_password"`
	CreateDate   time.Time `json:"user_createdat" sql:"user_createdat"`
	UpdateDate   time.Time `json:"user_updatedat" sql:"user_updatedat"`
	Role         int       `json:"user_role" sql:"user_role"`
}

type UserClaims struct {
	UserID        string
	UserRole      int
	UserFestivals []int
	UserArtists   []int
	UserLocations []int
	jwt.RegisteredClaims
}
