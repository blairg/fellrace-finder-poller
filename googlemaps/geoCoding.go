package googlemaps

import (
	"fmt"
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
		fmt.Println("error connecting to GEO API: ", err)

		return location, ""
	}

	r := &maps.GeocodingRequest{
		Address: address,
		Region:  "GB",
	}

	geoResult, err := client.Geocode(context.Background(), r)

	if err != nil {
		fmt.Println("error getting GEO data: ", err)

		return location, ""
	}

	if len(geoResult) == 0 {
		return location, ""
	}

	return geoResult[0].Geometry.Location, geoResult[0].FormattedAddress
}
