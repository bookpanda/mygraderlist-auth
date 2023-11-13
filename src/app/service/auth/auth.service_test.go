package auth

import (
	"context"
	"testing"
	"time"

	"github.com/bookpanda/mygraderlist-auth/src/client"
	mock "github.com/bookpanda/mygraderlist-auth/src/mocks/auth"
	"golang.org/x/oauth2"

	dto "github.com/bookpanda/mygraderlist-auth/src/app/dto/auth"
	"github.com/bookpanda/mygraderlist-auth/src/app/model"
	"github.com/bookpanda/mygraderlist-auth/src/app/model/auth"
	"github.com/bookpanda/mygraderlist-auth/src/app/utils"
	"github.com/bookpanda/mygraderlist-auth/src/config"
	role "github.com/bookpanda/mygraderlist-auth/src/constant/auth"
	auth_proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/auth"
	user_proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/user"
	"github.com/bxcodec/faker/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AuthServiceTest struct {
	suite.Suite
	Auth              *auth.Auth
	UserDto           *user_proto.User
	Credential        *auth_proto.Credential
	Payload           *dto.TokenPayloadAuth
	UserCredential    *dto.UserCredential
	conf              config.App
	oauthConf         oauth2.Config
	googleOauthClient *client.GoogleOauthClient
	UnauthorizedErr   error
	NotFoundErr       error
	ServiceDownErr    error
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTest))
}

func (t *AuthServiceTest) SetupTest() {
	t.Auth = &auth.Auth{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		UserID:       faker.UUIDDigit(),
		Role:         role.USER,
		RefreshToken: faker.Word(),
	}

	t.UserDto = &user_proto.User{
		Id:       t.Auth.UserID,
		Username: faker.Username(),
		Email:    faker.Email(),
	}

	t.Credential = &auth_proto.Credential{
		AccessToken:  faker.Word(),
		RefreshToken: t.Auth.RefreshToken,
		ExpiresIn:    3600,
	}

	t.Payload = &dto.TokenPayloadAuth{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    faker.Word(),
			ExpiresAt: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserId: t.Auth.UserID,
	}

	t.UserCredential = &dto.UserCredential{
		UserId: t.Auth.UserID,
		Role:   role.Role(t.Auth.Role),
	}

	t.UnauthorizedErr = errors.New("unauthorized")
	t.NotFoundErr = errors.New("not found user")
	t.ServiceDownErr = errors.New("service is down")

	t.conf = config.App{
		Port:   3001,
		Debug:  false,
		Secret: "asuperstrong32bitpasswordgohere!",
	}

	t.oauthConf = oauth2.Config{}
}

func (t *AuthServiceTest) TestValidateSuccess() {
	want := &auth_proto.ValidateResponse{
		UserId: t.UserDto.Id,
		Role:   t.Auth.Role,
	}
	token := faker.Word()

	repo := &mock.RepositoryMock{}

	userService := &mock.UserServiceMock{}

	tokenService := &mock.TokenServiceMock{}
	tokenService.On("Validate", token).Return(t.UserCredential, nil)

	srv := NewService(repo, tokenService, userService, t.conf, &t.oauthConf, t.googleOauthClient)

	actual, err := srv.Validate(context.Background(), &auth_proto.ValidateRequest{Token: token})

	assert.Nilf(t.T(), err, "error: %v", err)
	assert.Equal(t.T(), want, actual)
}

func (t *AuthServiceTest) TestValidateInvalidToken() {
	token := faker.Word()

	repo := &mock.RepositoryMock{}

	userService := &mock.UserServiceMock{}

	tokenService := &mock.TokenServiceMock{}
	tokenService.On("Validate", token).Return(nil, errors.New("Invalid token"))

	srv := NewService(repo, tokenService, userService, t.conf, &t.oauthConf, t.googleOauthClient)

	actual, err := srv.Validate(context.Background(), &auth_proto.ValidateRequest{Token: token})

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Unauthenticated, st.Code())
}

