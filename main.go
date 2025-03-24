package main

import (
	cardStats "clash-royale-bot/utils/cardStats"
)

func main() {
		cardStats.MakeCardListJson(cardStats.GetCardNames())
}
