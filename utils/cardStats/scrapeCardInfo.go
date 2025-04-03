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
type CardLevelStats struct {
	Level     string
	Hitpoints string
	Damage    string
}

type CardInfo struct {
	LevelStats  []CardLevelStats
	Hitspeed    string
	Speed       string
	Count       string
	Range       string
	SpellRadius string
	Duration    string
	TowerDamage string
}

// FUNCTIONS
func DownloadCardImage(cardURL string, cardName string, c *colly.Collector) {
	// Get image URL from card page
	pathToImage := "body > main > article > section.bg-gradient-to-br.from-gray-body.to-gray-dark.px-page.py-3 > div > div:nth-child(1) > div.flex.items-center.gap-3 > img"
	var cardImageURL string

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visitinggg:", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnHTML(pathToImage, func(e *colly.HTMLElement) {
		cardImageURL = urlPrefix + e.Attr("src")
		fmt.Println("Image source:", cardImageURL)
	})

	c.Visit(cardURL)

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
	var cardImageDirectoryFilepath string = "data/images/cardImages/"
	var cardImageFilepath string = cardImageDirectoryFilepath + cardName + ".png"

	if err := os.MkdirAll(cardImageDirectoryFilepath, os.ModePerm); err != nil {
		fmt.Println("Error while creating card image directory:", err)
		return
	}

	file, err := os.Create(cardImageFilepath)
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

func GetCardInfo(cardName string, c *colly.Collector) CardInfo {
	var cardInfo CardInfo
	var statNames []string

	// VARIABLES
	// --- URLs
	cardUrl := urlPrefix + "/card/detail/" + cardName

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
	// --- Get level stats (that increase with each level)
	levelStatsTable := "body > main > article > section.mb-10 > div.grid.md\\:grid-cols-2.gap-5 > div:nth-child(1) > table"
	c.OnHTML(levelStatsTable, func(e *colly.HTMLElement) {
		// Get all stat names
		e.ForEach("tr:not(tbody tr) > th", func(_ int, el *colly.HTMLElement) {
			// Get first th or first child if it has children
			var statName string
			if el.DOM.Children().Length() > 0 {
				statName = strings.TrimSpace(el.DOM.Children().First().Text())
			} else {
				statName = strings.TrimSpace(el.Text)
			}
			statNames = append(statNames, statName)
		})
		fmt.Println(statNames)

		// Get all stat values
		// Get all stat values from tbody rows
		e.ForEach("tbody tr", func(_ int, el *colly.HTMLElement) {
			var values []string

			el.ForEach("th", func(_ int, th *colly.HTMLElement) {
				values = append(values, strings.TrimSpace(th.Text))
			})

			el.ForEach("td", func(_ int, td *colly.HTMLElement) {
				values = append(values, strings.TrimSpace(td.Text))
			})
			// TODO: fix last tr entry (mirrored, level 15)
		})
	})

	// TODO: get other stats like radius, range, etc., from other table
	// otherStatsTable
	// c.OnHTML()

	// --- Visits the card URL and starts the scrape
	c.Visit(cardUrl)
	return cardInfo
}
