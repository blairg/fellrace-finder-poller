package parseresults

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"

	"github.com/blairg/fellrace-finder-poller/googlebucket"

	"googlemaps.github.io/maps"

	"golang.org/x/net/html"

	"github.com/blairg/fellrace-finder-poller/googlemaps"
	"github.com/blairg/fellrace-finder-poller/models"
	"github.com/blairg/fellrace-finder-poller/services"
)

// ParseRace extracts the races from the HTML
func ParseRace(raceID, htmlContent string) models.Race {
	raceIDParsed, _ := strconv.ParseInt(raceID, 10, 32)
	raceReader := strings.NewReader(htmlContent)

	var parsedRace models.Race
	parsedRace.ID = int(raceIDParsed)
	processRace(raceReader, &parsedRace)

	var race models.Race
	race.ID = int(raceIDParsed)
	race.Name = parsedRace.Name
	race.Date = parsedRace.Date
	race.Time = parsedRace.Time
	// race.Country = parsedRace.Country
	// race.Region = parsedRace.Region
	// race.Category = parsedRace.Category
	// race.Website = parsedRace.Website
	race.Distance = parsedRace.Distance
	race.Climb = parsedRace.Climb
	race.Venue = parsedRace.Venue
	race.GeoLocation = parsedRace.GeoLocation
	race.GMapImageURL = parsedRace.GMapImageURL
	// race.GridReference = parsedRace.GridReference
	// race.SkillsExperience = parsedRace.SkillsExperience
	// race.MinimumAge = parsedRace.MinimumAge
	// race.EntryFee = parsedRace.EntryFee
	race.Records = parsedRace.Records

	//fmt.Println(race)

	return race
}

func isValidHTMLTag(htmlTag string) bool {
	switch htmlTag {
	case
		"h2",
		"ul",
		"li":
		return true
	}
	return false
}

func processRace(reader io.Reader, race *models.Race) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func(race *models.Race) {
		defer wg.Done()
		parseHTML(reader, race)
	}(race)
	wg.Wait()
}

// Race Name
func isRaceNameToken(token html.Token) bool {
	if len(token.Attr) > 0 && token.Attr[0].Val == "title_races" {
		return true
	}

	return false
}

func splitRaceName(raceName string) string {
	nameSplit := strings.Split(raceName, "–")

	if len(nameSplit) != 2 {
		return ""
	}

	return strings.TrimLeft(strings.TrimRight(strings.Replace(nameSplit[1], "(R)", "", 1), " "), " ")
}

func getRaceName(isRaceName *bool, race *models.Race, token html.Token) string {
	if *isRaceName && race.Name == "" {
		*isRaceName = false

		return splitRaceName(token.Data)
	}

	if !*isRaceName && race.Name != "" {
		return race.Name
	}

	return ""
}

//

// Date and Time
func tidyDay(day string) string {
	formattedDay := day
	formattedDay = strings.Replace(formattedDay, "nd", "", 1)
	formattedDay = strings.Replace(formattedDay, "st", "", 1)
	formattedDay = strings.Replace(formattedDay, "rd", "", 1)
	formattedDay = strings.Replace(formattedDay, "th", "", 1)
	dayIntParsed, _ := strconv.ParseInt(formattedDay, 10, 32)
	formattedDay = fmt.Sprintf("%02d", dayIntParsed)

	return formattedDay
}

func getMonth(month string) string {
	switch month {
	case "Jan":
		return "01"
	case "Feb":
		return "02"
	case "Mar":
		return "03"
	case "Apr":
		return "04"
	case "May":
		return "05"
	case "Jun":
		return "06"
	case "Jul":
		return "07"
	case "Aug":
		return "08"
	case "Sep":
		return "09"
	case "Oct":
		return "10"
	case "Nov":
		return "11"
	case "Dec":
		return "12"
	}

	return "01"
}

