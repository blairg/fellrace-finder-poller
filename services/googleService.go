package services

import (
	"fmt"

	"github.com/blairg/fellrace-finder-poller/googlemaps"
	"github.com/blairg/fellrace-finder-poller/storage"
	"googlemaps.github.io/maps"
)

// GetCoordinates get geolocation from an address. Tries database first Google API second.
func GetCoordinates(address string) maps.LatLng {
	// check in database first, if not in database hit Google
	race := storage.GetRaceByAddress(address)

	if race.GeoLocation.Latitude != 0 {
		var geoLocation maps.LatLng
		geoLocation.Lat = race.GeoLocation.Latitude
		geoLocation.Lat = race.GeoLocation.Longitude

		return geoLocation
	}

	fmt.Println("Getting geo from Google")

	coordinatesResult, _ := googlemaps.GetCoordinates(address)

	return coordinatesResult
}
