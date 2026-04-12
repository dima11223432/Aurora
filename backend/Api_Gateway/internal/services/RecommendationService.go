package services

import (
	"API_Service/internal/domains/models"
	"context"
	"fmt"

	// v1 "github.com/dima11223432/Aurora_SSO_Protos/api/gen/v1"
	rsv1 "github.com/dima11223432/recommendationService_protos/api/gen/v1"
)

type RecommendationService struct {
	RecommendationClient rsv1.RecommendationServiceClient
	authinterceptor      AuthInterceptor
}

func NewRecommendationService(recommendationClient rsv1.RecommendationServiceClient, authinterceptor AuthInterceptor) *RecommendationService {
	return &RecommendationService{
		RecommendationClient: recommendationClient,
		authinterceptor:      authinterceptor,
	}
}

func (r *RecommendationService) GetUserRecommendatedPosts(ctx context.Context, cursor *models.Cursor) ([]models.Post, *models.Cursor, error) {
	const op = "Api_Service.internal.services.RecommendationService.GetUserPriorityPosts"
	userID, err := r.authinterceptor.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("%s, %w", op, err)
	}
	var newCursor *rsv1.Cursor
	if cursor != nil {
		newCursor = &rsv1.Cursor{
			Score: int64(cursor.Score),
			Id:    cursor.ID,
		}
	}
	res, err := r.RecommendationClient.GetRecommendatedPosts(
		ctx,
		&rsv1.GetRecommendatedPostsRequest{
			UserId: userID,
			Cursor: newCursor,
		},
	)
	var nextCursor *models.Cursor
	if res.NextCursor != nil {
		nextCursor = &models.Cursor{
			Score: float64(res.NextCursor.Score),
			ID:    res.NextCursor.Id,
		}
	}

	if err != nil {
		return nil, nil, fmt.Errorf("%s, %w", op, err)
	}
	var posts []models.Post
	for _, post := range res.Posts {
		var stocks []models.Stock
		for _, stock := range post.Stocks {
			stocks = append(stocks, models.Stock{
				StockName: stock.StockName,
				Side:      stock.Side,
			})
		}
		posts = append(posts, models.Post{
			Stocks:          stocks,
			PostText:        post.PostText,
			PostURI:         post.PostUri,
			ChannelUsername: post.ChannelUsername,
			Date:            post.Date.AsTime(),
		})
	}
	return posts, nextCursor, nil

}
