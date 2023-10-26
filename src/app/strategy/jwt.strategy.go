package strategy

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

type JwtStrategy struct {
	secret string
}

func NewJwtStrategy(secret string) *JwtStrategy {
	return &JwtStrategy{secret: secret}
}

func (s *JwtStrategy) AuthDecode(token *jwt.Token) (interface{}, error) {
	if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
		return nil, errors.New(fmt.Sprintf("invalid token %v\n", token.Header["alg"]))
	}

	return []byte(s.secret), nil
}