func (t *AuthServiceTest) TestRedeemRefreshTokenSuccess() {
	token := faker.Word()
	t.Auth.RefreshToken = utils.Hash([]byte(t.Credential.RefreshToken))

	want := &auth_proto.RefreshTokenResponse{Credential: t.Credential}

	repo := &mock.RepositoryMock{}
	repo.On("FindByRefreshToken", utils.Hash([]byte(token)), &auth.Auth{}).Return(t.Auth, nil)
	repo.On("Update", t.Auth).Return(t.Auth, nil)

	userService := &mock.UserServiceMock{}

	tokenService := &mock.TokenServiceMock{}
	tokenService.On("CreateRefreshToken").Return(token)
	tokenService.On("CreateCredentials", t.Auth, t.conf.Secret).Return(t.Credential, nil)

	srv := NewService(repo, tokenService, userService, t.conf, &t.oauthConf, t.googleOauthClient)

	actual, err := srv.RefreshToken(context.Background(), &auth_proto.RefreshTokenRequest{RefreshToken: token})

	assert.Nilf(t.T(), err, "error: %v", err)
	assert.Equal(t.T(), want, actual)
}

func (t *AuthServiceTest) TestRedeemRefreshTokenInvalidToken() {
	token := faker.Word()
	t.Credential.RefreshToken = utils.Hash([]byte(token))

	repo := &mock.RepositoryMock{}
	repo.On("FindByRefreshToken", t.Credential.RefreshToken, &auth.Auth{}).Return(nil, errors.New("Not found token"))
	repo.On("Update", t.Auth).Return(t.Auth, nil)

	userService := &mock.UserServiceMock{}

	tokenService := &mock.TokenServiceMock{}
	tokenService.On("CreateRefreshToken").Return(token)
	tokenService.On("CreateCredentials", t.Auth, t.conf.Secret).Return(t.Credential, nil)

	srv := NewService(repo, tokenService, userService, t.conf, &t.oauthConf, t.googleOauthClient)

	actual, err := srv.RefreshToken(context.Background(), &auth_proto.RefreshTokenRequest{RefreshToken: token})

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Unauthenticated, st.Code())
}

func (t *AuthServiceTest) TestRedeemRefreshTokenInternalErr() {
	token := faker.Word()
	t.Credential.RefreshToken = utils.Hash([]byte(token))

	repo := &mock.RepositoryMock{}
	repo.On("FindByRefreshToken", t.Credential.RefreshToken, &auth.Auth{}).Return(t.Auth, nil)

	userService := &mock.UserServiceMock{}

	tokenService := &mock.TokenServiceMock{}
	tokenService.On("CreateRefreshToken").Return(token)
	tokenService.On("CreateCredentials", t.Auth, t.conf.Secret).Return(nil, errors.New("Invalid secret key"))

	srv := NewService(repo, tokenService, userService, t.conf, &t.oauthConf, t.googleOauthClient)

	actual, err := srv.RefreshToken(context.Background(), &auth_proto.RefreshTokenRequest{RefreshToken: token})

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Internal, st.Code())
}

func (t *AuthServiceTest) TestCreateCredentialsSuccess() {
	token := faker.Word()
	t.Credential.RefreshToken = utils.Hash([]byte(faker.Word()))

	want := t.Credential

	repo := &mock.RepositoryMock{}
	repo.On("Update", t.Auth).Return(t.Auth, nil)

	userService := &mock.UserServiceMock{}

	tokenService := &mock.TokenServiceMock{}
	tokenService.On("CreateRefreshToken").Return(token)
	tokenService.On("CreateCredentials", t.Auth, t.conf.Secret).Return(t.Credential, nil)

	srv := NewService(repo, tokenService, userService, t.conf, &t.oauthConf, t.googleOauthClient)

	credentials, err := srv.CreateNewCredential(t.Auth)

	assert.Nilf(t.T(), err, "error: %v", err)
	assert.Equal(t.T(), want, credentials)
}

func (t *AuthServiceTest) TestCreateCredentialsInternalErr() {
	want := errors.New("Invalid secret key")

	token := faker.Word()
	t.Credential.RefreshToken = utils.Hash([]byte(faker.Word()))

	repo := &mock.RepositoryMock{}
	repo.On("Update", t.Auth).Return(t.Auth, nil)

	userService := &mock.UserServiceMock{}

	tokenService := &mock.TokenServiceMock{}
	tokenService.On("CreateRefreshToken").Return(token)
	tokenService.On("CreateCredentials", t.Auth, t.conf.Secret).Return(nil, errors.New("Invalid secret key"))

	srv := NewService(repo, tokenService, userService, t.conf, &t.oauthConf, t.googleOauthClient)

	credentials, err := srv.CreateNewCredential(t.Auth)

	assert.Nil(t.T(), credentials)
	assert.Equal(t.T(), want.Error(), err.Error())
}
