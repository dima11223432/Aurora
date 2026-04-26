package grpcauth

import (
	"context"
	"errors"
	recv1 "recommendationService/api/gen/v1"
	"recommendationService/internal/domain/models"
	"recommendationService/internal/storage"

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

type ParsingChannelsProvider interface {
	GetAllParsingChannels(ctx context.Context) ([]string, error)
	AddNewParsingChannel(ctx context.Context, channel string, category string) error
	DeleteParsingChannel(ctx context.Context, channel string) error
	GetParsingChannelsWithCategories(ctx context.Context) (map[string][]string, error)
	GetAllCategories(ctx context.Context) ([]string, error)
}

type NewsDataProvider interface {
	GetRecommendatedPosts(ctx context.Context, userID int64, cursor *models.Cursor) ([]models.Post, *models.Cursor, error)
}

type serverAPI struct {
	recv1.UnimplementedRecommendationServiceServer
	userDataProvider        UserDataProvider
	parsingChannelsProvider ParsingChannelsProvider
	newsDataProvider        NewsDataProvider
}

const (
	emptyValue = 0
)

func Register(gRPC *grpc.Server, userDataProvider UserDataProvider, parsingChannelsProvider ParsingChannelsProvider, newsNewsDataProvider NewsDataProvider) {
	recv1.RegisterRecommendationServiceServer(gRPC, &serverAPI{
		userDataProvider:        userDataProvider,
		parsingChannelsProvider: parsingChannelsProvider,
		newsDataProvider:        newsNewsDataProvider,
	})
}

func NewServerAPI(userDataProvider UserDataProvider, parsingChannelsProvider ParsingChannelsProvider, newsDataProvider NewsDataProvider) *serverAPI {
	return &serverAPI{
		userDataProvider:        userDataProvider,
		parsingChannelsProvider: parsingChannelsProvider,
		newsDataProvider:        newsDataProvider,
	}
}

func (s *serverAPI) GetUserPriorityChannels(ctx context.Context, req *recv1.GetUserPriorityChannelsRequest) (*recv1.GetUserPriorityChannelsResponse, error) {
	channels, err := s.userDataProvider.GetUserPriorityChannels(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &recv1.GetUserPriorityChannelsResponse{
		Channels: lo.Map(channels, func(item models.PriorityChannel, _ int) string {
			return item.Channel
		}),
	}, nil
}

func (s *serverAPI) GetRecommendatedPosts(ctx context.Context, req *recv1.GetRecommendatedPostsRequest) (*recv1.GetRecommendatedPostsResponse, error) {
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
	protoPosts := make([]*recv1.Post, 0)

	for _, post := range posts {
		stoks := post.Stocks
		protoStocks := make([]*recv1.Stock, 0)
		for _, stock := range stoks {
			protoStocks = append(protoStocks, &recv1.Stock{
				StockName: stock.StockName,
				Side:      stock.Side,
			})
		}
		protoPosts = append(protoPosts, &recv1.Post{
			Stocks:          protoStocks,
			PostText:        post.PostText,
			PostUri:         post.PostURI,
			ChannelUsername: post.ChannelUsername,
			Reasoning:       post.Reasoning,
			Date:            timestamppb.New(post.Date),
		})
	}
	return &recv1.GetRecommendatedPostsResponse{
		Posts: protoPosts,
		NextCursor: &recv1.Cursor{
			Score: int64(nextCursor.Score),
			Id:    nextCursor.ID,
		},
	}, nil
}

func (s *serverAPI) GetAllParsingChannels(ctx context.Context, req *recv1.GetAllParsingChannelsRequest) (*recv1.GetAllParsingChannelsResponse, error) {
	channels, err := s.parsingChannelsProvider.GetAllParsingChannels(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &recv1.GetAllParsingChannelsResponse{
		Channels: channels,
	}, nil
}

func (s *serverAPI) AddNewParsingChannel(ctx context.Context, req *recv1.AddNewParsingChannelRequest) (*recv1.AddNewParsingChannelResponse, error) {
	channelUsername := req.GetChannelUsername()
	channelCategory := req.GetCategory()
	if channelUsername == "" {
		return nil, status.Error(codes.InvalidArgument, "channel username is empty")
	}
	err := s.parsingChannelsProvider.AddNewParsingChannel(ctx, channelUsername, channelCategory)
	if err != nil {
		if errors.Is(err, storage.ErrChannelExists) {
			return nil, status.Error(codes.AlreadyExists, "channel already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &recv1.AddNewParsingChannelResponse{}, nil
}

func (s *serverAPI) DeleteParsingChannel(ctx context.Context, req *recv1.DeleteParsingChannelRequest) (*recv1.DeleteParsingChannelResponse, error) {
	channelUsername := req.GetChannelUsername()
	if channelUsername == "" {
		return nil, status.Error(codes.InvalidArgument, "channel username is empty")
	}
	err := s.parsingChannelsProvider.DeleteParsingChannel(ctx, channelUsername)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &recv1.DeleteParsingChannelResponse{}, nil
}

func (s *serverAPI) GetAllParsingChannelsWithCategories(ctx context.Context, req *recv1.GetAllParsingChannelsWithCategoriesRequest) (*recv1.GetAllParsingChannelsWithCategoriesResponse, error) {
	channelsWithCategories, err := s.parsingChannelsProvider.GetParsingChannelsWithCategories(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	protoChannels := make(map[string]*recv1.ChannelList)

	for category, channels := range channelsWithCategories {
		protoChannels[category] = &recv1.ChannelList{
			Usernames: channels,
		}
	}

	return &recv1.GetAllParsingChannelsWithCategoriesResponse{
		Channels: protoChannels,
	}, nil
}