func splitRaceDateAndTime(dateAndTime string) []string {
	trimmedDateAndTime := strings.TrimRight(strings.TrimLeft(dateAndTime, " "), " ")
	trimmedDateAndTime = strings.Replace(trimmedDateAndTime, "\n", "", -1)
	trimmedDateAndTime = strings.Replace(trimmedDateAndTime, "\t", "", -1)
	dateTimeSplit := strings.Split(trimmedDateAndTime, " at ")
	// Date is such - Sat 22nd Sep 2018
	dateSplit := strings.Split(dateTimeSplit[0], " ")
	day := tidyDay(dateSplit[1])
	month := getMonth(dateSplit[2])
	year := dateSplit[3]
	parsedDateFormat := day + "/" + month + "/" + year
	dateTimeSplit[0] = parsedDateFormat

	if len(dateTimeSplit) != 2 {
		return []string{parsedDateFormat, ""}
	}

	return dateTimeSplit
}

func isRaceDateAndTimeToken(token html.Token) bool {
	if strings.Contains(token.Data, "Date & time:") {
		return true
	}

	return false
}

func getRaceDateAndTime(isDateAndTime *bool, race *models.Race, token html.Token) []string {
	if *isDateAndTime && race.Date == "" {
		*isDateAndTime = false

		return splitRaceDateAndTime(token.Data)
	}

	if !*isDateAndTime && race.Date != "" {
		return []string{race.Date, race.Time}
	}

	return []string{"", ""}
}

//

// Distance
func splitDistance(distanceToSplit string) models.Distance {
	var distanceType models.Distance
	trimmedDistance := strings.TrimRight(strings.TrimLeft(distanceToSplit, " "), " ")
	trimmedDistance = strings.Replace(trimmedDistance, "\n", "", -1)
	trimmedDistance = strings.Replace(trimmedDistance, "\t", "", -1)
	distanceSplit := strings.Split(trimmedDistance, " / ")

	if len(distanceSplit) == 2 {
		kilometresParsed, _ := strconv.ParseFloat(strings.Replace(distanceSplit[0], "km", "", 1), 32)
		distanceType.Kilometers = float32(kilometresParsed)

		milesParsed, _ := strconv.ParseFloat(strings.Replace(distanceSplit[1], "m", "", 1), 32)
		distanceType.Miles = float32(milesParsed)
	}

	return distanceType
}

func isDistance(token html.Token) bool {
	if strings.Contains(token.Data, "Distance:") {
		return true
	}

	return false
}

func getDistance(isDistance *bool, race *models.Race, token html.Token) models.Distance {
	var distance models.Distance

	if *isDistance && race.Distance.Kilometers == 0 {
		*isDistance = false

		return splitDistance(token.Data)
	}

	if !*isDistance && race.Distance.Kilometers != 0 {
		return race.Distance
	}

	return distance
}

//

// Climb
func splitClimb(climbToSplit string) models.Climb {
	var climbType models.Climb
	trimmedClimb := strings.TrimRight(strings.TrimLeft(climbToSplit, " "), " ")
	trimmedClimb = strings.Replace(trimmedClimb, "\n", "", -1)
	trimmedClimb = strings.Replace(trimmedClimb, "\t", "", -1)
	climbSplit := strings.Split(trimmedClimb, " / ")

	if len(climbSplit) == 2 {
		metresParsed, _ := strconv.ParseInt(strings.Replace(climbSplit[0], "m", "", 1), 10, 32)
		climbType.Meters = int(metresParsed)

		feetParsed, _ := strconv.ParseInt(strings.Replace(climbSplit[1], "ft", "", 1), 10, 32)
		climbType.Feet = int(feetParsed)
	}

	return climbType
}

func isClimb(token html.Token) bool {
	if strings.Contains(token.Data, "Climb:") {
		return true
	}

	return false
}

func getClimb(isClimb *bool, race *models.Race, token html.Token) models.Climb {
	var climb models.Climb

	if *isClimb && race.Climb.Meters == 0 {
		*isClimb = false

		return splitClimb(token.Data)
	}

	if !*isClimb && race.Climb.Meters != 0 {
		return race.Climb
	}

	return climb
}

