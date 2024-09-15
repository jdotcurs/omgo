package omgo

import (
	"context"
	"fmt"
	"time"
)

type HistoricalData struct {
	StartDate  time.Time
	EndDate    time.Time
	Forecast   Forecast
	HourlyData HourlyData
	DailyData  DailyData
}

type HourlyData struct {
	Time                   []time.Time
	Temperature2m          []float64
	RelativeHumidity2m     []float64
	DewPoint2m             []float64
	ApparentTemperature    []float64
	Precipitation          []float64
	Rain                   []float64
	Snowfall               []float64
	WindSpeed10m           []float64
	WindDirection10m       []float64
	WindGusts10m           []float64
	ShortwaveRadiation     []float64
	DirectNormalIrradiance []float64
	DiffuseRadiation       []float64
	CloudCover             []float64
	Visibility             []float64
	WeatherCode            []int
}

type DailyData struct {
	Time                     []time.Time
	WeatherCode              []int
	Temperature2mMax         []float64
	Temperature2mMin         []float64
	ApparentTemperatureMax   []float64
	ApparentTemperatureMin   []float64
	Sunrise                  []time.Time
	Sunset                   []time.Time
	PrecipitationSum         []float64
	RainSum                  []float64
	SnowfallSum              []float64
	PrecipitationHours       []float64
	WindSpeed10mMax          []float64
	WindGusts10mMax          []float64
	WindDirection10mDominant []float64
	ShortwaveRadiationSum    []float64
}

func (c Client) GetHistoricalData(ctx context.Context, loc Location, opts *Options) (HistoricalData, error) {
	if opts == nil {
		return HistoricalData{}, ErrInvalidInput{Param: "options", Value: nil}
	}

	if opts.StartDate == "" || opts.EndDate == "" {
		return HistoricalData{}, ErrInvalidInput{Param: "start_date or end_date", Value: "empty"}
	}

	startDate, err := time.Parse("2006-01-02", opts.StartDate)
	if err != nil {
		return HistoricalData{}, ErrInvalidInput{Param: "start_date", Value: opts.StartDate}
	}

	endDate, err := time.Parse("2006-01-02", opts.EndDate)
	if err != nil {
		return HistoricalData{}, ErrInvalidInput{Param: "end_date", Value: opts.EndDate}
	}

	body, err := c.Get(ctx, loc, opts)
	if err != nil {
		return HistoricalData{}, fmt.Errorf("failed to get data: %w", err)
	}

	forecast, err := ParseBody(body)
	if err != nil {
		return HistoricalData{}, fmt.Errorf("failed to parse body: %w", err)
	}

	historicalData, err := ParseHistoricalBody(body)
	if err != nil {
		return HistoricalData{}, fmt.Errorf("failed to parse historical body: %w", err)
	}

	return HistoricalData{
		StartDate:  startDate,
		EndDate:    endDate,
		Forecast:   *forecast,
		HourlyData: historicalData.HourlyData,
		DailyData:  historicalData.DailyData,
	}, nil
}
