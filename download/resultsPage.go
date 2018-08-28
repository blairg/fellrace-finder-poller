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

// GetResults grabs the race results
func GetFellRunnerResults() (htmlResponse string, success bool) {
	htmlResponse = ""
	success = false
	resultsPageURL := os.Getenv("RESULTS_PAGE_URL")

	if resultsPageURL == "" {
		fmt.Println("RESULTS_PAGE_URL not found")

		return
	}

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

	response, error := netClient.Get(resultsPageURL)

	if error != nil {
		fmt.Println(error)

		return
	}

	htmlResponse = handleResultsResponse(response, resultsPageURL)
	success = true

	return htmlResponse, success
}

func handleResultsResponse(response *http.Response, urlToGet string) string {
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

	if !strings.Contains(bodyAsString, "Search for a race") {
		return ""
	}

	return bodyAsString
}
