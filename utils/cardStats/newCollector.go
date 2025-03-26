package cardStats

import (
	"github.com/gocolly/colly"
)

// Collector with default settings for testing modules that use colly
func NewCollector() *colly.Collector {
	// Instantiate default collector
	return colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
		colly.AllowURLRevisit(),
	)
}
