package user

import (
	"testing"

	mock "github.com/bookpanda/mygraderlist-auth/src/mocks/user"
	"github.com/bookpanda/mygraderlist-auth/src/proto"
	"github.com/bxcodec/faker/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceTest struct {
	suite.Suite
	UserDto         *proto.User
	UnauthorizedErr error
	NotFoundErr     error
	ServiceDownErr  error
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTest))
}

func (t *UserServiceTest) SetupTest() {
	t.UserDto = &proto.User{
		Id:       faker.UUIDDigit(),
		Username: faker.Username(),
		Email:    faker.Email(),
	}

	t.UnauthorizedErr = errors.New("Unauthorized")
	t.NotFoundErr = errors.New("Not found user")
	t.ServiceDownErr = errors.New("Service is down")
}

func (t *UserServiceTest) TestCreateSuccess() {
	want := t.UserDto

	c := &mock.ClientMock{}
	c.On("Create", &proto.CreateUserRequest{User: &proto.User{}}).
		Return(&proto.CreateUserResponse{User: t.UserDto}, nil)

	srv := NewUserService(c)

	actual, err := srv.Create(&proto.User{})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *UserServiceTest) TestCreateGrpcErr() {
	c := &mock.ClientMock{}
	c.On("Create", &proto.CreateUserRequest{User: &proto.User{}}).
		Return(nil, status.Error(codes.Unavailable, t.ServiceDownErr.Error()))

	srv := NewUserService(c)

	actual, err := srv.Create(&proto.User{})

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Unavailable, st.Code())
}