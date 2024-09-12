package omgo_test

import (
	"context"
	"testing"
	"time"

	"github.com/hectormalot/omgo"
	"github.com/stretchr/testify/require"
)

func TestSetAPIKey(t *testing.T) {
	c, err := omgo.NewClient()
	require.NoError(t, err)

	apiKey := "test_api_key"
	c.SetAPIKey(apiKey)
	require.Equal(t, apiKey, c.APIKey)
}

func TestClientCaching(t *testing.T) {
	client, err := omgo.NewClient()
	require.NoError(t, err)

	loc, err := omgo.NewLocation(52.3738, 4.8910) // Amsterdam
	require.NoError(t, err)

	opts := &omgo.Options{
		DailyMetrics: []string{"temperature_2m_max"},
	}

	// First request
	_, err = client.Get(context.Background(), loc, opts)
	require.NoError(t, err)

	// Second request (should be cached)
	start := time.Now()
	_, err = client.Get(context.Background(), loc, opts)
	require.NoError(t, err)
	duration := time.Since(start)

	// The second request should be significantly faster due to caching
	require.Less(t, duration, 10*time.Millisecond)
}
