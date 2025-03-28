package cardStats

// IMPORTS
import (
	"fmt"
	"io"
	"net/http"
	"os"
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

// FUNCTIONS
func DownloadCardImage(cardImageURL string, cardName string) {
	// Send GET request to image URL
	response, err := http.Get(cardImageURL)
	if err != nil {
		fmt.Println("Error while downloading card image:", err)
		return
	}
	defer response.Body.Close()

	// Response check
	if response.StatusCode != http.StatusOK {
		fmt.Println("Failed to download image:", response.StatusCode)
		return
	}

	// Create file to store image
	file, err := os.Create(cardName + ".png")
	if err != nil {
		fmt.Println("Error while creating image file:", err)
		return
	}
	defer file.Close()

	// Copy image to file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("Error while copying image to file:", err)
		return
	}
	fmt.Println("Image downloaded successfully:", cardName+".png")

}

func GetCardImageURL(cardURL string, c *colly.Collector) string {
	// VARIABLES
	var pathToImage string = "body > main > article > section.bg-gradient-to-br.from-gray-body.to-gray-dark.px-page.py-3 > div > div:nth-child(1) > div.flex.items-center.gap-3 > img"
	var imageSrc string

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visitinggg:", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnHTML(pathToImage, func(e *colly.HTMLElement) {
		imageSrc = urlPrefix + e.Attr("src")
		fmt.Println("Image source:", imageSrc)
	})

	c.Visit(cardURL)
	return imageSrc
}

func GetCardInfo(cardName string, c *colly.Collector) {
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
		// Loops through each row of the table
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

	// --- Starts the scrape, downloads the card's image, then returns the card's info
	// TODO: call WriteToJson in this function, do not export the card info struct
	c.Visit(cardUrl)
	DownloadCardImage(GetCardImageURL(cardUrl, c), cardName)
	MakeCardListJson(cardName, cardInfo)
}
