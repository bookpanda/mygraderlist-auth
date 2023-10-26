package jwt

import (
	"time"

	dto "github.com/bookpanda/mygraderlist-auth/src/app/dto/auth"
	model "github.com/bookpanda/mygraderlist-auth/src/app/model/auth"
	"github.com/bookpanda/mygraderlist-auth/src/config"
	_jwt "github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

type IStrategy interface {
	AuthDecode(*_jwt.Token) (interface{}, error)
}

type Service struct {
	conf     config.Jwt
	strategy IStrategy
}

func NewJwtService(conf config.Jwt, strategy IStrategy) *Service {
	return &Service{
		conf:     conf,
		strategy: strategy,
	}
}

func (s *Service) SignAuth(in *model.Auth) (string, error) {
	payloads := &dto.TokenPayloadAuth{
		RegisteredClaims: _jwt.RegisteredClaims{
			Issuer:    s.conf.Issuer,
			ExpiresAt: _jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(s.conf.ExpiresIn))),
			IssuedAt:  _jwt.NewNumericDate(time.Now()),
		},
		UserId: in.UserID,
	}
	token := _jwt.NewWithClaims(_jwt.SigningMethodHS256, payloads)

	tokenStr, err := token.SignedString([]byte(s.conf.Secret))
	if err != nil {
		return "", errors.New("Error while signing the token")
	}

	return tokenStr, nil
}

func (s *Service) VerifyAuth(token string) (*_jwt.Token, error) {
	return _jwt.Parse(token, s.strategy.AuthDecode)
}

func (s *Service) GetConfig() *config.Jwt {
	return &s.conf
}
