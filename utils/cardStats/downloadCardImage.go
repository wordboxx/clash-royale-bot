package cardStats

// IMPORTS
import (
	"fmt"

	"github.com/gocolly/colly"
)

// FUNCTIONS
func GetCardImage(cardName string, c *colly.Collector) (imgSrc string) {
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Getting image for:", cardName)
	})

	c.OnHTML("img", func(e *colly.HTMLElement) {
		// TODO: Download file from URL
		// imgSrc = e.Attr("src")
		// imageURL := "deckshop.pro" + imgSrc
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Download successful", r.Request.URL)
	})

	return imgSrc
}
