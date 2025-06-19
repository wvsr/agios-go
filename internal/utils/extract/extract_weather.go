package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DailySummary represents basic daily weather data.
type DailySummary struct {
	Date            string   `json:"date"`
	TemperatureMaxC *float64 `json:"temperature_max_c,omitempty"`
	TemperatureMinC *float64 `json:"temperature_min_c,omitempty"`
	PrecipitationMM *float64 `json:"precipitation_mm,omitempty"`
	WindspeedMaxKMH *float64 `json:"windspeed_max_kmh,omitempty"`
	Sunrise         *string  `json:"sunrise,omitempty"`
	Sunset          *string  `json:"sunset,omitempty"`
}

// DailyWeather extends DailySummary with current conditions.
type DailyWeather struct {
	DailySummary
	TemperatureCurrentC        *float64 `json:"temperature_current_c,omitempty"`
	WindspeedCurrentKMH        *float64 `json:"windspeed_current_kmh,omitempty"`
	WinddirectionCurrentDeg    *int     `json:"winddirection_current_deg,omitempty"`
	WeathercodeCurrent         *int     `json:"weathercode_current,omitempty"`
	PressureMSLCurrentHPA      *float64 `json:"pressure_msl_current_hpa,omitempty"`
	RelativeHumidityCurrentPct *int     `json:"relative_humidity_current_percent,omitempty"`
}

const (
	geocodeAPIURL       = "https://geocoding-api.open-meteo.com/v1/search"
	forecastAPIURL      = "https://api.open-meteo.com/v1/forecast"
	defaultTimeout      = 10 * time.Second
	defaultForecastDays = 6
)

// Config holds configuration options for weather requests.
type Config struct {
	Timeout      time.Duration
	ForecastDays int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Timeout:      defaultTimeout,
		ForecastDays: defaultForecastDays,
	}
}

// GetWeatherForecast fetches weather forecast for a city or coordinates.
func GetWeatherForecast(ctx context.Context, city string, latitude, longitude *float64) ([]DailyWeather, error) {
	return GetWeatherForecastWithConfig(ctx, city, latitude, longitude, DefaultConfig())
}

// GetWeatherForecastWithConfig fetches weather forecast with custom configuration.
func GetWeatherForecastWithConfig(ctx context.Context, city string, latitude, longitude *float64, config *Config) ([]DailyWeather, error) {
	// Validate input
	if city == "" && (latitude == nil || longitude == nil) {
		return nil, errors.New("either city or coordinates must be provided")
	}

	if config == nil {
		config = DefaultConfig()
	}

	// Geocode if needed
	if latitude == nil || longitude == nil {
		coords, err := geocodeWithTimeout(ctx, city, config.Timeout)
		if err != nil {
			return nil, fmt.Errorf("geocoding failed: %w", err)
		}
		latitude, longitude = &coords.Lat, &coords.Lon
	}

	// Build query parameters
	params := url.Values{}
	params.Set("latitude", fmt.Sprintf("%.4f", *latitude))
	params.Set("longitude", fmt.Sprintf("%.4f", *longitude))
	params.Set("timezone", "auto")
	params.Set("current_weather", "true")
	params.Set("daily", "temperature_2m_max,temperature_2m_min,precipitation_sum,windspeed_10m_max,sunrise,sunset")
	params.Set("hourly", "pressure_msl,relative_humidity_2m")
	params.Set("forecast_days", fmt.Sprintf("%d", config.ForecastDays))

	// HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, forecastAPIURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	client := &http.Client{Timeout: config.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("forecast request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("forecast API returned status %d", resp.StatusCode)
	}

	// Response payload
	var payload forecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	dates := payload.Daily.Time
	if len(dates) == 0 {
		return nil, errors.New("no daily forecast data returned")
	}

	// Align current data with hourly data
	pressureCurr := alignFloat(payload.Hourly.Time, payload.Current.Time, payload.Hourly.Pressure)
	humidityCurr := alignInt(payload.Hourly.Time, payload.Current.Time, payload.Hourly.Humidity)

	// Construct result
	result := make([]DailyWeather, len(dates))
	for i, date := range dates {
		summary := DailySummary{
			Date:            date,
			TemperatureMaxC: safeFloat(payload.Daily.TempMax, i),
			TemperatureMinC: safeFloat(payload.Daily.TempMin, i),
			PrecipitationMM: safeFloat(payload.Daily.Precipitation, i),
			WindspeedMaxKMH: safeFloat(payload.Daily.WindspeedMax, i),
			Sunrise:         safeString(payload.Daily.Sunrise, i),
			Sunset:          safeString(payload.Daily.Sunset, i),
		}

		dw := DailyWeather{DailySummary: summary}

		// Add current conditions if this is today's forecast
		if isCurrentDay(payload.Current.Time, date) {
			t := payload.Current.Temperature
			w := payload.Current.Windspeed
			d := payload.Current.Winddirection
			c := payload.Current.Weathercode
			dw.TemperatureCurrentC = &t
			dw.WindspeedCurrentKMH = &w
			dw.WinddirectionCurrentDeg = &d
			dw.WeathercodeCurrent = &c
			dw.PressureMSLCurrentHPA = pressureCurr
			dw.RelativeHumidityCurrentPct = humidityCurr
		}
		result[i] = dw
	}

	return result, nil
}

