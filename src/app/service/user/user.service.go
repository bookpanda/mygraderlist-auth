package user

import (
	"context"
	"time"

	user_proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/user"
)

type Service struct {
	client user_proto.UserServiceClient
}

func NewUserService(client user_proto.UserServiceClient) *Service {
	return &Service{client: client}
}

func (s *Service) FindByEmail(email string) (*user_proto.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5000)
	defer cancel()

	res, err := s.client.FindByEmail(ctx, &user_proto.FindByEmailUserRequest{Email: email})
	if err != nil {
		return nil, err
	}

	return res.User, nil
}

func (s *Service) Create(user *user_proto.User) (*user_proto.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5000)
	defer cancel()

	res, err := s.client.Create(ctx, &user_proto.CreateUserRequest{User: user})
	if err != nil {
		return nil, err
	}

	return res.User, nil
}
