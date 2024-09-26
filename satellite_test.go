package omgo_test

import (
	"context"
	"testing"

	"github.com/jdotcurs/omgo"
	"github.com/stretchr/testify/require"
)

func TestGetSatelliteData(t *testing.T) {
	c, err := omgo.NewClient()
	require.NoError(t, err)

	loc, err := omgo.NewLocation(52.3738, 4.8910) // Amsterdam
	require.NoError(t, err)

	opts := &omgo.Options{
		SatelliteMetrics: []string{"cloud_cover", "infrared", "visible_light", "water_vapor"},
	}

	satData, err := c.GetSatelliteData(context.Background(), loc, opts)
	require.NoError(t, err)

	require.IsType(t, omgo.SatelliteData{}, satData)
	require.GreaterOrEqual(t, satData.CloudCover, float64(0))
	require.GreaterOrEqual(t, satData.Infrared, float64(0))
	require.GreaterOrEqual(t, satData.VisibleLight, float64(0))
	require.GreaterOrEqual(t, satData.WaterVapor, float64(0))
}

func TestGetSatelliteData_InvalidLocation(t *testing.T) {
	c, err := omgo.NewClient()
	require.NoError(t, err)

	loc, err := omgo.NewLocation(1000, 1000) // Invalid location
	require.NoError(t, err)

	opts := &omgo.Options{
		SatelliteMetrics: []string{"cloud_cover", "infrared", "visible_light", "water_vapor"},
	}

	_, err = c.GetSatelliteData(context.Background(), loc, opts)
	require.Error(t, err)
}

func TestGetSatelliteData_NilOptions(t *testing.T) {
	c, err := omgo.NewClient()
	require.NoError(t, err)

	loc, err := omgo.NewLocation(52.3738, 4.8910) // Amsterdam
	require.NoError(t, err)

	satData, err := c.GetSatelliteData(context.Background(), loc, nil)
	require.NoError(t, err)
	require.IsType(t, omgo.SatelliteData{}, satData)
}
