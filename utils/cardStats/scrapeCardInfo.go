package cardStats

// IMPORTS
import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

// VARIABLES
var urlPrefix string = "https://www.deckshop.pro"

// STRUCTS
type CardInfo struct {
	Level     string
	Hitpoints string
	Damage    string
}

func GetCardInfo(cardName string, c *colly.Collector) []CardInfo {
	// VARIABLES
	// --- URLs
	var cardUrl string = urlPrefix + "/card/detail/" + cardName
	// --- CardInfo object to store card stats
	var cardInfo []CardInfo

	// BEFORE, AFTER, AND ERROR FUNCTIONS
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
		fmt.Println("Finished scraping card info for:", cardName)
	})

	// SCRAPING
	// --- Selector for the card's stat table
	statTable := "body > main > article > section.mb-10 > div.grid.md\\:grid-cols-2.gap-5 > div:first-of-type"

	// --- Iterates through each row of the table and collects the card stats
	// --- for each level
	c.OnHTML(statTable, func(e *colly.HTMLElement) {
		e.ForEach("tbody:first-of-type tr", func(index int, row *colly.HTMLElement) {
			// Extracts values from table
			// (isolates level number string for each level, but done like this to
			// remove excess formatting/text from level 15 strings, which are weird)
			level := strings.TrimSpace(row.DOM.Find("th:first-of-type").Children().Remove().End().Text())
			hitpoints := row.ChildText("td:first-of-type")
			damage := row.ChildText("td:last-of-type")

			// Appends values to cardInfo struct
			// TODO: get all different stats and types, like air/ground, etc.
			cardInfo = append(cardInfo, CardInfo{
				Level:     level,
				Hitpoints: hitpoints,
				Damage:    damage,
			})
		})
	})

	// --- Starts the scrape and returns the card info
	c.Visit(cardUrl)
	return cardInfo
}
