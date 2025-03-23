package scraper

import (
	"fmt"

	"github.com/gocolly/colly"
)

func GetCardInfo(cardUrl string, cardName string) {
	var urlPrefix string = "https://www.deckshop.pro"
	cardUrl = urlPrefix + cardUrl

	// INSTANTIATE DEFAULT COLLECTOR
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
		colly.AllowURLRevisit(),
	)

	// BEFORE, AFTER, ON ERROR FUNCTIONS
	// --- Before making a request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
	})

	// --- If an error occurs
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	// --- After making a request
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished:", r.Request.URL)
	})

	// SCRAPING FUNCTIONS
	// TODO: Extract card info
	c.OnHTML(func(e *colly.HTMLElement) {
		// Extract card info
		cardInfo := e.ChildText("div.card-info")
		fmt.Println(cardInfo)
	})
}
