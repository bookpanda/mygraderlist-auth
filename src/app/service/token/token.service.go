package token

import (
	"time"

	dto "github.com/bookpanda/mygraderlist-auth/src/app/dto/auth"
	model "github.com/bookpanda/mygraderlist-auth/src/app/model/auth"
	"github.com/bookpanda/mygraderlist-auth/src/config"
	role "github.com/bookpanda/mygraderlist-auth/src/constant/auth"
	auth_proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/auth"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Service struct {
	jwtService      IJwtService
	cacheRepository ICacheRepository
}

type IJwtService interface {
	SignAuth(*model.Auth) (string, error)
	VerifyAuth(string) (*jwt.Token, error)
	GetConfig() *config.Jwt
}

type ICacheRepository interface {
	SaveCache(string, interface{}, int) error
	GetCache(string, interface{}) error
}

func NewTokenService(jwtService IJwtService, cacheRepository ICacheRepository) *Service {
	return &Service{
		jwtService:      jwtService,
		cacheRepository: cacheRepository,
	}
}

func (s *Service) CreateCredentials(auth *model.Auth, secret string) (*auth_proto.Credential, error) {
	token, err := s.jwtService.SignAuth(auth)
	if err != nil {
		return nil, err
	}

	cache := dto.CacheAuth{
		Token: token,
		Role:  role.Role(auth.Role),
	}

	err = s.cacheRepository.SaveCache(auth.UserID, &cache, int(s.jwtService.GetConfig().ExpiresIn))
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "auth").
			Str("module", "validate").
			Msg("Cannot connect to cache server")
		return nil, errors.New("Internal service error")
	}

	credential := &auth_proto.Credential{
		AccessToken:  token,
		RefreshToken: s.CreateRefreshToken(),
		ExpiresIn:    s.jwtService.GetConfig().ExpiresIn,
	}

	return credential, nil
}

func (s *Service) Validate(token string) (*dto.UserCredential, error) {
	t, err := s.jwtService.VerifyAuth(token)
	if err != nil {
		return nil, err
	}

	payload := t.Claims.(jwt.MapClaims)

	if payload["iss"] != s.jwtService.GetConfig().Issuer {
		return nil, errors.New("Invalid token")
	}

	if time.Unix(int64(payload["exp"].(float64)), 0).Before(time.Now()) {
		return nil, errors.New("Token is expired")
	}

	cache := dto.CacheAuth{}
	err = s.cacheRepository.GetCache(payload["user_id"].(string), &cache)
	if err != nil {
		if err != redis.Nil {
			log.Error().
				Err(err).
				Str("service", "auth").
				Str("module", "validate").
				Msg("Cannot connect to cache server")
			return nil, errors.New("Internal service error")
		}

		return nil, errors.New("Invalid token")
	}

	if cache.Token != token {
		return nil, errors.New("Invalid token")
	}

	return &dto.UserCredential{
		UserId: payload["user_id"].(string),
		Role:   cache.Role,
	}, nil
}

func (s *Service) CreateRefreshToken() string {
	return uuid.New().String()
}
