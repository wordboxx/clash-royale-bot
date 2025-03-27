package main

import (
	cardStats "clash-royale-bot/utils/cardStats"
)

func main() {

	c := cardStats.NewCollector()

	// var urlPrefix string = "https://www.deckshop.pro"
	// var cardUrl string = urlPrefix + "/card/detail/" + "barbarians"

	cardStats.GetCardInfo("ice-spirit", c)
}
