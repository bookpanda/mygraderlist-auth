package user

import (
	"context"
	"time"

	"github.com/bookpanda/mygraderlist-auth/src/proto"
)

type Service struct {
	client proto.UserServiceClient
}

func NewUserService(client proto.UserServiceClient) *Service {
	return &Service{client: client}
}

func (s *Service) Create(user *proto.User) (*proto.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5000)
	defer cancel()

	res, err := s.client.Create(ctx, &proto.CreateUserRequest{User: user})
	if err != nil {
		return nil, err
	}

	return res.User, nil
}
