package services

import (
	"API_Service/internal/domains/models"
	"context"

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

func (r *RecommendationService) GetUserPriorityChannels(ctx context.Context) ([]models.Post, error) {
	return nil, nil
}
