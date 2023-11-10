package auth

import (
	dto "github.com/bookpanda/mygraderlist-auth/src/app/dto/auth"
	model "github.com/bookpanda/mygraderlist-auth/src/app/model/auth"
	"github.com/bookpanda/mygraderlist-auth/src/config"
	"github.com/bookpanda/mygraderlist-auth/src/proto"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) FindByRefreshToken(id string, result *model.Auth) error {
	args := r.Called(id, result)

	if args.Get(0) != nil {
		*result = *args.Get(0).(*model.Auth)
	}

	return args.Error(1)
}

func (r *RepositoryMock) FindByUserID(id string, in *model.Auth) error {
	args := r.Called(id, in)

	if args.Get(0) != nil {
		*in = *args.Get(0).(*model.Auth)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Create(in *model.Auth) error {
	args := r.Called(in)

	if args.Get(0) != nil {
		*in = *args.Get(0).(*model.Auth)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Update(id string, in *model.Auth) error {
	args := r.Called(in)

	if args.Get(0) != nil {
		*in = *args.Get(0).(*model.Auth)
	}

	return args.Error(1)
}

type UserServiceMock struct {
	mock.Mock
}

func (c *UserServiceMock) FindByEmail(email string) (result *proto.User, err error) {
	args := c.Called(email)

	if args.Get(0) != nil {
		result = args.Get(0).(*proto.User)
	}

	return result, args.Error(1)
}

func (c *UserServiceMock) Create(in *proto.User) (result *proto.User, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		result = args.Get(0).(*proto.User)
	}

	return result, args.Error(1)
}

type JwtServiceMock struct {
	mock.Mock
}

func (s *JwtServiceMock) SignAuth(in *model.Auth) (token string, err error) {
	args := s.Called(in)

	return args.String(0), args.Error(1)
}

func (s *JwtServiceMock) VerifyAuth(token string) (decode *jwt.Token, err error) {
	args := s.Called(token)

	if args.Get(0) != nil {
		decode = args.Get(0).(*jwt.Token)
	}

	return decode, args.Error(1)
}

func (s *JwtServiceMock) GetConfig() *config.Jwt {
	args := s.Called()

	return args.Get(0).(*config.Jwt)
}

type TokenServiceMock struct {
	mock.Mock
}

func (s *TokenServiceMock) CreateCredentials(in *model.Auth, secret string) (credential *proto.Credential, err error) {
	args := s.Called(in, secret)

	if args.Get(0) != nil {
		credential = args.Get(0).(*proto.Credential)
	}

	return credential, args.Error(1)
}

func (s *TokenServiceMock) Validate(token string) (payload *dto.UserCredential, err error) {
	args := s.Called(token)

	if args.Get(0) != nil {
		payload = args.Get(0).(*dto.UserCredential)
	}

	return payload, args.Error(1)
}
