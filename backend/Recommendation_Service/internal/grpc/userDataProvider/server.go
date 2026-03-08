package grpcauth

import (
	"context"
	"github.com/samber/lo"
	ssov1 "recommendationService/api/gen/v1"
	"recommendationService/internal/domain/models"
	grpcErr "recommendationService/internal/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserDataProvider interface {
	GetUserPriorityChannels(ctx context.Context, userID int64) ([]models.PriorityChannel, error)
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

func (s *serverAPI) GetUserPriorityChanneld(ctx context.Context, req *ssov1.GetUserPriorityChannelsRequest) (*ssov1.GetUserPriorityChannelsResponse, error) {

	const op = "internal.transport.grpc.serverAPI.GetUserPriorityChanneld"

	if req.GetUserId() == emptyValue {
		return nil, grpcErr.ErrUserIDEmpty
	}

	channels, err := s.userDataProvider.GetUserPriorityChannels(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	//NOTE: формировать gRPC ответ лучше сразу в конце
	resp := &ssov1.GetUserPriorityChannelsResponse{}

	for _, ch := range channels {
		//FIX: тип данных в resp.Channels -> string, а не *ssov1.PriorityChannel
		resp.Channels = append(resp.Channels, &ssov1.PriorityChannel{
			ChannelId: ch.ChannelID,
			Priority:  int32(ch.Priority),
		})
	}

	return resp, nil
}
