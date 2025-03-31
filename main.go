package main

import (
	cardStats "clash-royale-bot/utils/cardStats"
)

func main() {

	c := cardStats.NewCollector()

	cardStats.GetCardInfo("poison", c)
}
