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
	LevelStats []CardLevelStats
	Hitspeed   string
	Speed      string
	Count      string
	Range      string
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

func GetCardLevelStats(levelStats []CardLevelStats, e *colly.HTMLElement) []CardLevelStats {
	e.ForEach("tbody:first-of-type tr", func(index int, row *colly.HTMLElement) {
		level := strings.TrimSpace(row.DOM.Find("th:first-of-type").Children().Remove().End().Text())
		hitpoints := row.ChildText("td:first-of-type")
		damage := row.ChildText("td:last-of-type")

		levelStats = append(levelStats, CardLevelStats{
			Level:     level,
			Hitpoints: hitpoints,
			Damage:    damage,
		})
	})
	return levelStats
}

func GetCardInfo(cardName string, c *colly.Collector) {
	var cardInfo CardInfo

	// VARIABLES
	// --- URLs
	var cardUrl string = urlPrefix + "/card/detail/" + cardName

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

	// --- Adds card level stats to cardInfo
	c.OnHTML(statTable, func(e *colly.HTMLElement) {
		cardInfo.LevelStats = GetCardLevelStats(cardInfo.LevelStats, e)
	})
	// --- Gets other card info
	var cardOtherStats string = "body > main > article > section.bg-gradient-to-br.from-gray-body.to-gray-dark.px-page.py-3 > div > div:nth-child(2) > table > tbody"
	c.OnHTML(cardOtherStats, func(e *colly.HTMLElement) {
		e.ForEach("tr", func(index int, row *colly.HTMLElement) {
			statName := row.ChildText("th:first-of-type")
			statValue := row.ChildText("td:last-of-type")
			switch statName {
			case "Hitspeed":
				cardInfo.Hitspeed = statValue
			case "Speed":
				cardInfo.Speed = statValue
			case "Count":
				cardInfo.Count = statValue
			case "Range":
				cardInfo.Range = statValue
			}
		})
	})

	// --- Starts the scrape, downloads the card's image, then writes the card info to a JSON file
	c.Visit(cardUrl)
	DownloadCardImage(GetCardImageURL(cardUrl, c), cardName)
	MakeCardListJson(cardName, []CardInfo{cardInfo})
}
