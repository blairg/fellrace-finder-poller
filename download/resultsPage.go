package download

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
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

// GetRacePageList gets paginated race page
func GetRacePageList(pageID int) (htmlResponse string, success bool) {
	htmlResponse = ""
	success = false
	racePageURL := os.Getenv("RACE_PAGE_URL")

	if racePageURL == "" {
		fmt.Println("RACE_PAGE_URL not found")

		return
	}

	if strings.Contains(racePageURL, "?m=") || strings.Contains(racePageURL, "?y=") || strings.Contains(racePageURL, "?all") {
		racePageURL = racePageURL + "&p=" + strconv.Itoa(pageID)
	} else {
		racePageURL = racePageURL + "?p=" + strconv.Itoa(pageID)
	}

	fmt.Println("Getting " + racePageURL)

	response := getHTMLResponse(racePageURL)
	htmlResponse = handleResultsResponse(response, racePageURL, "races.php?id=")

	if htmlResponse != "" {
		success = true
	}

	return htmlResponse, success
}

func getHTMLResponse(url string) (response *http.Response) {
	errorFound := false
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		time.Sleep(time.Millisecond * 50)

		defer wg.Done()
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

		var err error
		response, err = netClient.Get(url)

		if err != nil {
			fmt.Println(err)

			errorFound = true
		}
	}()
	wg.Wait()

	if errorFound {
		return nil
	}

	return response
}

func handleResultsResponse(response *http.Response, urlToGet, textToLookFor string) string {
	if response.StatusCode != 200 && response.StatusCode != 304 {
		fmt.Println("Failed for " + urlToGet + " with status code " + strconv.Itoa(response.StatusCode))

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
