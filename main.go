package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/urfave/cli"

	"github.com/blairg/fellrace-finder-poller/download"
	"github.com/blairg/fellrace-finder-poller/helpers"
	"github.com/blairg/fellrace-finder-poller/parseresults"
	"github.com/blairg/fellrace-finder-poller/storage"
)

var results []parseresults.Result
var races []parseresults.Race

func main() {
	app := cli.NewApp()
	app.Name = "fellrace-poller"
	app.Usage = "Downloads the races and results from https://fellrunner.org.uk"
	app.Action = func(c *cli.Context) error {
		//Results
		var resultsWaitGroup sync.WaitGroup
		resultsWaitGroup.Add(1)

		go func() {
			defer resultsWaitGroup.Done()

			getResults()

			fmt.Println("Going to store " + strconv.Itoa(len(results)) + " results")

			if len(results) > 0 {
				storage.StoreManyResults(results)
			}

			fmt.Println("Finished getting results")
		}()
		resultsWaitGroup.Wait()

		// Races
		var racesWaitGroup sync.WaitGroup
		racesWaitGroup.Add(1)

		go func() {
			defer racesWaitGroup.Done()

			getRaces()

			fmt.Println("Going to store " + strconv.Itoa(len(races)) + " races")

			if len(races) > 0 {
				storage.StoreManyRaces(races)
			}

			fmt.Println("Finished getting races")
		}()
		racesWaitGroup.Wait()

		return nil
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

// Results
func getResults() {
	resultsHTML, success := download.GetFellRunnerResults()

	if !success {
		fmt.Println("Failed to get results")

		return
	} else {
		fmt.Println("Race Results found")
	}

	resultLinks := parseresults.GetResultLinks(resultsHTML)
	fmt.Println(resultLinks)

	newResultIds := storage.FilterIds(resultLinks, "races")

	fmt.Println(newResultIds)

	var wg sync.WaitGroup
	wg.Add(len(newResultIds))

	for _, resultID := range newResultIds {
		go func(resultID string) {
			defer wg.Done()
			//fmt.Println(resultId, element)
			raceIndex64, err := strconv.ParseInt(resultID, 10, 32)

			if err != nil {
				log.Fatal(err)
			}

			raceIndex := int(raceIndex64)

			if helpers.ArrayStringContains(newResultIds, resultID) {
				getAndStoreResult(raceIndex)
			}
		}(resultID)
	}
	wg.Wait()
}

func getAndStoreResult(raceID int) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func(raceID int) {
		defer wg.Done()

		downloadAndStoreResultFiles(raceID)
	}(raceID)
	wg.Wait()
}

func downloadAndStoreResultFiles(raceID int) {
	raceIDString := strconv.Itoa(raceID)
	pathToSaveTo := "./results/"
	emptyPathToSaveTo := "./noResult/"
	fileExtension := ".html"
	fileLocation := pathToSaveTo + raceIDString + fileExtension
	resultsHTML, success := download.GetRace("https://fellrunner.org.uk/results.php?id=", raceID)

	if resultsHTML != "" {
		fmt.Println("Storing " + fileLocation)
		storage.Store(raceIDString, pathToSaveTo, resultsHTML)
		results = append(results, parseresults.ParseResult(raceIDString, resultsHTML))
	} else {
		if success {
			fmt.Println("Storing No Data " + fileLocation)
			storage.Store(raceIDString, emptyPathToSaveTo, "")
		}
	}
}

// Races
func getRaces() {
	resultsHTML, success := download.GetFellRunnerRaces()

	if !success {
		fmt.Println("Failed to get races")

		return
	}

	fmt.Println("Races found")

	raceLinks := parseresults.GetRaceLinks(resultsHTML)
	fmt.Println(raceLinks)

	raceLinks = checkForRacePagination(raceLinks)

	newRaceIds := storage.FilterIds(raceLinks, "raceinfo")
	fmt.Println(newRaceIds)

	var wg sync.WaitGroup
	wg.Add(len(raceLinks))

	for _, raceID := range raceLinks {
		go func(raceID string) {
			defer wg.Done()
			//fmt.Println(index, element)
			raceIndex64, err := strconv.ParseInt(raceID, 10, 32)

			if err != nil {
				log.Fatal(err)
			}

			if helpers.ArrayStringContains(newRaceIds, raceID) {
				getAndStoreRace(int(raceIndex64))
			}
		}(raceID)
	}
	wg.Wait()
}

func checkForRacePagination(raceLinks []string) []string {
	// Check if results are paginated
	for i := 2; i < 1000000; i++ {
		resultsHTML := ""
		success := false
		var wg sync.WaitGroup
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			resultsHTML, success = download.GetRacePageList(i)
		}(i)
		wg.Wait()

		if !success {
			fmt.Println("Failed to get race page " + strconv.Itoa(i))

			break
		}

		fmt.Println("Races found for page " + strconv.Itoa(i))

		raceLinksFound := parseresults.GetRaceLinks(resultsHTML)

		for _, link := range raceLinksFound {
			raceLinks = append(raceLinks, link)
		}
	}

	return raceLinks
}

func getAndStoreRace(raceID int) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func(raceID int) {
		defer wg.Done()

		downloadAndStoreRaceFiles(raceID)
	}(raceID)
	wg.Wait()
}

func downloadAndStoreRaceFiles(raceID int) {
	raceIDString := strconv.Itoa(raceID)
	pathToSaveTo := "./races/"
	emptyPathToSaveTo := "./noRace/"
	// fileExtension := ".html"
	// fileLocation := pathToSaveTo + raceIDString + fileExtension
	resultsHTML, success := download.GetRace("https://fellrunner.org.uk/races.php?id=", raceID)

	if resultsHTML != "" {
		//fmt.Println("Storing " + fileLocation)
		storage.Store(raceIDString, pathToSaveTo, resultsHTML)
		races = append(races, parseresults.ParseRace(raceIDString, resultsHTML))
	} else {
		if success {
			//fmt.Println("Storing No Data " + fileLocation)
			storage.Store(raceIDString, emptyPathToSaveTo, "")
		}
	}
}
