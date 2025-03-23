package utils

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gocolly/colly"
) // IMPORTS

// GLOBAL VARIABLES
var link string = "https://www.deckshop.pro/card/list"
var cardNames []string

func GetCardNames() []string {
	/*
	* This function finds the links to the card details page
	* and extracts the card name from the URL.
	 */

	// INSTANTIATE DEFAULT COLLECTOR
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
		colly.AllowURLRevisit(),
	)
	// BEFORE, AFTER, ON ERROR FUNCTIONS
	// --- Before making a request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// --- If an error occurs
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	// --- After making a request
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	// SCRAPING STARTS HERE
	// --- Find all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// Extract the link
		href := e.Attr("href")

		// If the link is not empty
		if href != "" {

			// Isolate links that start with "/card/detail/",
			// because they contain the card name at the end
			if strings.HasPrefix(href, "/card/detail/") {

				// Extract card name from URL
				// Example: /card/detail/ice-spirit
				cardName := strings.TrimPrefix(href, "/card/detail/")

				// Make sure card name is unique in the list
				if slices.Contains(cardNames, cardName) {
					return
				}
				cardNames = append(cardNames, cardName)
			}
		}
	})

	c.Visit(link)
	return cardNames
}

func WriteToJson(cardNames []string) {
	//TODO: Write the card names to a JSON file
}
