package googlemaps

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

// GetCoordinates get geo-location from the Google API
func GetCoordinates(address string) (maps.LatLng, string) {
	var location maps.LatLng

	apiKey := os.Getenv("GOOGLE_API_KEY")

	if apiKey == "" {
		fmt.Println("GOOGLE_API_KEY not found")

		return location, ""
	}

	client, err := maps.NewClient(maps.WithAPIKey(apiKey))

	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	r := &maps.GeocodingRequest{
		Address: address,
	}

	geoResult, err := client.Geocode(context.Background(), r)

	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	return geoResult[0].Geometry.Location, geoResult[0].FormattedAddress
}