//

// Records
func splitRecordToken(r rune) bool {
	return r == '$'
}

func splitRecords(recordToSplit string) models.RecordDetails {
	var recordDetailsType models.RecordDetails

	if strings.Contains(recordToSplit, "No record information") {
		return recordDetailsType
	}

	trimmedRecord := strings.TrimRight(strings.TrimLeft(recordToSplit, " "), " ")
	trimmedRecord = strings.Replace(trimmedRecord, "\n", "", -1)
	trimmedRecord = strings.Replace(trimmedRecord, "\t", "", -1)
	trimmedRecord = strings.Replace(trimmedRecord, " – ", "$", -1)
	recordSplit := strings.FieldsFunc(trimmedRecord, splitRecordToken)

	if len(recordSplit) == 3 {
		recordDetailsType.Name = strings.TrimRight(strings.TrimLeft(recordSplit[0], " "), " ")
		recordDetailsType.Time = strings.TrimRight(strings.TrimLeft(recordSplit[1], " "), " ")
		yearParsed, _ := strconv.ParseInt(strings.TrimRight(strings.TrimLeft(recordSplit[2], " "), " "), 10, 32)
		recordDetailsType.Year = int(yearParsed)
	}

	return recordDetailsType
}

func isFemaleRecord(token html.Token) bool {
	if strings.Contains(token.Data, "Female:") {
		return true
	}

	return false
}

func isMaleRecord(token html.Token) bool {
	if strings.Contains(token.Data, "Male:") {
		return true
	}

	return false
}

func getFemaleRecord(isRecord *bool, race *models.Race, token html.Token) models.RecordDetails {
	var femaleRecord models.RecordDetails

	if *isRecord && race.Records.Female.Name == "" {
		*isRecord = false

		return splitRecords(token.Data)
	}

	if !*isRecord && race.Records.Female.Name != "" {
		return race.Records.Female
	}

	return femaleRecord
}

func getMaleRecord(isRecord *bool, race *models.Race, token html.Token) models.RecordDetails {
	var maleRecord models.RecordDetails

	if *isRecord && race.Records.Male.Name == "" {
		*isRecord = false

		return splitRecords(token.Data)
	}

	if !*isRecord && race.Records.Male.Name != "" {
		return race.Records.Male
	}

	return maleRecord
}

//

// Venue
func splitVenue(venueToSplit string) string {
	trimmedVenue := strings.TrimRight(strings.TrimLeft(venueToSplit, " "), " ")
	trimmedVenue = strings.Replace(trimmedVenue, "\n", "", -1)
	trimmedVenue = strings.Replace(trimmedVenue, "\t", "", -1)

	return trimmedVenue
}

func isVenue(token html.Token) bool {
	if strings.Contains(token.Data, "Venue:") {
		return true
	}

	return false
}

func getVenue(isVenue *bool, race *models.Race, token html.Token) string {
	var venue string

	if *isVenue && race.Venue == "" {
		*isVenue = false

		return splitVenue(token.Data)
	}

	if !*isVenue && race.Venue != "" {
		return race.Venue
	}

	return venue
}

//

