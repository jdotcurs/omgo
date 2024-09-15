# Open-Meteo-Go (Fork)

This is a fork of the [hectormalot/omgo](https://github.com/hectormalot/omgo) repository, expanding on the original work by hectormalot. This fork aims to enhance the Go client for the [Open-Meteo](https://open-meteo.com) API with additional features.

## Fork Information

- Original Repository: [hectormalot/omgo](https://github.com/hectormalot/omgo)
- Fork Maintainer: [jdotcurs](https://github.com/jdotcurs)
- Last Updated: September 10, 2024

We would love to potentially merge these changes back into the original repository in the future.

## Features

- Current weather data retrieval
- Weather forecasts
- Historical weather data retrieval
- Air quality information
- Seasonal forecasts
- Customizable options for data retrieval
- Support for multiple locations
- Temperature unit conversion (Celsius, Fahrenheit)
- Wind speed unit options (km/h, m/s, mph, knots)
- Precipitation unit options (mm, inch)
- Timezone support

## Installation

To install the Open-Meteo-Go client, use the following command:

```bash
go get github.com/jdotcurs/omgo
```


## Usage

Here's a basic example of how to use this forked version of the Open-Meteo-Go client:

```go
package main
import (
"context"
"fmt"
"log"
"github.com/jdotcurs/omgo"
)
func main() {
client, err := omgo.NewClient()
if err != nil {
log.Fatalf("Failed to create client: %v", err)
}
loc, err := omgo.NewLocation(52.52, 13.41) // Berlin
if err != nil {
log.Fatalf("Failed to create location: %v", err)
}
opts := &omgo.Options{
TemperatureUnit: "celsius",
WindspeedUnit: "kmh",
DailyMetrics: []string{"temperature_2m_max", "temperature_2m_min"},
}
forecast, err := client.Forecast(context.Background(), loc, opts)
if err != nil {
log.Fatalf("Failed to get forecast: %v", err)
}
fmt.Printf("Current temperature in Berlin: %.1f°C\n", forecast.CurrentWeather.Temperature)
fmt.Printf("Max temperature tomorrow: %.1f°C\n", forecast.DailyMetrics["temperature_2m_max"][1])
fmt.Printf("Min temperature tomorrow: %.1f°C\n", forecast.DailyMetrics["temperature_2m_min"][1])
}
```

## Advanced Usage

To retrieve historical weather data:

```go
endDate := time.Now().AddDate(0, 0, -1)
startDate := endDate.AddDate(0, 0, -30)
opts := &omgo.Options{
StartDate: startDate.Format("2006-01-02"),
EndDate: endDate.Format("2006-01-02"),
DailyMetrics: []string{"temperature_2m_max", "temperature_2m_min"},
}
historicalData, err := client.GetHistoricalData(context.Background(), loc, opts)
if err != nil {
log.Fatalf("Failed to get historical data: %v", err)
}
// Process historical data
```

To get a seasonal forecast:

```go
opts := &omgo.Options{
SeasonalForecast: true,
ForecastMonths: 3,
DailyMetrics: []string{"temperature_2m_max", "temperature_2m_min"},
}
seasonalForecast, err := client.GetSeasonalForecast(context.Background(), loc, opts)
if err != nil {
log.Fatalf("Failed to get seasonal forecast: %v", err)
}
// Process seasonal forecast data
```

For advanced usage examples, including forecasts, historical data, satellite data and air quality, please refer to the [`example/main.go`](example/main.go) file in this repository.

```go
func main() {
	client, err := omgo.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	cities := []struct {
		Name string
		Lat  float64
		Lon  float64
	}{
		{"New York", 40.7128, -74.0060},
		{"London", 51.5074, -0.1278},
		{"Tokyo", 35.6762, 139.6503},
		{"Sydney", -33.8688, 151.2093},
		{"Rio de Janeiro", -22.9068, -43.1729},
	}

	var cityWeathers []CityWeather

	for _, city := range cities {
		weather, err := getCityWeather(client, city.Name, city.Lat, city.Lon)
		if err != nil {
			log.Printf("Failed to get weather for %s: %v", city.Name, err)
			continue
		}
		cityWeathers = append(cityWeathers, weather)
	}

	printWeatherComparison(cityWeathers)
	printDetailedWeatherInfo(cityWeathers)
	printHistoricalData(client, cityWeathers[0], cities[0].Lat, cities[0].Lon)
	printSeasonalForecast(client, cities[0].Name, cities[0].Lat, cities[0].Lon)
}

func getCityWeather(client omgo.Client, cityName string, lat, lon float64) (CityWeather, error) {
	loc, err := omgo.NewLocation(lat, lon)
	if err != nil {
		return CityWeather{}, fmt.Errorf("failed to create location: %w", err)
	}

	opts := &omgo.Options{
		TemperatureUnit:   "celsius",
		WindspeedUnit:     "kmh",
		PrecipitationUnit: "mm",
		Timezone:          "UTC",
		HourlyMetrics:     []string{"relativehumidity_2m", "cloudcover"},
		DailyMetrics:      []string{"precipitation_sum"},
		AirQualityMetrics: []string{"pm2_5"},
	}

	forecast, err := client.Forecast(context.Background(), loc, opts)
	if err != nil {
		return CityWeather{}, fmt.Errorf("failed to get forecast: %w", err)
	}

	airQuality, err := client.GetAirQuality(context.Background(), loc, opts)
	if err != nil {
		return CityWeather{}, fmt.Errorf("failed to get air quality data: %w", err)
	}

	return CityWeather{
		Name:             cityName,
		Temperature:      forecast.CurrentWeather.Temperature,
		Humidity:         forecast.HourlyMetrics["relativehumidity_2m"][0],
		WindSpeed:        forecast.CurrentWeather.WindSpeed,
		AirQuality:       airQuality.PM2_5,
		CloudCover:       forecast.HourlyMetrics["cloudcover"][0],
		PrecipitationSum: forecast.DailyMetrics["precipitation_sum"][0],
	}, nil
}

func printWeatherComparison(cityWeathers []CityWeather) {
	fmt.Println("Global Weather Comparison")
	fmt.Println("=========================")

	sort.Slice(cityWeathers, func(i, j int) bool {
		return cityWeathers[i].Temperature > cityWeathers[j].Temperature
	})

	fmt.Println("\nCities ranked by temperature (hottest to coldest):")
	for _, cw := range cityWeathers {
		fmt.Printf("%s: %.1f°C\n", cw.Name, cw.Temperature)
	}
}

func printDetailedWeatherInfo(cityWeathers []CityWeather) {
	fmt.Println("\nDetailed Weather Information:")
	for _, cw := range cityWeathers {
		fmt.Printf("\n%s:\n", cw.Name)
		fmt.Printf("  Temperature: %.1f°C\n", cw.Temperature)
		fmt.Printf("  Humidity: %.1f%%\n", cw.Humidity)
		fmt.Printf("  Wind Speed: %.1f km/h\n", cw.WindSpeed)
		fmt.Printf("  Air Quality (PM2.5): %.2f\n", cw.AirQuality)
		fmt.Printf("  Cloud Cover: %.1f%%\n", cw.CloudCover)
		fmt.Printf("  Precipitation Sum: %.1f mm\n", cw.PrecipitationSum)
	}
}

func printHistoricalData(client omgo.Client, city CityWeather, lat, lon float64) {
	fmt.Printf("\nHistorical Data for %s (Last 7 days):\n", city.Name)
	historicalData, err := getHistoricalData(client, lat, lon)
	if err != nil {
		log.Printf("Failed to get historical data for %s: %v", city.Name, err)
	} else {
		for date, temp := range historicalData {
			fmt.Printf("%s: %.1f°C\n", date, temp)
		}
	}
}

func getHistoricalData(client omgo.Client, lat, lon float64) (map[string]float64, error) {
	loc, err := omgo.NewLocation(lat, lon)
	if err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}

	endDate := time.Now().AddDate(0, 0, -1)
	startDate := endDate.AddDate(0, 0, -6)

	opts := &omgo.Options{
		StartDate:    startDate.Format("2006-01-02"),
		EndDate:      endDate.Format("2006-01-02"),
		DailyMetrics: []string{"temperature_2m_max"},
	}

	historicalData, err := client.GetHistoricalData(context.Background(), loc, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical data: %w", err)
	}

	result := make(map[string]float64)
	for i, date := range historicalData.Forecast.DailyTimes {
		result[date.Format("2006-01-02")] = historicalData.Forecast.DailyMetrics["temperature_2m_max"][i]
	}

	return result, nil
}

func printSeasonalForecast(client omgo.Client, cityName string, lat, lon float64) {
	fmt.Printf("\nSeasonal Forecast for %s (Next 3 months):\n", cityName)
	seasonalData, err := getSeasonalForecast(client, lat, lon)
	if err != nil {
		log.Printf("Failed to get seasonal forecast for %s: %v", cityName, err)
	} else {
		fmt.Printf("Start Date: %s\n", seasonalData.StartDate.Format("2006-01-02"))
		fmt.Printf("End Date: %s\n", seasonalData.EndDate.Format("2006-01-02"))
		fmt.Printf("Average Max Temperature: %.1f°C\n", calculateAverage(seasonalData.Forecast.DailyMetrics["temperature_2m_max"]))
		fmt.Printf("Average Min Temperature: %.1f°C\n", calculateAverage(seasonalData.Forecast.DailyMetrics["temperature_2m_min"]))
	}
}

func getSeasonalForecast(client omgo.Client, lat, lon float64) (omgo.SeasonalForecast, error) {
	loc, err := omgo.NewLocation(lat, lon)
	if err != nil {
		return omgo.SeasonalForecast{}, fmt.Errorf("failed to create location: %w", err)
	}

	opts := &omgo.Options{
		SeasonalForecast: true,
		ForecastMonths:   3,
		DailyMetrics:     []string{"temperature_2m_max", "temperature_2m_min"},
	}

	seasonalForecast, err := client.GetSeasonalForecast(context.Background(), loc, opts)
	if err != nil {
		return omgo.SeasonalForecast{}, fmt.Errorf("failed to get seasonal forecast: %w", err)
	}

	return seasonalForecast, nil
}

func calculateAverage(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}
```

Output

```bash
➜  example git:(main) ✗ go run main.go
Global Weather Comparison
=========================

Cities ranked by temperature (hottest to coldest):
Tokyo: 31.5°C
Rio de Janeiro: 20.9°C
Sydney: 20.4°C
New York: 18.2°C
London: 11.2°C

Detailed Weather Information:

Tokyo:
  Temperature: 31.5°C
  Humidity: 79.0%
  Wind Speed: 6.5 km/h
  Air Quality (PM2.5): 0.00
  Cloud Cover: 13.0%
  Precipitation Sum: 0.0 mm

Rio de Janeiro:
  Temperature: 20.9°C
  Humidity: 77.0%
  Wind Speed: 3.8 km/h
  Air Quality (PM2.5): 0.00
  Cloud Cover: 0.0%
  Precipitation Sum: 0.0 mm

Sydney:
  Temperature: 20.4°C
  Humidity: 47.0%
  Wind Speed: 11.1 km/h
  Air Quality (PM2.5): 0.00
  Cloud Cover: 64.0%
  Precipitation Sum: 0.0 mm

New York:
  Temperature: 18.2°C
  Humidity: 39.0%
  Wind Speed: 10.7 km/h
  Air Quality (PM2.5): 0.00
  Cloud Cover: 1.0%
  Precipitation Sum: 0.0 mm

London:
  Temperature: 11.2°C
  Humidity: 77.0%
  Wind Speed: 12.4 km/h
  Air Quality (PM2.5): 0.00
  Cloud Cover: 0.0%
  Precipitation Sum: 3.5 mm

Historical Data for Tokyo (Last 7 days):
2024-09-06: 25.7°C
2024-09-07: 24.8°C
2024-09-08: 22.0°C
2024-09-09: 25.4°C
2024-09-03: 24.2°C
2024-09-04: 25.6°C
2024-09-05: 25.7°C

Seasonal Forecast for New York (Next 3 months):
Start Date: 2024-09-10
End Date: 2024-12-10
Average Max Temperature: 25.5°C
Average Min Temperature: 17.5°C
```

## Contributing

Contributions to this fork are welcome! Please feel free to submit issues, fork the repository, and send pull requests. If you're interested in helping merge these improvements back into the original repository, please reach out to discuss the best approach.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Thanks to [hectormalot](https://github.com/hectormalot) for the original Open-Meteo-Go client
- Thanks to [Open-Meteo](https://open-meteo.com) for providing the weather API

For more detailed information on how to use this forked version of the Open-Meteo-Go client, please refer to the full documentation in the repository.

## Future work

I will continue to add features from the main Open Meteo repo to the Go SDK to make it more feature rich and useful.