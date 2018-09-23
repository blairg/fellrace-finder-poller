package download

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// GetFellRunnerResults grabs the race results
func GetFellRunnerResults() (htmlResponse string, success bool) {
	htmlResponse = ""
	success = false
	resultsPageURL := os.Getenv("RESULTS_PAGE_URL")

	if resultsPageURL == "" {
		fmt.Println("RESULTS_PAGE_URL not found")

		return
	}

	response := getHTMLResponse(resultsPageURL)
	htmlResponse = handleResultsResponse(response, resultsPageURL, "Search for a race")
	success = true

	return htmlResponse, success
}

// GetFellRunnerRaces grabs the race results
func GetFellRunnerRaces() (htmlResponse string, success bool) {
	htmlResponse = ""
	success = false
	racePageURL := os.Getenv("RACE_PAGE_URL")

	if racePageURL == "" {
		fmt.Println("RACE_PAGE_URL not found")

		return
	}

	response := getHTMLResponse(racePageURL)
	htmlResponse = handleResultsResponse(response, racePageURL, "Races this month")
	success = true

	return htmlResponse, success
}

func getHTMLResponse(url string) *http.Response {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	response, error := netClient.Get(url)

	if error != nil {
		fmt.Println(error)

		return nil
	}

	return response
}

func handleResultsResponse(response *http.Response, urlToGet, textToLookFor string) string {
	if response.StatusCode != 200 {
		fmt.Println("Failed for " + urlToGet)

		return ""
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		fmt.Println(err)
	}

	bodyAsString := string(body)

	if !strings.Contains(bodyAsString, textToLookFor) {
		return ""
	}

	return bodyAsString
}
