package user_data_provider

import (
	"context"
	"errors"
	"fmt"

)
type PriorityChannel struct {
    ID       int64
    Name     string
    Priority int
}
type PriorityChannelsProvider interface {
    GetPriorityChannelsByUserID(ctx context.Context, userID int64) ([]models.PriorityChannel, error)
}

type UserDataProvider struct {
    priorityChannelsProvider PriorityChannelsProvider 
}

func NewUserDataProvider(provider PriorityChannelsProvider) *UserDataProvider {
    return &UserDataProvider{
        priorityChannelsProvider: provider,
    }
}

func (u *UserDataProvider) GetUserPriorityChannels(ctx context.Context, userID int64) ([]models.PriorityChannel, error) {
    const op = "internal.services.user_data_provider.userDataProvider.go.GetUserPriorityChannels"

    channels, err := u.priorityChannelsProvider.GetPriorityChannelsByUserID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }
    
    return channels, nil
}