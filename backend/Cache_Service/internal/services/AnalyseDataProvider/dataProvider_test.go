package analyseDataProvider

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"CacheService/internal/domain/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"log/slog"
)

// Mock provider for testing
type MockAnalysedDataProvider struct {
	mock.Mock
}

func (m *MockAnalysedDataProvider) SetValue(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	if len(args) > 0 {
		return args.Error(0)
	}
	return nil
}

func (m *MockAnalysedDataProvider) GetValue(ctx context.Context, key string) (interface{}, error) {
	args := m.Called(ctx, key)
	if len(args) > 0 {
		return args.Get(0), args.Error(1)
	}
	return nil, nil
}

func (m *MockAnalysedDataProvider) SetCard(ctx context.Context, value models.AnalysedData) error {
	args := m.Called(ctx, value)
	if len(args) > 0 {
		return args.Error(0)
	}
	return nil
}

// Test suite for RedisService
type RedisServiceTestSuite struct {
	suite.Suite
	service      *RedisService
	mockProvider *MockAnalysedDataProvider
	logger       *slog.Logger
	ctx          context.Context
}

func (s *RedisServiceTestSuite) SetupTest() {
	s.logger = slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	s.mockProvider = new(MockAnalysedDataProvider)
	s.service = NewRedisService(s.logger, s.mockProvider, 1*time.Hour)
	s.ctx = context.Background()
}

func TestRedisServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RedisServiceTestSuite))
}

func (s *RedisServiceTestSuite) TestSetAnalysedData_Error() {
	testData := models.AnalysedData{
		Stocks: []models.Stock{
			{StockName: "AAPL", Side: "BUY"},
		},
		PostText:        "Test post",
		PostURI:         "https://example.com",
		ChannelUsername: "testchannel",
		Date:            time.Now(),
		Reasoning:       "Test reasoning",
	}
	dataTitle := "test-data-title"

	expectedErr := fmt.Errorf("provider error")
	s.mockProvider.On("SetCard", s.ctx, testData).Return(expectedErr).Once()

	err := s.service.SetAnalysedData(s.ctx, dataTitle, testData)

	s.Error(err)
	s.ErrorContains(err, "provider error")
	s.mockProvider.AssertExpectations(s.T())
}

func (s *RedisServiceTestSuite) TestGetAnalysedData_Error() {
	dataTitle := "test-data-title"
	expectedErr := fmt.Errorf("provider error")

	s.mockProvider.On("GetValue", s.ctx, dataTitle).Return(nil, expectedErr).Once()
	data, err := s.service.GetAnalysedData(s.ctx, dataTitle)

	s.Error(err)
	s.Contains(err.Error(), "provider error")
	s.Nil(data)
	s.mockProvider.AssertExpectations(s.T())
}

func (s *RedisServiceTestSuite) TestGetAnalysedData_Success() {
	dataTitle := "test-data-title"
	expectedData := []byte("test-value")

	s.mockProvider.On("GetValue", s.ctx, dataTitle).Return(expectedData, nil).Once()

	data, err := s.service.GetAnalysedData(s.ctx, dataTitle)

	s.NoError(err)
	s.Equal(expectedData, data)
	s.mockProvider.AssertExpectations(s.T())
}


