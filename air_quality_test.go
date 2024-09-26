package omgo_test

import (
	"context"
	"testing"

	"github.com/jdotcurs/omgo"
	"github.com/stretchr/testify/require"
)

func TestGetAirQuality(t *testing.T) {
	c, err := omgo.NewClient()
	require.NoError(t, err)

	loc, err := omgo.NewLocation(52.3738, 4.8910) // Amsterdam
	require.NoError(t, err)

	opts := &omgo.Options{
		AirQualityMetrics: []string{"pm10", "pm2_5", "o3", "no2"},
	}

	aqData, err := c.GetAirQuality(context.Background(), loc, opts)
	require.NoError(t, err)

	require.IsType(t, omgo.AirQualityData{}, aqData)
	require.GreaterOrEqual(t, aqData.PM10, float64(0))
	require.GreaterOrEqual(t, aqData.PM2_5, float64(0))
	require.GreaterOrEqual(t, aqData.O3, float64(0))
	require.GreaterOrEqual(t, aqData.NO2, float64(0))
}

func TestGetAirQuality_InvalidLocation(t *testing.T) {
	c, err := omgo.NewClient()
	require.NoError(t, err)

	loc, err := omgo.NewLocation(1000, 1000) // Invalid location
	require.NoError(t, err)

	opts := &omgo.Options{
		AirQualityMetrics: []string{"pm10", "pm2_5", "o3", "no2"},
	}

	_, err = c.GetAirQuality(context.Background(), loc, opts)
	require.Error(t, err)
}

func TestGetAirQuality_NilOptions(t *testing.T) {
	c, err := omgo.NewClient()
	require.NoError(t, err)

	loc, err := omgo.NewLocation(52.3738, 4.8910) // Amsterdam
	require.NoError(t, err)

	aqData, err := c.GetAirQuality(context.Background(), loc, nil)
	require.NoError(t, err)
	require.IsType(t, omgo.AirQualityData{}, aqData)
}

func TestParseAirQualityData_InvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{"air_quality": {invalid}}`)
	_, err := omgo.ParseAirQualityData(invalidJSON)
	require.Error(t, err)
	require.IsType(t, omgo.ErrAPIResponse{}, err)
}
