package omgo

import (
	"encoding/json"
	"fmt"
	"time"
)

type ForecastJSON struct {
	Latitude       float64
	Longitude      float64
	Elevation      float64
	GenerationTime float64                    `json:"generationtime_ms"`
	CurrentWeather CurrentWeather             `json:"current_weather"`
	HourlyUnits    map[string]string          `json:"hourly_units"`
	HourlyMetrics  map[string]json.RawMessage `json:"hourly"` // Parsed later, the API returns both Time and floats here
	DailyUnits     map[string]string          `json:"daily_units"`
	DailyMetrics   map[string]json.RawMessage `json:"daily"` // Parsed later, the API returns both Time and floats here

}

type Forecast struct {
	Latitude       float64
	Longitude      float64
	Elevation      float64
	GenerationTime float64
	CurrentWeather CurrentWeather
	HourlyUnits    map[string]string
	HourlyMetrics  map[string][]float64 // Parsed from ForecastJSON.HourlyMetrics
	HourlyTimes    []time.Time          // Parsed from ForecastJSON.HourlyMetrics
	DailyUnits     map[string]string
	DailyMetrics   map[string][]float64 // Parsed from ForecastJSON.DailyMetrics
	DailyTimes     []time.Time          // Parsed from ForecastJSON.DailyMetrics
}

type CurrentWeather struct {
	Temperature   float64
	Time          ApiTime
	WeatherCode   float64
	WindDirection float64
	WindSpeed     float64
}

// ParseBody converts the API response body into a Forecast struct
// Rationale: The API returns a map with both times as well as floats, this function
// unmarshalls in 2 steps in order to not return a map[string][]interface{}
func ParseBody(body []byte) (*Forecast, error) {
	f := &ForecastJSON{}
	err := json.Unmarshal(body, f)
	if err != nil {
		return nil, err
	}

	fc := &Forecast{
		Latitude:       f.Latitude,
		Longitude:      f.Longitude,
		Elevation:      f.Elevation,
		GenerationTime: f.GenerationTime,
		CurrentWeather: f.CurrentWeather,
		HourlyUnits:    f.HourlyUnits,
		HourlyTimes:    []time.Time{},
		HourlyMetrics:  make(map[string][]float64),
		DailyUnits:     f.DailyUnits,
		DailyTimes:     []time.Time{},
		DailyMetrics:   make(map[string][]float64),
	}

	for k, v := range f.HourlyMetrics {
		if k == "time" {
			// We unmarshal into an ApiTime array because of the custom formatting
			// of the timestamp in the API response
			target := []ApiTime{}
			err := json.Unmarshal(v, &target)
			if err != nil {
				return nil, err
			}

			for _, at := range target {
				fc.HourlyTimes = append(fc.HourlyTimes, at.Time)
			}

			continue
		}
		target := []float64{}
		err := json.Unmarshal(v, &target)
		if err != nil {
			return nil, err
		}
		fc.HourlyMetrics[k] = target
	}

	for k, v := range f.DailyMetrics {
		if k == "time" {
			// We unmarshal into an ApiTime array because of the custom formatting
			// of the timestamp in the API response
			target := []ApiDate{}
			err := json.Unmarshal(v, &target)
			if err != nil {
				return nil, err
			}

			for _, at := range target {
				fc.DailyTimes = append(fc.DailyTimes, at.Time)
			}

			continue
		}
		target := []float64{}
		err := json.Unmarshal(v, &target)
		if err != nil {
			return nil, err
		}
		fc.DailyMetrics[k] = target
	}

	return fc, nil
}

