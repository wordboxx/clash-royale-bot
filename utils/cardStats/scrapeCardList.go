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
func GetCardNames() map[string]interface{} {
	/*
	* This function finds the links to the card details page
	* and extracts the card name from the URL.
	* It returns a list of card names.
	 */

	// VARIABLES
	cardNames := make(map[string]interface{})

	// FUNCTIONS
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
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	// --- After making a request
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished scraping card names")
	})

	// SCRAPING FUNCTIONS
	// --- Find all links
	// TODO: Add the card scrape function here OR send the info to a function that doesn't make a request/new collector
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// Select link
		href := e.Attr("href")

		// Isolate card name from end of link
		if strings.HasPrefix(href, "/card/detail/") {
			cardName := strings.TrimPrefix(href, "/card/detail/")

			// Extract card's stats
			cardNames[cardName] = GetCardInfo(cardName)

			// Sleep for 3 seconds to avoid getting blocked
			time.Sleep(3 * time.Second)
		}
	})

	// GOES TO LINK TO START SCRAPING
	c.Visit(link)
	return cardNames
}
