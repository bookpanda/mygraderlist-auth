package auth

import (
	"context"
	"net/url"
	"strings"

	dto "github.com/bookpanda/mygraderlist-auth/src/app/dto/auth"
	model "github.com/bookpanda/mygraderlist-auth/src/app/model/auth"
	"github.com/bookpanda/mygraderlist-auth/src/app/utils"
	"github.com/bookpanda/mygraderlist-auth/src/client"
	"github.com/bookpanda/mygraderlist-auth/src/config"
	role "github.com/bookpanda/mygraderlist-auth/src/constant/auth"
	auth_proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/auth"
	user_proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/user"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	repo              IRepository
	tokenService      ITokenService
	userService       IUserService
	conf              config.App
	oauthConfig       *oauth2.Config
	googleOauthClient *client.GoogleOauthClient
}

type IRepository interface {
	FindByRefreshToken(string, *model.Auth) error
	FindByUserID(string, *model.Auth) error
	Create(*model.Auth) error
	Update(string, *model.Auth) error
}

type IUserService interface {
	FindByEmail(string) (*user_proto.User, error)
	Create(*user_proto.User) (*user_proto.User, error)
}

type ITokenService interface {
	CreateCredentials(*model.Auth, string) (*auth_proto.Credential, error)
	Validate(string) (*dto.UserCredential, error)
}

func NewService(
	repo IRepository,
	tokenService ITokenService,
	userService IUserService,
	conf config.App,
	oauthConfig *oauth2.Config,
	googleOauthClient *client.GoogleOauthClient,
) *Service {
	return &Service{
		repo:              repo,
		tokenService:      tokenService,
		userService:       userService,
		conf:              conf,
		oauthConfig:       oauthConfig,
		googleOauthClient: googleOauthClient,
	}
}

func (s *Service) Validate(_ context.Context, req *auth_proto.ValidateRequest) (res *auth_proto.ValidateResponse, err error) {
	credential, err := s.tokenService.Validate(req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &auth_proto.ValidateResponse{
		UserId: credential.UserId,
		Role:   string(credential.Role),
	}, nil
}

func (s *Service) RefreshToken(_ context.Context, req *auth_proto.RefreshTokenRequest) (res *auth_proto.RefreshTokenResponse, err error) {
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

	return &auth_proto.RefreshTokenResponse{Credential: credentials}, nil
}

func (s *Service) CreateNewCredential(auth *model.Auth) (*auth_proto.Credential, error) {
	credentials, err := s.tokenService.CreateCredentials(auth, s.conf.Secret)
	log.Print("credentials FROM TOKENSERVICE", credentials)

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

func (s *Service) GetGoogleLoginUrl(context.Context, *auth_proto.GetGoogleLoginUrlRequest) (*auth_proto.GetGoogleLoginUrlResponse, error) {
	URL, err := url.Parse(s.oauthConfig.Endpoint.AuthURL)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse url")
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	parameters := url.Values{}
	parameters.Add("client_id", s.oauthConfig.ClientID)
	parameters.Add("scope", strings.Join(s.oauthConfig.Scopes, " "))
	parameters.Add("redirect_uri", s.oauthConfig.RedirectURL)
	parameters.Add("response_type", "code")
	URL.RawQuery = parameters.Encode()
	url := URL.String()

	return &auth_proto.GetGoogleLoginUrlResponse{
		Url: url,
	}, nil
}

func (s *Service) VerifyGoogleLogin(ctx context.Context, req *auth_proto.VerifyGoogleLoginRequest) (*auth_proto.VerifyGoogleLoginResponse, error) {
	code := req.GetCode()
	auth := model.Auth{}

	if code == "" {
		return nil, status.Error(codes.InvalidArgument, "No code is provided")
	}

	response, err := s.googleOauthClient.GetUserEmail(code)
	if err != nil {
		switch err.Error() {
		case "Invalid code":
			return nil, status.Error(codes.InvalidArgument, "Invalid code")
		default:
			log.Error().Err(err).Msg("Unable to get user info")
			return nil, status.Error(codes.Internal, "Internal server error")
		}
	}

	email := response.Email
	user, err := s.userService.FindByEmail(email)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				in := &user_proto.User{
					Email:    email,
					Username: response.Firstname,
				}

				user, err = s.userService.Create(in)
				if err != nil {
					return nil, status.Error(codes.InvalidArgument, st.Message())
				}

				auth = model.Auth{
					Role:   role.USER,
					UserID: user.Id,
				}

				err = s.repo.Create(&auth)
				if err != nil {
					log.Error().
						Err(err).
						Str("service", "auth").
						Str("module", "google").
						Msg("Error creating the auth data")
					return nil, status.Error(codes.Unavailable, st.Message())
				}

			default:
				log.Error().
					Err(err).
					Str("service", "auth").
					Str("module", "google").
					Msg("Service is down")
				return nil, status.Error(codes.Unavailable, st.Message())
			}
		} else {
			log.Error().
				Err(err).
				Str("service", "auth").
				Str("module", "google").
				Msg("Error connect to sso")
			return nil, status.Error(codes.Unavailable, "Service is down")
		}
	} else {
		err := s.repo.FindByUserID(user.Id, &auth)
		if err != nil {
			return nil, status.Error(codes.NotFound, "not found user")
		}
	}

	credentials, err := s.CreateNewCredential(&auth)
	log.Print("credentials", credentials)
	if err != nil {
		log.Error().Err(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Info().
		Str("service", "auth").
		Msg("User login to the service")

	return &auth_proto.VerifyGoogleLoginResponse{Credential: credentials}, err
}
