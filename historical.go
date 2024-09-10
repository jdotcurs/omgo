package omgo

import (
	"context"
	"time"
)

type HistoricalData struct {
	StartDate time.Time
	EndDate   time.Time
	Forecast  Forecast // Reuse the existing Forecast struct
}

func (c Client) GetHistoricalData(ctx context.Context, loc Location, opts *Options) (HistoricalData, error) {
	if opts == nil {
		return HistoricalData{}, ErrInvalidInput{Param: "options", Value: nil}
	}

	if opts.StartDate == "" || opts.EndDate == "" {
		return HistoricalData{}, ErrInvalidInput{Param: "start_date or end_date", Value: "empty"}
	}

	body, err := c.Get(ctx, loc, opts)
	if err != nil {
		return HistoricalData{}, ErrInvalidInput{}
	}

	forecast, err := ParseBody(body)
	if err != nil {
		return HistoricalData{}, err
	}

	startDate, err := time.Parse("2006-01-02", opts.StartDate)
	if err != nil {
		return HistoricalData{}, ErrInvalidInput{Param: "start_date", Value: opts.StartDate}
	}

	endDate, err := time.Parse("2006-01-02", opts.EndDate)
	if err != nil {
		return HistoricalData{}, ErrInvalidInput{Param: "end_date", Value: opts.EndDate}
	}

	return HistoricalData{
		StartDate: startDate,
		EndDate:   endDate,
		Forecast:  *forecast,
	}, nil
}
