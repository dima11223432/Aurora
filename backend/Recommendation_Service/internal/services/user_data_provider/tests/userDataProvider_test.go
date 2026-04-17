package auth_test

import (
	"context"
	"log/slog"
	auth "recommendationService/internal/services/user_data_provider"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type userDataProviderSuite struct {
	suite.Suite
	mock    *userDataProviderMock
	service *auth.UserDataProvider
}

func (u *userDataProviderSuite) SetupTest() {
	u.mock = new(userDataProviderMock)
	u.service = auth.New(slog.Default(), u.mock, u.mock, u.mock, 5*time.Minute)
}

func (u *userDataProviderSuite) TestGetUserPriorityChannels() {
	ctx := context.Background()
	u.mock.On("GetAllParsingChannels", ctx).Return([]string{"channel1", "channel2"}, nil)
	channels, err := u.service.GetAllParsingChannels(ctx)
	u.NoError(err)
	u.Equal([]string{"channel1", "channel2"}, channels)
}

func TestUserDataProviderTestSuite(t *testing.T) {
	suite.Run(t, new(userDataProviderSuite))
}
