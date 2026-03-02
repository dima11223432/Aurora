package suite

import (
	"context"
	"recomendationService/internal/domain/models"
	"recomendationService/tests/suite"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func Test_GetPriotiryChannelsByUserID(t *testing.T) {
	st := suite.NewSuite(t)

	var expectedChannels = []models.PriorityChannel{
		{
			Channel: "channel1",
		},
		{
			Channel: "channel2",
		},
	}

	channels, err := st.PostgresController.GetPriorityChannelsByUserID(context.Background(), int64(1))
	assert.NoError(t, err)
	assert.Equal(t, expectedChannels, channels)

}
