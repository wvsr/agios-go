package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type PlaceDetails struct {
	PlaceID                  string         `json:"place_id"`
	Name                     string         `json:"name"`
	SimpleAddress            string         `json:"simple_address"`
	Location                 map[string]any `json:"location"`
	Types                    []string       `json:"types"`
	FormattedAddress         string         `json:"formatted_address"`
	InternationalPhoneNumber string         `json:"international_phone_number"`
	Website                  string         `json:"website"`
	Rating                   float64        `json:"rating"`
	OpeningHours             map[string]any `json:"opening_hours"`
	PhotoURLs                []string       `json:"photo_urls"`
	PriceLevel               int            `json:"price_level"`
	BusinessStatus           string         `json:"business_status"`
	URL                      string         `json:"url"`
	UserRatingsTotal         int            `json:"user_ratings_total"`
}

var apiKey string

func init() {
	_ = godotenv.Load()
	apiKey = os.Getenv("GOOGLE_MAP_KEY")
	if apiKey == "" {
		log.Fatal("GOOGLE_MAP_KEY not found in environment variables")
	}
}

func getPlaceDetails(ctx context.Context, placeID string, maxPhotos int) (map[string]any, error) {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/details/json?place_id=%s&fields=name,formatted_address,international_phone_number,website,rating,opening_hours,photo,price_level,business_status,url,user_ratings_total&key=%s", placeID, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	place := result["result"].(map[string]any)
	photos := []string{}
	if p, ok := place["photos"].([]any); ok {
		for i, ph := range p {
			if i >= maxPhotos {
				break
			}
			ref := ph.(map[string]any)["photo_reference"].(string)
			photoURL := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/photo?maxwidth=400&photoreference=%s&key=%s", ref, apiKey)
			photos = append(photos, photoURL)
		}
	}
	place["photo_urls"] = photos
	return place, nil
}

func GetNearbyPlaces(ctx context.Context, lat, lng float64, radius int, placeType, keyword string, maxResults int) ([]PlaceDetails, error) {
	baseURL := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%.6f,%.6f&radius=%d&key=%s", lat, lng, radius, apiKey)
	if placeType != "" {
		baseURL += "&type=" + placeType
	}
	if keyword != "" {
		baseURL += "&keyword=" + keyword
	}

	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	summaries := data["results"].([]any)
	if len(summaries) == 0 {
		return nil, nil
	}

	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	combined := []PlaceDetails{}

	for i, s := range summaries {
		if i >= maxResults {
			break
		}
		summary := s.(map[string]any)
		placeID := summary["place_id"].(string)

		wg.Add(1)
		go func(summary map[string]any, placeID string) {
			defer wg.Done()
			details, err := getPlaceDetails(ctx, placeID, 3)
			if err != nil || details == nil {
				log.Printf("Skipping place_id %s due to error or missing details", placeID)
				return
			}

			mu.Lock()
			defer mu.Unlock()
			combined = append(combined, PlaceDetails{
				PlaceID:                  placeID,
				Name:                     summary["name"].(string),
				SimpleAddress:            summary["vicinity"].(string),
				Location:                 summary["geometry"].(map[string]any)["location"].(map[string]any),
				Types:                    toStringSlice(summary["types"]),
				FormattedAddress:         getString(details["formatted_address"]),
				InternationalPhoneNumber: getString(details["international_phone_number"]),
				Website:                  getString(details["website"]),
				Rating:                   getFloat(details["rating"]),
				OpeningHours:             getMap(details["opening_hours"]),
				PhotoURLs:                toStringSlice(details["photo_urls"]),
				PriceLevel:               getInt(details["price_level"]),
				BusinessStatus:           getString(details["business_status"]),
				URL:                      getString(details["url"]),
				UserRatingsTotal:         getInt(details["user_ratings_total"]),
			})
		}(summary, placeID)
	}
	wg.Wait()
	return combined, nil
}

func toStringSlice(raw any) []string {
	if raw == nil {
		return nil
	}
	arr := raw.([]any)
	s := []string{}
	for _, v := range arr {
		s = append(s, v.(string))
	}
	return s
}

func getString(v any) string {
	if v == nil {
		return ""
	}
	return v.(string)
}

func getFloat(v any) float64 {
	if v == nil {
		return 0
	}
	return v.(float64)
}

func getInt(v any) int {
	if v == nil {
		return 0
	}
	return int(v.(float64))
}

func getMap(v any) map[string]any {
	if v == nil {
		return nil
	}
	return v.(map[string]any)
}
