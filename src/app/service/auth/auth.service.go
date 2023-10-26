package auth

import (
	"context"

	dto "github.com/bookpanda/mygraderlist-auth/src/app/dto/auth"
	model "github.com/bookpanda/mygraderlist-auth/src/app/model/auth"
	"github.com/bookpanda/mygraderlist-auth/src/app/utils"
	"github.com/bookpanda/mygraderlist-auth/src/config"
	"github.com/bookpanda/mygraderlist-auth/src/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	repo         IRepository
	tokenService ITokenService
	userService  IUserService
	conf         config.App
}

type IRepository interface {
	FindByRefreshToken(string, *model.Auth) error
	FindByUserID(string, *model.Auth) error
	Create(*model.Auth) error
	Update(string, *model.Auth) error
}

type IUserService interface {
	FindByStudentID(string) (*proto.User, error)
	Create(*proto.User) (*proto.User, error)
}

type ITokenService interface {
	CreateCredentials(*model.Auth, string) (*proto.Credential, error)
	Validate(string) (*dto.UserCredential, error)
}

func NewService(
	repo IRepository,
	tokenService ITokenService,
	userService IUserService,
	conf config.App,
) *Service {
	return &Service{
		repo:         repo,
		tokenService: tokenService,
		userService:  userService,
		conf:         conf,
	}
}

func (s *Service) Validate(_ context.Context, req *proto.ValidateRequest) (res *proto.ValidateResponse, err error) {
	credential, err := s.tokenService.Validate(req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &proto.ValidateResponse{
		UserId: credential.UserId,
		Role:   string(credential.Role),
	}, nil
}

func (s *Service) RefreshToken(_ context.Context, req *proto.RefreshTokenRequest) (res *proto.RefreshTokenResponse, err error) {
	auth := model.Auth{}

	err = s.repo.FindByRefreshToken(utils.Hash([]byte(req.RefreshToken)), &auth)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid refresh token")
	}

	credentials, err := s.CreateNewCredential(&auth)
	if err != nil {
		log.Error().Err(err).
			Str("service", "auth").
			Str("module", "refresh token").
			Msg("Error while create new token")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.RefreshTokenResponse{Credential: credentials}, nil
}

func (s *Service) CreateNewCredential(auth *model.Auth) (*proto.Credential, error) {
	credentials, err := s.tokenService.CreateCredentials(auth, s.conf.Secret)
	if err != nil {
		return nil, err
	}

	auth.RefreshToken = utils.Hash([]byte(credentials.RefreshToken))

	err = s.repo.Update(auth.ID.String(), auth)
	if err != nil {
		return nil, err
	}

	return credentials, nil
}
