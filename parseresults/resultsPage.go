package parseresults

import (
	"fmt"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

// GetResultLinks get all the race links from the results page
func GetResultLinks(htmlContent string) (resultLinks []string) {
	if htmlContent == "" {
		return resultLinks
	}

	resultsNodes, err := html.Parse(strings.NewReader(htmlContent))

	if err != nil {
		fmt.Println("html read error")
	}

	processRaceResult(resultsNodes, &resultLinks)

	return resultLinks
}

func processRaceResult(node *html.Node, resultsLinks *[]string) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func(resultsLinks *[]string) {
		defer wg.Done()
		parseRaceResultDetails(node, resultsLinks)
	}(resultsLinks)
	wg.Wait()
}

func parseRaceResultDetails(node *html.Node, resultsLinks *[]string) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attribute := range node.Attr {
			if attribute.Key == "href" {
				if strings.Contains(attribute.Val, "results.php?id=") {
					splitString := strings.Split(attribute.Val, "=")
					*resultsLinks = append(*resultsLinks, splitString[1])
				}
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		parseRaceResultDetails(child, resultsLinks)
	}
}
