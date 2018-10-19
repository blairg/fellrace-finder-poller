package googlemaps

import (
	"fmt"
	"image"
	"log"
	"os"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

// GetStaticMap gets a static map based on an address
func GetStaticMap(address string, location maps.LatLng) image.Image {
	var mapImage image.Image

	apiKey := os.Getenv("GOOGLE_API_KEY")

	if apiKey == "" {
		fmt.Println("GOOGLE_API_KEY not found")

		return mapImage
	}

	var client *maps.Client
	var err error
	client, err = maps.NewClient(maps.WithAPIKey(apiKey))

	if err != nil {
		fmt.Println("error getting static map: ", err)

		return nil
	}

	var marker maps.Marker
	marker.Color = "green"
	marker.Location = []maps.LatLng{location}

	r := &maps.StaticMapRequest{
		Center: address,
		Zoom:   10,
		Size:   "600x300",
		// Scale:    *scale,
		// Format:   maps.Format(*format),
		// Language: *language,
		// Region:   *region,
		Markers: []maps.Marker{marker},
		MapType: maps.MapType("terrain"),
	}

	mapImageResult, err := client.StaticMap(context.Background(), r)

	if err != nil {
		log.Fatalf("fatal error: %s", err)

		return nil
	}

	return mapImageResult
}
