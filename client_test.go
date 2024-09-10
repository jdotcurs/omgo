package omgo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetAPIKey(t *testing.T) {
	c, err := NewClient()
	require.NoError(t, err)

	apiKey := "test_api_key"
	c.SetAPIKey(apiKey)
	require.Equal(t, apiKey, c.APIKey)
}
