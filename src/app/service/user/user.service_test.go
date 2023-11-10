package user

import (
	"testing"

	mock "github.com/bookpanda/mygraderlist-auth/src/mocks/user"
	user_proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/user"
	"github.com/bxcodec/faker/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceTest struct {
	suite.Suite
	UserDto         *user_proto.User
	UnauthorizedErr error
	NotFoundErr     error
	ServiceDownErr  error
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTest))
}

func (t *UserServiceTest) SetupTest() {
	t.UserDto = &user_proto.User{
		Id:       faker.UUIDDigit(),
		Username: faker.Username(),
		Email:    faker.Email(),
	}

	t.UnauthorizedErr = errors.New("Unauthorized")
	t.NotFoundErr = errors.New("Not found user")
	t.ServiceDownErr = errors.New("Service is down")
}

func (t *UserServiceTest) TestFindByEmailSuccess() {
	want := t.UserDto

	c := &mock.ClientMock{}
	c.On("FindByEmail", &user_proto.FindByEmailUserRequest{Email: t.UserDto.Email}).
		Return(&user_proto.FindByEmailUserResponse{User: t.UserDto}, nil)

	srv := NewUserService(c)

	actual, err := srv.FindByEmail(t.UserDto.Email)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *UserServiceTest) TestFindByEmailUnauthorized() {
	c := &mock.ClientMock{}
	c.On("FindByEmail", &user_proto.FindByEmailUserRequest{Email: t.UserDto.Email}).
		Return(nil, status.Error(codes.Unauthenticated, t.NotFoundErr.Error()))

	srv := NewUserService(c)

	actual, err := srv.FindByEmail(t.UserDto.Email)

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Unauthenticated, st.Code())
}

func (t *UserServiceTest) TestFindByEmailNotFound() {
	c := &mock.ClientMock{}
	c.On("FindByEmail", &user_proto.FindByEmailUserRequest{Email: t.UserDto.Email}).
		Return(nil, status.Error(codes.NotFound, t.NotFoundErr.Error()))

	srv := NewUserService(c)

	actual, err := srv.FindByEmail(t.UserDto.Email)

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, st.Code())
}

func (t *UserServiceTest) TestFindByEmailGrpcError() {
	c := &mock.ClientMock{}
	c.On("FindByEmail", &user_proto.FindByEmailUserRequest{Email: t.UserDto.Email}).
		Return(nil, status.Error(codes.Unavailable, t.ServiceDownErr.Error()))

	srv := NewUserService(c)

	actual, err := srv.FindByEmail(t.UserDto.Email)

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Unavailable, st.Code())
}

func (t *UserServiceTest) TestCreateSuccess() {
	want := t.UserDto

	c := &mock.ClientMock{}
	c.On("Create", &user_proto.CreateUserRequest{User: &user_proto.User{}}).
		Return(&user_proto.CreateUserResponse{User: t.UserDto}, nil)

	srv := NewUserService(c)

	actual, err := srv.Create(&user_proto.User{})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *UserServiceTest) TestCreateGrpcErr() {
	c := &mock.ClientMock{}
	c.On("Create", &user_proto.CreateUserRequest{User: &user_proto.User{}}).
		Return(nil, status.Error(codes.Unavailable, t.ServiceDownErr.Error()))

	srv := NewUserService(c)

	actual, err := srv.Create(&user_proto.User{})

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Unavailable, st.Code())
}
