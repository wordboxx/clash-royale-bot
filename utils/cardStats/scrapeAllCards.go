package cardStats

// IMPORTS
import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// GLOBAL VARIABLES
var link string = "https://www.deckshop.pro/card/list"

// FUNCTIONS
func ScrapeAllCards() {
	/*
	* This function finds the links to the card details page
	* and extracts the card name from the URL.
	* It returns a list of card names.
	 */

	// FUNCTIONS
	// Instantiate collector
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
		colly.AllowURLRevisit(),
	)

	// Before scrape, after scrape, and error handling
	// --- Before making a request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
	})

	// --- If an error occurs
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	// --- After making a request
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished scraping card names")
	})

	// Scraping functions
	// --- Loops through all cards
	var cardList string = "body > main > article > div:nth-child(2) > section:nth-child(2)"
	c.OnHTML(cardList+" a[href]", func(e *colly.HTMLElement) {
		// Select link
		href := e.Attr("href")

		// If link is to card details page
		if strings.HasPrefix(href, "/card/detail/") {
			// Isolate card name from end of link
			cardName := strings.TrimPrefix(href, "/card/detail/")

			// TODO: Do all card manipulation here rather than exporting
			GetCardInfo(cardName, c)

			// Sleep for 1 second to avoid getting blocked
			time.Sleep(1 * time.Second)
		}
	})

	// GOES TO LINK TO START SCRAPING
	c.Visit(link)
}
