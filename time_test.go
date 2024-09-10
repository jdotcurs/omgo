package omgo_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/hectormalot/omgo"
	"github.com/stretchr/testify/require"
)

func TestApiTime_MarshalJSON(t *testing.T) {
	// Test non-zero time
	nonZeroTime := omgo.ApiTime{Time: time.Date(2023, 5, 1, 12, 0, 0, 0, time.UTC)}
	data, err := json.Marshal(nonZeroTime)
	require.NoError(t, err)
	require.Equal(t, `"2023-05-01T12:00:00Z"`, string(data))

	// Test zero time
	zeroTime := &omgo.ApiTime{}
	data, err = json.Marshal(zeroTime)
	require.NoError(t, err)
	require.Equal(t, `null`, string(data))
}

func TestApiTime_IsSet(t *testing.T) {
	nonZeroTime := omgo.ApiTime{Time: time.Now()}
	require.True(t, nonZeroTime.IsSet())

	zeroTime := omgo.ApiTime{}
	require.False(t, zeroTime.IsSet())
}

func TestApiDate_MarshalJSON(t *testing.T) {
	// Test non-zero date
	nonZeroDate := omgo.ApiDate{Time: time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)}
	data, err := json.Marshal(nonZeroDate)
	require.NoError(t, err)
	require.Equal(t, `"2023-05-01T00:00:00Z"`, string(data))

	// Test zero date
	zeroDate := &omgo.ApiDate{}
	data, err = json.Marshal(zeroDate)
	require.NoError(t, err)
	require.Equal(t, `null`, string(data))
}

func TestApiDate_IsSet(t *testing.T) {
	nonZeroDate := omgo.ApiDate{Time: time.Now()}
	require.True(t, nonZeroDate.IsSet())

	zeroDate := omgo.ApiDate{}
	require.False(t, zeroDate.IsSet())
}