// forecastResponse represents the API response structure.
type forecastResponse struct {
	Current struct {
		Time          string  `json:"time"`
		Temperature   float64 `json:"temperature"`
		Windspeed     float64 `json:"windspeed"`
		Winddirection int     `json:"winddirection"`
		Weathercode   int     `json:"weathercode"`
	} `json:"current_weather"`
	Daily struct {
		Time          []string  `json:"time"`
		TempMax       []float64 `json:"temperature_2m_max"`
		TempMin       []float64 `json:"temperature_2m_min"`
		Precipitation []float64 `json:"precipitation_sum"`
		WindspeedMax  []float64 `json:"windspeed_10m_max"`
		Sunrise       []string  `json:"sunrise"`
		Sunset        []string  `json:"sunset"`
	} `json:"daily"`
	Hourly struct {
		Time     []string  `json:"time"`
		Pressure []float64 `json:"pressure_msl"`
		Humidity []int     `json:"relative_humidity_2m"`
	} `json:"hourly"`
}

// geocodeResult holds geocoding response.
type geocodeResult struct {
	Results []struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Name      string  `json:"name"`
		Country   string  `json:"country"`
	} `json:"results"`
}

// geocodeWithTimeout fetches coordinates for a city with specified timeout.
func geocodeWithTimeout(ctx context.Context, city string, timeout time.Duration) (struct{ Lat, Lon float64 }, error) {
	apiURL := fmt.Sprintf("%s?name=%s&count=1", geocodeAPIURL, url.QueryEscape(city))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return struct{ Lat, Lon float64 }{}, fmt.Errorf("creating geocode request: %w", err)
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return struct{ Lat, Lon float64 }{}, fmt.Errorf("geocode request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return struct{ Lat, Lon float64 }{}, fmt.Errorf("geocoding API returned status %d", resp.StatusCode)
	}

	var res geocodeResult
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return struct{ Lat, Lon float64 }{}, fmt.Errorf("decoding geocode response: %w", err)
	}

	if len(res.Results) == 0 {
		return struct{ Lat, Lon float64 }{}, fmt.Errorf("no results found for city: %s", city)
	}

	return struct{ Lat, Lon float64 }{res.Results[0].Latitude, res.Results[0].Longitude}, nil
}

// isCurrentDay checks if the current timestamp belongs to the given date.
func isCurrentDay(currentTime, date string) bool {
	// Extract date part from current time (assuming ISO format)
	if len(currentTime) >= 10 {
		currentDate := currentTime[:10]
		return currentDate == date
	}
	// Fallback to prefix check
	return strings.HasPrefix(currentTime, date)
}

// alignFloat finds a float value at the matching timestamp.
func alignFloat(times []string, target string, vals []float64) *float64 {
	for i, t := range times {
		if t == target && i < len(vals) {
			return &vals[i]
		}
	}
	return nil
}

// alignInt finds an int value at the matching timestamp.
func alignInt(times []string, target string, vals []int) *int {
	for i, t := range times {
		if t == target && i < len(vals) {
			return &vals[i]
		}
	}
	return nil
}

// safeFloat safely accesses a float slice element.
func safeFloat(slice []float64, i int) *float64 {
	if i >= 0 && i < len(slice) {
		return &slice[i]
	}
	return nil
}

// safeString safely accesses a string slice element.
func safeString(slice []string, i int) *string {
	if i >= 0 && i < len(slice) {
		return &slice[i]
	}
	return nil
}
