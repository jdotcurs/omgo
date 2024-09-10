package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/hectormalot/omgo"
)

type CityWeather struct {
	Name             string
	Temperature      float64
	Humidity         float64
	WindSpeed        float64
	AirQuality       float64
	CloudCover       float64
	PrecipitationSum float64
}

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
