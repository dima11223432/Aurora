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
	ssov1.UnimplementedRecommendationServiceServer
	userDataProvider UserDataProvider
}

const (
	emptyValue = 0
)

func Register(gRPC *grpc.Server, userDataProvider UserDataProvider) {
	ssov1.RegisterRecommendationServiceServer(gRPC, &serverAPI{
		userDataProvider: userDataProvider,
	})
}

func (s *serverAPI) GetUserPriorityChanneld(ctx context.Context,req *ssov1.GetUserPriotiryChannelsRequest) (*ssov1.GetUserPriotiryChannelsResponse, error) {

	const op = "internal.transport.grpc.serverAPI.GetUserPriorityChanneld"

	if req.GetUserId() == 0 {
		return nil, errors.New(op + ": user_id is empty")
	}

	channels, err := s.userDataProvider.GetUserPriorityChanneld(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	resp := &ssov1.GetUserPriotiryChannelsResponse{}

	for _, ch := range channels {
		resp.Channels = append(resp.Channels, &ssov1.PriorityChannel{
			ChannelId: ch.ChannelID,
			Priority:  int32(ch.Priority),
		})
	}

	return resp, nil
}
