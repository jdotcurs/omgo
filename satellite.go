package omgo

import (
	"context"
	"encoding/json"
	"fmt"
)

type SatelliteData struct {
	CloudCover   float64 `json:"cloud_cover"`
	Infrared     float64 `json:"infrared"`
	VisibleLight float64 `json:"visible_light"`
	WaterVapor   float64 `json:"water_vapor"`
}

func (c Client) GetSatelliteData(ctx context.Context, loc Location, opts *Options) (SatelliteData, error) {
	if opts == nil {
		opts = &Options{}
	}
	if len(opts.SatelliteMetrics) == 0 {
		opts.SatelliteMetrics = []string{"cloud_cover", "infrared", "visible_light", "water_vapor"}
	}

	body, err := c.Get(ctx, loc, opts)
	if err != nil {
		return SatelliteData{}, fmt.Errorf("failed to get satellite data: %w", err)
	}

	satData, err := ParseSatelliteData(body)
	if err != nil {
		return SatelliteData{}, fmt.Errorf("failed to parse satellite data: %w", err)
	}

	return satData, nil
}

func ParseSatelliteData(body []byte) (SatelliteData, error) {
	var data struct {
		Satellite SatelliteData `json:"satellite"`
	}

	err := json.Unmarshal(body, &data)
	if err != nil {
		return SatelliteData{}, ErrAPIResponse{StatusCode: 0, Message: "Failed to parse JSON response"}
	}

	return data.Satellite, nil
}
