package main

import (
	cardStats "clash-royale-bot/utils/cardStats"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()
	cardStats.GetCardInfo("witch", c)
}
