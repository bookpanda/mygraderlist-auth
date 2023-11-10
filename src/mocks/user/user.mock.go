package user

import (
	"context"

	user_proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/user"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) FindOne(_ context.Context, in *user_proto.FindOneUserRequest, _ ...grpc.CallOption) (res *user_proto.FindOneUserResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*user_proto.FindOneUserResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) FindByEmail(_ context.Context, in *user_proto.FindByEmailUserRequest, _ ...grpc.CallOption) (res *user_proto.FindByEmailUserResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*user_proto.FindByEmailUserResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) Create(_ context.Context, in *user_proto.CreateUserRequest, _ ...grpc.CallOption) (res *user_proto.CreateUserResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*user_proto.CreateUserResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) Update(_ context.Context, in *user_proto.UpdateUserRequest, _ ...grpc.CallOption) (res *user_proto.UpdateUserResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*user_proto.UpdateUserResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) Delete(_ context.Context, in *user_proto.DeleteUserRequest, _ ...grpc.CallOption) (res *user_proto.DeleteUserResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*user_proto.DeleteUserResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) Verify(_ context.Context, in *user_proto.VerifyUserRequest, _ ...grpc.CallOption) (res *user_proto.VerifyUserResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*user_proto.VerifyUserResponse)
	}

	return res, args.Error(1)
}
