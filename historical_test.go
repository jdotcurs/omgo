package omgo_test

import (
	"context"
	"testing"
	"time"

	"github.com/hectormalot/omgo"
	"github.com/stretchr/testify/require"
)

func TestGetHistoricalData(t *testing.T) {
	c, err := omgo.NewClient()
	require.NoError(t, err)

	loc, err := omgo.NewLocation(52.3738, 4.8910) // Amsterdam
	require.NoError(t, err)

	endDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

	opts := &omgo.Options{
		StartDate:     startDate,
		EndDate:       endDate,
		HourlyMetrics: []string{"temperature_2m", "precipitation", "wind_speed_10m"},
		DailyMetrics:  []string{"temperature_2m_max", "temperature_2m_min", "precipitation_sum"},
	}

	historicalData, err := c.GetHistoricalData(context.Background(), loc, opts)
	require.NoError(t, err)

	// Test Forecast data
	require.NotEmpty(t, historicalData.Forecast.HourlyTimes)
	require.NotEmpty(t, historicalData.Forecast.HourlyMetrics["temperature_2m"])
	require.NotEmpty(t, historicalData.Forecast.DailyTimes)
	require.NotEmpty(t, historicalData.Forecast.DailyMetrics["temperature_2m_max"])

	// Test HourlyData
	require.Equal(t, len(historicalData.HourlyData.Time), len(historicalData.HourlyData.Temperature2m))
	require.Equal(t, len(historicalData.HourlyData.Time), len(historicalData.HourlyData.Precipitation))
	require.Equal(t, len(historicalData.HourlyData.Time), len(historicalData.HourlyData.WindSpeed10m))

	// Test DailyData
	require.Equal(t, len(historicalData.DailyData.Time), len(historicalData.DailyData.Temperature2mMax))
	require.Equal(t, len(historicalData.DailyData.Time), len(historicalData.DailyData.Temperature2mMin))
	require.Equal(t, len(historicalData.DailyData.Time), len(historicalData.DailyData.PrecipitationSum))

	// Test start and end dates
	require.Equal(t, startDate, historicalData.StartDate.Format("2006-01-02"))
	require.Equal(t, endDate, historicalData.EndDate.Format("2006-01-02"))
}

func TestGetHistoricalData_InvalidDates(t *testing.T) {
	c, err := omgo.NewClient()
	require.NoError(t, err)

	loc, err := omgo.NewLocation(52.3738, 4.8910) // Amsterdam
	require.NoError(t, err)

	opts := &omgo.Options{
		StartDate: "invalid",
		EndDate:   "2023-04-30",
	}

	_, err = c.GetHistoricalData(context.Background(), loc, opts)
	require.Error(t, err)
	require.IsType(t, omgo.ErrInvalidInput{}, err)

	opts.StartDate = "2023-04-01"
	opts.EndDate = "invalid"

	_, err = c.GetHistoricalData(context.Background(), loc, opts)
	require.Error(t, err)
	require.IsType(t, omgo.ErrInvalidInput{}, err)
}

func TestGetHistoricalData_NilOptions(t *testing.T) {
	c, err := omgo.NewClient()
	require.NoError(t, err)

	loc, err := omgo.NewLocation(52.3738, 4.8910) // Amsterdam
	require.NoError(t, err)

	_, err = c.GetHistoricalData(context.Background(), loc, nil)
	require.Error(t, err)
	require.IsType(t, omgo.ErrInvalidInput{}, err)
}
