package omgo

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

type Client struct {
	URL         string
	UserAgent   string
	Client      *http.Client
	APIKey      string
	LastRequest time.Time
	RateLimiter *rate.Limiter
}

const DefaultUserAgent = "Open-Meteo_Go_Client"
const MinRequestInterval = time.Second / 10 // 10 requests per second

func NewClient() (Client, error) {
	return Client{
		URL:         "https://api.open-meteo.com/v1/forecast",
		UserAgent:   DefaultUserAgent,
		Client:      http.DefaultClient,
		RateLimiter: rate.NewLimiter(rate.Every(time.Second/10), 1), // 10 requests per second
	}, nil
}

// Add a method to set the API key
func (c *Client) SetAPIKey(key string) {
	c.APIKey = key
}

type Location struct {
	lat, lon float64
}

func NewLocation(lat, lon float64) (Location, error) {
	return Location{lat: lat, lon: lon}, nil
}

type Options struct {
	TemperatureUnit   string   // Default "celsius"
	WindspeedUnit     string   // Default "kmh",
	PrecipitationUnit string   // Default "mm"
	Timezone          string   // Default "UTC"
	PastDays          int      // Default 0
	HourlyMetrics     []string // Lists required hourly metrics, see https://open-meteo.com/en/docs for valid metrics
	DailyMetrics      []string // Lists required daily metrics, see https://open-meteo.com/en/docs for valid metrics
	AirQualityMetrics []string // List of required air quality metrics
	SatelliteMetrics  []string // List of required satellite metrics
	StartDate         string   // Start date for historical data (format: YYYY-MM-DD)
	EndDate           string   // End date for historical data (format: YYYY-MM-DD)
	SeasonalForecast  bool     // Enable seasonal forecast
	ForecastMonths    int      // Number of months to forecast (1-6)
}

func urlFromOptions(baseURL string, loc Location, opts *Options) string {
	// TODO: Validate the Options are valid
	url := fmt.Sprintf(`%s?latitude=%f&longitude=%f&current_weather=true`, baseURL, loc.lat, loc.lon)
	if opts == nil {
		return url
	}

	if opts.TemperatureUnit != "" {
		url = fmt.Sprintf(`%s&temperature_unit=%s`, url, opts.TemperatureUnit)
	}
	if opts.WindspeedUnit != "" {
		url = fmt.Sprintf(`%s&windspeed_unit=%s`, url, opts.WindspeedUnit)
	}
	if opts.PrecipitationUnit != "" {
		url = fmt.Sprintf(`%s&precipitation_unit=%s`, url, opts.PrecipitationUnit)
	}
	if opts.Timezone != "" {
		url = fmt.Sprintf(`%s&timezone=%s`, url, opts.Timezone)
	}
	if opts.PastDays != 0 {
		url = fmt.Sprintf(`%s&past_days=%d`, url, opts.PastDays)
	}

	if len(opts.HourlyMetrics) > 0 {
		metrics := strings.Join(opts.HourlyMetrics, ",")
		url = fmt.Sprintf(`%s&hourly=%s`, url, metrics)
	}

	if len(opts.DailyMetrics) > 0 {
		metrics := strings.Join(opts.DailyMetrics, ",")
		url = fmt.Sprintf(`%s&daily=%s`, url, metrics)
	}

	if len(opts.AirQualityMetrics) > 0 {
		metrics := strings.Join(opts.AirQualityMetrics, ",")
		url = fmt.Sprintf(`%s&air_quality=%s`, url, metrics)
	}

	if len(opts.SatelliteMetrics) > 0 {
		metrics := strings.Join(opts.SatelliteMetrics, ",")
		url = fmt.Sprintf(`%s&satellite=%s`, url, metrics)
	}

	if opts.StartDate != "" {
		url = fmt.Sprintf(`%s&start_date=%s`, url, opts.StartDate)
	}
	if opts.EndDate != "" {
		url = fmt.Sprintf(`%s&end_date=%s`, url, opts.EndDate)
	}

	if opts.SeasonalForecast {
		url = fmt.Sprintf(`%s&seasonal=true`, url)
		if opts.ForecastMonths > 0 && opts.ForecastMonths <= 6 {
			url = fmt.Sprintf(`%s&forecast_months=%d`, url, opts.ForecastMonths)
		}
	}

	return url
}

func (c *Client) Get(ctx context.Context, loc Location, opts *Options) ([]byte, error) {
	if err := c.RateLimiter.Wait(ctx); err != nil {
		return nil, ErrRateLimit{Message: "Rate limit exceeded"}
	}

	url := urlFromOptions(c.URL, loc, opts)
	if c.APIKey != "" {
		url = fmt.Sprintf("%s&apikey=%s", url, c.APIKey)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 429 {
		return nil, ErrRateLimit{Message: "Rate limit exceeded"}
	}

	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("%s - %s", res.Status, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