func ParseHistoricalBody(body []byte) (HistoricalData, error) {
	var data struct {
		Hourly struct {
			Time                   []string  `json:"time"`
			Temperature2m          []float64 `json:"temperature_2m"`
			RelativeHumidity2m     []float64 `json:"relative_humidity_2m"`
			DewPoint2m             []float64 `json:"dew_point_2m"`
			ApparentTemperature    []float64 `json:"apparent_temperature"`
			Precipitation          []float64 `json:"precipitation"`
			Rain                   []float64 `json:"rain"`
			Snowfall               []float64 `json:"snowfall"`
			WindSpeed10m           []float64 `json:"wind_speed_10m"`
			WindDirection10m       []float64 `json:"wind_direction_10m"`
			WindGusts10m           []float64 `json:"wind_gusts_10m"`
			ShortwaveRadiation     []float64 `json:"shortwave_radiation"`
			DirectNormalIrradiance []float64 `json:"direct_normal_irradiance"`
			DiffuseRadiation       []float64 `json:"diffuse_radiation"`
			CloudCover             []float64 `json:"cloud_cover"`
			Visibility             []float64 `json:"visibility"`
			WeatherCode            []int     `json:"weather_code"`
		} `json:"hourly"`
		Daily struct {
			Time                     []string  `json:"time"`
			WeatherCode              []int     `json:"weather_code"`
			Temperature2mMax         []float64 `json:"temperature_2m_max"`
			Temperature2mMin         []float64 `json:"temperature_2m_min"`
			ApparentTemperatureMax   []float64 `json:"apparent_temperature_max"`
			ApparentTemperatureMin   []float64 `json:"apparent_temperature_min"`
			Sunrise                  []string  `json:"sunrise"`
			Sunset                   []string  `json:"sunset"`
			PrecipitationSum         []float64 `json:"precipitation_sum"`
			RainSum                  []float64 `json:"rain_sum"`
			SnowfallSum              []float64 `json:"snowfall_sum"`
			PrecipitationHours       []float64 `json:"precipitation_hours"`
			WindSpeed10mMax          []float64 `json:"wind_speed_10m_max"`
			WindGusts10mMax          []float64 `json:"wind_gusts_10m_max"`
			WindDirection10mDominant []float64 `json:"wind_direction_10m_dominant"`
			ShortwaveRadiationSum    []float64 `json:"shortwave_radiation_sum"`
		} `json:"daily"`
	}

	err := json.Unmarshal(body, &data)
	if err != nil {
		return HistoricalData{}, err
	}

	historicalData := HistoricalData{
		HourlyData: HourlyData{
			Time:                   make([]time.Time, len(data.Hourly.Time)),
			Temperature2m:          data.Hourly.Temperature2m,
			RelativeHumidity2m:     data.Hourly.RelativeHumidity2m,
			DewPoint2m:             data.Hourly.DewPoint2m,
			ApparentTemperature:    data.Hourly.ApparentTemperature,
			Precipitation:          data.Hourly.Precipitation,
			Rain:                   data.Hourly.Rain,
			Snowfall:               data.Hourly.Snowfall,
			WindSpeed10m:           data.Hourly.WindSpeed10m,
			WindDirection10m:       data.Hourly.WindDirection10m,
			WindGusts10m:           data.Hourly.WindGusts10m,
			ShortwaveRadiation:     data.Hourly.ShortwaveRadiation,
			DirectNormalIrradiance: data.Hourly.DirectNormalIrradiance,
			DiffuseRadiation:       data.Hourly.DiffuseRadiation,
			CloudCover:             data.Hourly.CloudCover,
			Visibility:             data.Hourly.Visibility,
			WeatherCode:            data.Hourly.WeatherCode,
		},
		DailyData: DailyData{
			Time:                     make([]time.Time, len(data.Daily.Time)),
			WeatherCode:              data.Daily.WeatherCode,
			Temperature2mMax:         data.Daily.Temperature2mMax,
			Temperature2mMin:         data.Daily.Temperature2mMin,
			ApparentTemperatureMax:   data.Daily.ApparentTemperatureMax,
			ApparentTemperatureMin:   data.Daily.ApparentTemperatureMin,
			Sunrise:                  make([]time.Time, len(data.Daily.Sunrise)),
			Sunset:                   make([]time.Time, len(data.Daily.Sunset)),
			PrecipitationSum:         data.Daily.PrecipitationSum,
			RainSum:                  data.Daily.RainSum,
			SnowfallSum:              data.Daily.SnowfallSum,
			PrecipitationHours:       data.Daily.PrecipitationHours,
			WindSpeed10mMax:          data.Daily.WindSpeed10mMax,
			WindGusts10mMax:          data.Daily.WindGusts10mMax,
			WindDirection10mDominant: data.Daily.WindDirection10mDominant,
			ShortwaveRadiationSum:    data.Daily.ShortwaveRadiationSum,
		},
	}

	for i, timeStr := range data.Hourly.Time {
		t, err := time.Parse("2006-01-02T15:04", timeStr)
		if err != nil {
			return HistoricalData{}, fmt.Errorf("failed to parse hourly time: %w", err)
		}
		historicalData.HourlyData.Time[i] = t
	}

	for i, timeStr := range data.Daily.Time {
		t, err := time.Parse("2006-01-02", timeStr)
		if err != nil {
			return HistoricalData{}, fmt.Errorf("failed to parse daily time: %w", err)
		}
		historicalData.DailyData.Time[i] = t
	}

	for i, timeStr := range data.Daily.Sunrise {
		t, err := time.Parse("2006-01-02T15:04", timeStr)
		if err != nil {
			return HistoricalData{}, fmt.Errorf("failed to parse sunrise time: %w", err)
		}
		historicalData.DailyData.Sunrise[i] = t
	}

	for i, timeStr := range data.Daily.Sunset {
		t, err := time.Parse("2006-01-02T15:04", timeStr)
		if err != nil {
			return HistoricalData{}, fmt.Errorf("failed to parse sunset time: %w", err)
		}
		historicalData.DailyData.Sunset[i] = t
	}

	return historicalData, nil
}
