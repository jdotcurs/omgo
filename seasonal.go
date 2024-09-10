package omgo

import (
	"context"
	"time"
)

type SeasonalForecast struct {
	StartDate time.Time
	EndDate   time.Time
	Forecast  Forecast
}

func (c Client) GetSeasonalForecast(ctx context.Context, loc Location, opts *Options) (SeasonalForecast, error) {
	if opts == nil {
		return SeasonalForecast{}, ErrInvalidInput{Param: "options", Value: nil}
	}

	if !opts.SeasonalForecast {
		return SeasonalForecast{}, ErrInvalidInput{Param: "SeasonalForecast", Value: false}
	}

	if opts.ForecastMonths < 1 || opts.ForecastMonths > 6 {
		return SeasonalForecast{}, ErrInvalidInput{Param: "ForecastMonths", Value: opts.ForecastMonths}
	}

	body, err := c.Get(ctx, loc, opts)
	if err != nil {
		return SeasonalForecast{}, err
	}

	forecast, err := ParseBody(body)
	if err != nil {
		return SeasonalForecast{}, err
	}

	startDate := time.Now()
	endDate := startDate.AddDate(0, opts.ForecastMonths, 0)

	return SeasonalForecast{
		StartDate: startDate,
		EndDate:   endDate,
		Forecast:  *forecast,
	}, nil
}
