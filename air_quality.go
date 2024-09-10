package omgo

import (
	"context"
	"encoding/json"
	"fmt"
)

type AirQualityData struct {
	PM10  float64 `json:"pm10"`
	PM2_5 float64 `json:"pm2_5"`
	O3    float64 `json:"o3"`
	NO2   float64 `json:"no2"`
	// Add more fields as needed
}

func (c Client) GetAirQuality(ctx context.Context, loc Location, opts *Options) (AirQualityData, error) {
	if opts == nil {
		opts = &Options{}
	}
	if len(opts.AirQualityMetrics) == 0 {
		opts.AirQualityMetrics = []string{"pm10", "pm2_5", "o3", "no2"}
	}

	body, err := c.Get(ctx, loc, opts)
	if err != nil {
		return AirQualityData{}, fmt.Errorf("failed to get air quality data: %w", err)
	}

	aqData, err := ParseAirQualityData(body)
	if err != nil {
		return AirQualityData{}, fmt.Errorf("failed to parse air quality data: %w", err)
	}

	return aqData, nil
}

func ParseAirQualityData(body []byte) (AirQualityData, error) {
	var data struct {
		AirQuality AirQualityData `json:"air_quality"`
	}

	err := json.Unmarshal(body, &data)
	if err != nil {
		return AirQualityData{}, ErrAPIResponse{StatusCode: 0, Message: "Failed to parse JSON response"}
	}

	return data.AirQuality, nil
}
