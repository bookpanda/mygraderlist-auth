package auth

import (
	"github.com/bookpanda/mygraderlist-auth/src/constant/auth"
	"github.com/golang-jwt/jwt/v4"
)

type TokenPayloadAuth struct {
	jwt.RegisteredClaims
	UserId string `json:"user_id"`
}

type UserCredential struct {
	UserId string    `json:"user_id"`
	Role   auth.Role `json:"role"`
}

type CacheAuth struct {
	Token string    `json:"token"`
	Role  auth.Role `json:"role"`
}
