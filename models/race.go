package models

import (
	"googlemaps.github.io/maps"
)

type GeoLocationSearch struct {
	GeoLocation maps.LatLng
	Address     string
}

type Distance struct {
	Kilometers float32 `json:"kilometres"`
	Miles      float32 `json:"miles"`
}

type Climb struct {
	Meters int `json:"meters"`
	Feet   int `json:"feet"`
}

type EntryFee struct {
	OnDay    float32 `json:"onDay"`
	PreEntry float32 `json:"preEntry"`
}

type RecordDetails struct {
	Name string `json:"name"`
	Time string `json:"time"`
	Year int    `json:"year"`
}

type Records struct {
	Male   RecordDetails `json:"male"`
	Female RecordDetails `json:"female"`
}

type GeoLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Race type
type Race struct {
	ID               int         `json:"id"`
	Name             string      `json:"name"`
	Date             string      `json:"date"`
	Time             string      `json:"time"`
	Country          string      `json:"country"`
	Region           string      `json:"region"`
	Category         string      `json:"category"`
	Website          string      `json:"website"`
	Distance         Distance    `json:"distance"`
	Climb            Climb       `json:"climb"`
	Venue            string      `json:"venue"`
	GeoLocation      GeoLocation `json:"geoLocation"`
	GMapImageURL     string      `json:"gmapImageUrl"`
	GridReference    string      `json:"gridReference"`
	SkillsExperience string      `json:"skillsExperience"`
	MinimumAge       int         `json:"minimumAge"`
	EntryFee         EntryFee    `json:"entryFee"`
	Records          Records     `json:"records"`
}
