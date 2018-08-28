package download

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetRace Runs a HTTP GET Request at an endpoint and returns the value
func GetRace(raceID int) (htmlResponse string, success bool) {
	htmlResponse = ""
	success = false
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

	urlToGet := "http://www.fellrunner.org.uk/results.php?id=" + strconv.Itoa(raceID)

	response, error := netClient.Get(urlToGet)

	if error != nil {
		fmt.Println(error)
		return
	}

	htmlResponse = handleResponse(response, urlToGet)
	success = true

	return
}

func handleResponse(response *http.Response, urlToGet string) string {
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

	if strings.Contains(bodyAsString, "This page is not accessible.") {

		return ""
	}

	return bodyAsString
}
