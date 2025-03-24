package cardStats

import (
	"fmt"

	"github.com/gocolly/colly"
)

type CardInfo struct {
	Level     int
	Hitpoints string
	Damage    string
}

var urlPrefix string = "https://www.deckshop.pro"

func GetCardInfo(cardName string) []CardInfo {
	// VARIABLES
	// --- URLs
	var cardUrl string = urlPrefix + "/card/detail/" + cardName
	// --- CardInfo object to store card stats
	var cardInfo []CardInfo

	// BEFORE, AFTER, AND ERROR FUNCTIONS
	// --- Creates a new collector
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
		colly.AllowURLRevisit(),
	)

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
	// TODO: incorporate this function into scrapeCardList.go
	c.OnHTML(statTable, func(e *colly.HTMLElement) {
		e.ForEach("tbody:first-of-type tr", func(index int, row *colly.HTMLElement) {
			// TODO: Some cards start at weird levels. EX: archer-queen
			level := index + 1
			hitpoints := row.ChildText("td:first-of-type")
			damage := row.ChildText("td:last-of-type")
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