func parseHTML(r io.Reader, race *models.Race) {
	d := html.NewTokenizer(r)
	isRaceName := false
	isDateAndTime := false
	isDistanceFound := false
	isClimbFound := false
	isFemaleRecordFound := false
	isMaleRecordFound := false
	isVenueFound := false

	for {
		tokenType := d.Next()
		if tokenType == html.ErrorToken {
			return
		}
		token := d.Token()

		switch tokenType {
		case html.StartTagToken: // <tag>
			if isValidHTMLTag(token.Data) {
				isRaceName = isRaceNameToken(token)
			}
		case html.TextToken: // text between start and end tag
			race.Name = getRaceName(&isRaceName, race, token)

			// Date and Time
			if isRaceDateAndTimeToken(token) {
				isDateAndTime = true
			} else {
				if isDateAndTime == true {
					dateAndTime := getRaceDateAndTime(&isDateAndTime, race, token)
					race.Date = dateAndTime[0]
					race.Time = dateAndTime[1]

					isDateAndTime = false
				}
			}

			// Distance
			if isDistance(token) {
				isDistanceFound = true
			} else {
				if isDistanceFound == true {
					distance := getDistance(&isDistanceFound, race, token)
					race.Distance.Kilometers = distance.Kilometers
					race.Distance.Miles = distance.Miles

					isDistanceFound = false
				}
			}

			// Climb
			if isClimb(token) {
				isClimbFound = true
			} else {
				if isClimbFound == true {
					climb := getClimb(&isClimbFound, race, token)
					race.Climb.Meters = climb.Meters
					race.Climb.Feet = climb.Feet

					isClimbFound = false
				}
			}

			// Records - Female
			if isFemaleRecord(token) {
				isFemaleRecordFound = true
			} else {
				if isFemaleRecordFound == true {
					femaleRecord := getFemaleRecord(&isFemaleRecordFound, race, token)
					race.Records.Female = femaleRecord

					isFemaleRecordFound = false
				}
			}

			// Records - Male
			if isMaleRecord(token) {
				isMaleRecordFound = true
			} else {
				if isMaleRecordFound == true {
					maleRecord := getMaleRecord(&isMaleRecordFound, race, token)
					race.Records.Male = maleRecord

					isMaleRecordFound = false
				}
			}

			// Venue
			if isVenue(token) {
				isVenueFound = true
			} else {
				if isVenueFound == true {

					// @TODO: Pull this into a function
					venue := getVenue(&isVenueFound, race, token)
					race.Venue = venue

					geoLocationChannel := make(chan models.GeoLocationSearch)
					go getCoordinates(venue, geoLocationChannel)
					geoLocationResult := <-geoLocationChannel

					if geoLocationResult.Address != "" {
						race.GeoLocation.Latitude = geoLocationResult.GeoLocation.Lat
						race.GeoLocation.Longitude = geoLocationResult.GeoLocation.Lng

						staticMapChannel := make(chan image.Image)
						go getStaticMap(geoLocationResult.Address, geoLocationResult.GeoLocation, staticMapChannel)
						staticMapImage := <-staticMapChannel

						if staticMapImage != nil {
							storeChannel := make(chan bool)
							imageName := strconv.Itoa(race.ID) + ".png"
							go storeImage(imageName, staticMapImage, storeChannel)
							stored := <-storeChannel

							if !stored {
								fmt.Println("Failed to upload " + imageName)
							} else {
								race.GMapImageURL = "https://storage.googleapis.com/fellrace-finder/maps/" + imageName
							}
						}
					}

					isVenueFound = false
				}
			}

		case html.EndTagToken: // </tag>
		case html.SelfClosingTagToken: // <tag/>

		}
	}
}

func getCoordinates(address string, geoLocationChannel chan models.GeoLocationSearch) {
	geoLocation := services.GetCoordinates(address)

	var geoSearch models.GeoLocationSearch
	geoSearch.GeoLocation = geoLocation
	geoSearch.Address = address

	geoLocationChannel <- geoSearch
}

func getStaticMap(address string, location maps.LatLng, staticMapChannel chan image.Image) {
	staticMapChannel <- googlemaps.GetStaticMap(address, location)
}

// StoreImage stores an image to the filesystem and the bucket
func storeImage(filename string, image image.Image, storeChannel chan bool) {
	buff := new(bytes.Buffer)

	// encode image to buffer
	err := png.Encode(buff, image)

	if err != nil {
		storeChannel <- false
		fmt.Println("failed to create buffer for image "+filename, err)
	}

	err = ioutil.WriteFile("./"+filename, buff.Bytes(), 0644)

	if err != nil {
		storeChannel <- false
		fmt.Println("failed to write for image "+filename, err)
	}

	storeChannel <- googlebucket.StoreObject(filename, filename)
}
