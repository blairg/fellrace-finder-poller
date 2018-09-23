package parseresults

import (
	"fmt"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

// GetRaceLinks get all the race links from the race page
func GetRaceLinks(htmlContent string) (raceLinks []string) {
	if htmlContent == "" {
		return raceLinks
	}

	raceNodes, err := html.Parse(strings.NewReader(htmlContent))

	if err != nil {
		fmt.Println("html read error")
	}

	processRacePage(raceNodes, &raceLinks)

	return raceLinks
}

func processRacePage(node *html.Node, raceLinks *[]string) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func(raceLinks *[]string) {
		defer wg.Done()
		parseRaceDetails(node, raceLinks)
	}(raceLinks)
	wg.Wait()
}

func parseRaceDetails(node *html.Node, raceLinks *[]string) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attribute := range node.Attr {
			if attribute.Key == "href" {
				if strings.Contains(attribute.Val, "races.php?id=") {
					splitString := strings.Split(attribute.Val, "=")
					*raceLinks = append(*raceLinks, splitString[1])
				}
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		parseRaceDetails(child, raceLinks)
	}
}
