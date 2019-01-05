package parseresults

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/blairg/fellrace-finder-poller/models"
	"golang.org/x/net/html"
)

// ParseResult extracts the runners from the HTML
func ParseResult(resultName, htmlContent string) models.Result {
	articleNodes, err := html.Parse(strings.NewReader(htmlContent))

	if err != nil {
		fmt.Println("html read error")
	}

	var parsedResult models.Result
	processResult(articleNodes, &parsedResult)

	raceID, _ := strconv.ParseInt(resultName, 10, 32)

	var result models.Result
	result.ID = int(raceID)
	result.Race = parsedResult.Race
	result.Date = parsedResult.Date
	result.Runners = parsedResult.Runners
	result.NumberOfRunners = parsedResult.NumberOfRunners

	return result
}

func isValidResultNode(node string) bool {
	switch node {
	case
		"h2",
		"table":
		return true
	}
	return false
}

func processResult(node *html.Node, result *models.Result) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func(result *models.Result) {
		parseResultDetails(node, result)
		wg.Done()
	}(result)
	wg.Wait()
}

func parseResultDetails(node *html.Node, result *models.Result) {
	if (node.Type == html.ElementNode) && isValidResultNode(node.Data) {
		for _, attribute := range node.Attr {
			if attribute.Key == "class" || attribute.Key == "id" {
				if attribute.Val == "title_committee" && node.FirstChild.Data != "Sponsors" {
					splitString := strings.Split(node.FirstChild.Data, " – ")
					result.Race = splitString[1]
					result.Date = splitString[0]
				}

				if attribute.Val == "posts-table" {
					var runners []models.Runner
					parseRunners(node, &runners)

					result.Runners = runners
					result.NumberOfRunners = len(runners)
				}
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		parseResultDetails(child, result)
	}
}

func parseRunners(node *html.Node, runners *[]models.Runner) []models.Runner {
	if node.Type == html.ElementNode && node.Data == "tr" {
		var runnerDetails []string
		parseRunner(node, &runnerDetails)

		if len(runnerDetails) == 5 {
			var eachRunner models.Runner
			eachRunner.Position = runnerDetails[0]
			eachRunner.Name = runnerDetails[1]
			eachRunner.Category = runnerDetails[2]
			eachRunner.Club = runnerDetails[3]
			eachRunner.Time = runnerDetails[4]

			*runners = append(*runners, eachRunner)
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		parseRunners(child, runners)
	}

	return *runners
}

func parseRunner(node *html.Node, eachRunner *[]string) []string {
	if node.Type == html.ElementNode && node.Data == "td" {
		if node.FirstChild != nil {
			*eachRunner = append(*eachRunner, node.FirstChild.Data)
		} else {
			*eachRunner = append(*eachRunner, "")
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		parseRunner(child, eachRunner)
	}

	return *eachRunner
}
