package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/blairg/fellrace-finder-poller/download"
	"github.com/blairg/fellrace-finder-poller/parseresults"
	"github.com/blairg/fellrace-finder-poller/storage"
	"github.com/urfave/cli"
)

var results []parseresults.Result

func main() {
	app := cli.NewApp()
	app.Name = "fellrace-poller"
	app.Usage = "Downloads the results from http://www.fellrunner.org.uk/results.php"
	app.Action = func(c *cli.Context) error {
		getResults()

		fmt.Println("Going to store " + strconv.Itoa(len(results)) + " races")

		if len(results) > 0 {
			storage.StoreManyResults(results)
		}

		fmt.Println("Finished getting results")

		return nil
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

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

	newRaceIds := storage.FilterRaceIds(resultLinks)

	fmt.Println(newRaceIds)

	var wg sync.WaitGroup
	wg.Add(len(newRaceIds))

	for index, element := range newRaceIds {
		go func(element string) {
			defer wg.Done()
			fmt.Println(index, element)
			raceIndex64, err := strconv.ParseInt(element, 10, 32)

			if err != nil {
				log.Fatal(err)
			}

			raceIndex := int(raceIndex64)

			getAndStore(raceIndex)
		}(element)
	}
	wg.Wait()
}

func getAndStore(raceID int) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func(raceID int) {
		defer wg.Done()

		downloadAndStoreFiles(raceID)
	}(raceID)
	wg.Wait()
}

func downloadAndStoreFiles(raceID int) {
	raceIDString := strconv.Itoa(raceID)
	pathToSaveTo := "./results/"
	emptyPathToSaveTo := "./noResult/"
	fileExtension := ".html"
	fileLocation := pathToSaveTo + raceIDString + fileExtension
	resultsHTML, success := download.GetRace(raceID)

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
