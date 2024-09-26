package omgo_test

import (
	"context"
	"testing"

	"github.com/jdotcurs/omgo"
	"github.com/stretchr/testify/require"
)

func TestGetSeasonalForecast(t *testing.T) {
	c, err := omgo.NewClient()
	require.NoError(t, err)

	loc, err := omgo.NewLocation(52.3738, 4.8910) // Amsterdam
	require.NoError(t, err)

	opts := &omgo.Options{
		SeasonalForecast: true,
		ForecastMonths:   3,
		DailyMetrics:     []string{"temperature_2m_max", "temperature_2m_min"},
	}

	seasonalForecast, err := c.GetSeasonalForecast(context.Background(), loc, opts)
	require.NoError(t, err)

	require.Equal(t, 3, int(seasonalForecast.EndDate.Sub(seasonalForecast.StartDate).Hours()/24/30))
	require.NotEmpty(t, seasonalForecast.Forecast.DailyMetrics["temperature_2m_max"])
	require.NotEmpty(t, seasonalForecast.Forecast.DailyMetrics["temperature_2m_min"])
}

func TestGetSeasonalForecast_InvalidOptions(t *testing.T) {
	c, err := omgo.NewClient()
	require.NoError(t, err)

	loc, err := omgo.NewLocation(52.3738, 4.8910) // Amsterdam
	require.NoError(t, err)

	opts := &omgo.Options{
		SeasonalForecast: true,
		ForecastMonths:   7, // Invalid: should be 1-6
	}

	_, err = c.GetSeasonalForecast(context.Background(), loc, opts)
	require.Error(t, err)
	require.IsType(t, omgo.ErrInvalidInput{}, err)

	opts.ForecastMonths = 0 // Invalid: should be 1-6
	_, err = c.GetSeasonalForecast(context.Background(), loc, opts)
	require.Error(t, err)
	require.IsType(t, omgo.ErrInvalidInput{}, err)

	opts.SeasonalForecast = false // Invalid: should be true
	_, err = c.GetSeasonalForecast(context.Background(), loc, opts)
	require.Error(t, err)
	require.IsType(t, omgo.ErrInvalidInput{}, err)
}
