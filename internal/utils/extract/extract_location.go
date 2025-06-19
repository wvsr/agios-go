package utils

import (
	"log"
	"net"
	"os"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

var (
	dbPath = "GeoLite2-City.mmdb"
	reader *geoip2.Reader
	once   sync.Once
)

type LocationData struct {
	IP      string  `json:"ip"`
	City    string  `json:"city,omitempty"`
	Country string  `json:"country,omitempty"`
	Lat     float64 `json:"lat,omitempty"`
	Lon     float64 `json:"lon,omitempty"`
}

func initializeReader() {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Printf("GeoIP mmdb file not found at path: %s. Geolocation will be limited.", dbPath)
		return
	}

	var err error
	reader, err = geoip2.Open(dbPath)
	if err != nil {
		log.Printf("Failed to load GeoIP mmdb file from %s: %v", dbPath, err)
		reader = nil
		return
	}

	log.Printf("GeoIP database loaded successfully from: %s", dbPath)
}

func ExtractLocationFromIP(ipStr string) LocationData {
	once.Do(initializeReader)

	location := LocationData{
		IP: ipStr,
	}

	if reader == nil {
		log.Printf("GeoIP reader not available. Returning basic IP info for %s", ipStr)
		return location
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		log.Printf("Invalid IP address format: %s", ipStr)
		return location
	}

	record, err := reader.City(ip)
	if err != nil {
		log.Printf("Failed to resolve IP %s using GeoIP: %v", ipStr, err)
		return location
	}

	if record.City.Names != nil {
		location.City = record.City.Names["en"]
	}
	location.Country = record.Country.Names["en"]
	location.Lat = record.Location.Latitude
	location.Lon = record.Location.Longitude

	return location
}
