package grpcauth

import (
	"context"
	ssov1 "recommendationService/api/gen/v1"
	"recommendationService/internal/domain/models"

	"github.com/samber/lo"

	// grpcErr "recommendationService/internal/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserDataProvider interface {
	GetUserPriorityChannels(ctx context.Context, userID int64) ([]models.PriorityChannel, error)
}

type NewsDataProvider interface {
	GetRecommendatedPosts(ctx context.Context, userID int64, cursor *models.Cursor) ([]models.Post, *models.Cursor, error)
}

type serverAPI struct {
	ssov1.UnimplementedRecommendationServiceServer
	userDataProvider UserDataProvider
	newsDataProvider NewsDataProvider
}

const (
	emptyValue = 0
)

func Register(gRPC *grpc.Server, userDataProvider UserDataProvider, newsNewsDataProvider NewsDataProvider) {
	ssov1.RegisterRecommendationServiceServer(gRPC, &serverAPI{
		userDataProvider: userDataProvider,
		newsDataProvider: newsNewsDataProvider,
	})
}

func (s *serverAPI) GetUserPriorityChannels(ctx context.Context, req *ssov1.GetUserPriorityChannelsRequest) (*ssov1.GetUserPriorityChannelsResponse, error) {
	channels, err := s.userDataProvider.GetUserPriorityChannels(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.GetUserPriorityChannelsResponse{
		Channels: lo.Map(channels, func(item models.PriorityChannel, _ int) string {
			return item.Channel
		}),
	}, nil
}

func (s *serverAPI) GetRecommendatedPosts(ctx context.Context, req *ssov1.GetRecommendatedPostsRequest) (*ssov1.GetRecommendatedPostsResponse, error) {
	const op = "Recommendation_Service.internal.grpc.UserDataProvider.server.GetRecommendatedPosts"

	var cursor *models.Cursor
	if req.GetCursor() != nil && (req.GetCursor().GetScore() != 0 || req.GetCursor().GetId() != "") {
		cursor = &models.Cursor{
			Score: float64(req.GetCursor().GetScore()),
			ID:    req.GetCursor().GetId(),
		}
	}

	posts, nextCursor, err := s.newsDataProvider.GetRecommendatedPosts(ctx, req.GetUserId(), cursor)
	if nextCursor == nil {
		return nil, status.New(codes.OK, "no more posts").Err()
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get recommended posts")
	}
	protoPosts := make([]*ssov1.Post, 0)

	for _, post := range posts {
		stoks := post.Stocks
		protoStocks := make([]*ssov1.Stock, 0)
		for _, stock := range stoks {
			protoStocks = append(protoStocks, &ssov1.Stock{
				StockName: stock.StockName,
				Side:      stock.Side,
			})
		}
		protoPosts = append(protoPosts, &ssov1.Post{
			Stocks:          protoStocks,
			PostText:        post.PostText,
			PostUri:         post.PostURI,
			ChannelUsername: post.ChannelUsername,
			Reasoning:       post.Reasoning,
			Date:            timestamppb.New(post.Date),
		})
	}
	return &ssov1.GetRecommendatedPostsResponse{
		Posts: protoPosts,
		NextCursor: &ssov1.Cursor{
			Score: int64(nextCursor.Score),
			Id:    nextCursor.ID,
		},
	}, nil
}
