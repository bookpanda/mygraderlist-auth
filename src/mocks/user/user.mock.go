package user

import (
	"context"

	"github.com/bookpanda/mygraderlist-auth/src/proto"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) FindByStudentID(_ context.Context, in *proto.FindByStudentIDUserRequest, _ ...grpc.CallOption) (res *proto.FindByStudentIDUserResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.FindByStudentIDUserResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) Create(_ context.Context, in *proto.CreateUserRequest, _ ...grpc.CallOption) (res *proto.CreateUserResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.CreateUserResponse)
	}

	return res, args.Error(1)
}
