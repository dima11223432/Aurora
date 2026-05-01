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
	GetAllDefaultParsingChannels(ctx context.Context) ([]string, error)
	AddNewDefaultParsingChannel(ctx context.Context, channel string, category string) error
	AddNewUserCustomParsingChannel(ctx context.Context, userID int64, channel string) error
	DeleteDefaultParsingChannel(ctx context.Context, channel string) error
	DeleteUserCustomParsingChannel(ctx context.Context, userID int64, channel string) error
	GetDefaultParsingChannelsWithCategories(ctx context.Context) (map[string][]string, error)
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

func (s *serverAPI) GetAllDefaultParsingChannels(ctx context.Context, req *recv1.GetAllDefaultParsingChannelsRequest) (*recv1.GetAllDefaultParsingChannelsResponse, error) {
	channels, err := s.parsingChannelsProvider.GetAllDefaultParsingChannels(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &recv1.GetAllDefaultParsingChannelsResponse{
		Channels: channels,
	}, nil
}

func (s *serverAPI) AddNewUserCustomParsingChannel(
	ctx context.Context,
	req *recv1.AddNewUserCustomParsingChannelRequest,
) (*recv1.AddNewUserCustomParsingChannelResponse, error) {
	userID := req.GetUserId()
	channelUsername := req.GetChannelUsername()

	if err := s.parsingChannelsProvider.AddNewUserCustomParsingChannel(ctx, userID, channelUsername); err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &recv1.AddNewUserCustomParsingChannelResponse{}, nil
}

func (s *serverAPI) AddNewParsingChannel(ctx context.Context, req *recv1.AddNewDefaultParsingChannelRequest) (*recv1.AddNewDefaultParsingChannelResponse, error) {
	channelUsername := req.GetChannelUsername()
	channelCategory := req.GetCategory()
	if channelUsername == "" {
		return nil, status.Error(codes.InvalidArgument, "channel username is empty")
	}
	err := s.parsingChannelsProvider.AddNewDefaultParsingChannel(ctx, channelUsername, channelCategory)
	if err != nil {
		if errors.Is(err, storage.ErrChannelExists) {
			return nil, status.Error(codes.AlreadyExists, "channel already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &recv1.AddNewDefaultParsingChannelResponse{}, nil
}

func (s *serverAPI) DeleteUserCustomParsingChannel(
	ctx context.Context,
	req *recv1.DeleteUserCustomParsingChannelRequest,
) (*recv1.DeleteUserCustomParsingChannelResponse, error) {
	channelUsername := req.GetChannelUsername()
	userID := req.GetUserId()
	err := s.parsingChannelsProvider.DeleteUserCustomParsingChannel(ctx, userID, channelUsername)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &recv1.DeleteUserCustomParsingChannelResponse{}, nil
}

func (s *serverAPI) DeleteDefaultParsingChannel(ctx context.Context, req *recv1.DeleteDefaultParsingChannelRequest) (*recv1.DeleteDefaultParsingChannelResponse, error) {
	channelUsername := req.GetChannelUsername()
	if channelUsername == "" {
		return nil, status.Error(codes.InvalidArgument, "channel username is empty")
	}
	err := s.parsingChannelsProvider.DeleteDefaultParsingChannel(ctx, channelUsername)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &recv1.DeleteDefaultParsingChannelResponse{}, nil
}

func (s *serverAPI) GetAllDefaultParsingChannelsWithCategories(ctx context.Context, req *recv1.GetAllDefaultParsingChannelsWithCategoriesRequest) (*recv1.GetAllDefaultParsingChannelsWithCategoriesResponse, error) {
	channelsWithCategories, err := s.parsingChannelsProvider.GetDefaultParsingChannelsWithCategories(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	protoChannels := make(map[string]*recv1.ChannelList)

	for category, channels := range channelsWithCategories {
		protoChannels[category] = &recv1.ChannelList{
			Usernames: channels,
		}
	}

	return &recv1.GetAllDefaultParsingChannelsWithCategoriesResponse{
		Channels: protoChannels,
	}, nil
}
