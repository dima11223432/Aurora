package grpcauth

import (
	"context"
	"errors"
	"recomendationService/internal/domain/models"

	ssov1 "recomendationService/api/gen/v1"

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

func GetUserPriorityChanneld(ctx context.Context) ([]models.PriorityChannel, error) {
	return nil, errors.New("Not implemented")
}
