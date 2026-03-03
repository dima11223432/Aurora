package grpcauth

import (
	"context"
	"errors"
	"recommendationService/internal/domain/models"

	ssov1 "recommendationService/api/gen/v1"

	"google.golang.org/grpc"
)

type UserDataProvider interface {
	GetUserPriorityChanneld(ctx context.Context, userID int64) ([]models.PriorityChannel, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServiceServer
	userDataProvider UserDataProvider
}

const (
	emptyValue = 0
)

func Register(gRPC *grpc.Server, userDataProvider UserDataProvider) {
	ssov1.RegisterAuthServiceServer(gRPC, &serverAPI{
		userDataProvider: userDataProvider,
	})
}

func GetUserPriorityChanneld(ctx context.Context, req *ssov1.GetUserPriotiryChannelsRequest) ([]models.PriorityChannel, error) {
	return nil, errors.New("Not implemented")
}
