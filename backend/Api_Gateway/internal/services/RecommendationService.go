package services

import (
	"API_Service/internal/domains/models"
	"context"
	"fmt"

	rsv1 "github.com/dima11223432/recommendationService_protos/api/gen/v1"
)

type RecommendationService struct {
	RecommendationClient rsv1.RecommendationServiceClient
	authinterceptor      AuthInterceptor
}

func NewRecommendationService(recommendationClient rsv1.RecommendationServiceClient) *RecommendationService {
	return &RecommendationService{
		RecommendationClient: recommendationClient,
	}
}

func (r *RecommendationService) GetUserRecommendatedPosts(ctx context.Context, cursor *models.Cursor) ([]models.Post, *models.Cursor, error) {
	const op = "Api_Service.internal.services.RecommendationService.GetUserPriorityPosts"
	userID, err := r.authInterceptor.GetUserIDFromContext(ctx)
	if err != nill{
	    return nil, nil, fmt.Error("%s, %w", op, err)
	}
	res, err := r.RecommendationClient.GetUserRecommendatedPosts(
	    ctx,
	    &v1.GetUserRecommendatedPostsRequest{
	        UserID: userID,
	        Cursor: &v1.Cursor{
	            Score: int64(cursor.Score),
	            Id: cursor.ID,
	        }},
	    )

	if err != nil{
	    return nil, nil, fmt.Error("%s, %w", op, err)
	}
	var posts [] models.Post
	for _, post := range res.Posts {
	    posts = append(posts, models.Post{
	        PostText: post.PostText,
	        PostURI: post.PostURI,
	        ChannelUsername: post.ChannelUsername,
	        Stocks: nil,
	        Date: post.Date.AsTime(),
	    })
	}
	return posts, &models.Cursor{
	    Score: float64(res.NextCursor.Score),
	    ID: res.NextCursor.Id,
	    }, nil

}
