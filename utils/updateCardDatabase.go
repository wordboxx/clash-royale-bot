package utils

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func GetCardNames(link string) {
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
		href := e.Attr("href")
		if href != "" {
			if strings.HasPrefix(href, "/card/detail/") {
				cardName := strings.TrimPrefix(href, "/card/detail/")
				fmt.Println(cardName)
			}
		}
	})

	c.Visit(link)
}
