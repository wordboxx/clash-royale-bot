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
	Level string            `json:"Level"`
	Stats map[string]string `json:"Stats"`
}

type CardInfo struct {
	LevelStats   []CardLevelStats
	Hitspeed     string
	Speed        string
	Count        string
	Range        string
	SpellRadius  string
	Duration     string
	TowerDamage  string
	AirTroop     bool
	WinCondition bool
	Spell        bool
	IsSpawner    bool
}

// FUNCTIONS
func DownloadCardImage(cardURL string, cardName string, c *colly.Collector) {
	// Get image URL from card page
	pathToImage := "body > main > article > section.bg-gradient-to-br.from-gray-body.to-gray-dark.px-page.py-3 > div > div:nth-child(1) > div.flex.items-center.gap-3 > img"
	var cardImageURL string

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
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

func getGeneralStats(c *colly.Collector, cardInfo *CardInfo) {
	var generalStatsTable string = "body > main > article > section.bg-gradient-to-br.from-gray-body.to-gray-dark.px-page.py-3 > div > div:nth-child(2) > table > tbody"
	c.OnHTML(generalStatsTable, func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			el.ForEach("th", func(_ int, th *colly.HTMLElement) {
				statName := strings.TrimSpace(th.Text)
				statValue := strings.TrimSpace(el.ChildText("td"))

				// Map general stats to CardInfo fields
				switch statName {
				case "Hit speed":
					cardInfo.Hitspeed = statValue
				case "Speed":
					cardInfo.Speed = statValue
				case "Count":
					cardInfo.Count = statValue
				case "Range":
					cardInfo.Range = statValue
				case "Spell radius":
					cardInfo.SpellRadius = statValue
				case "Duration":
					cardInfo.Duration = statValue
				}
			})
		})

	})
}

func getLevelStats(c *colly.Collector, cardInfo *CardInfo) {
	var levelStatsTable string = "body > main > article > section.mb-10 > div.grid.md\\:grid-cols-2.gap-5 > div:nth-child(1) > table"
	c.OnHTML(levelStatsTable, func(e *colly.HTMLElement) {

		// --- Get the stat names in the header row
		var statNames []string
		var levelStatsTableHeader string = "body > main > article > section.mb-10 > div.grid.md\\:grid-cols-2.gap-5 > div:nth-child(1) > table > thead > tr"
		e.ForEach(levelStatsTableHeader, func(_ int, el *colly.HTMLElement) {
			el.ForEach("th", func(_ int, th *colly.HTMLElement) {
				var statName string

				if th.DOM.Children().Length() > 0 {
					// Check if there's a link and use its href as the stat name
					if link := th.ChildAttr("a", "href"); link != "" {
						statName = link
					} else {
						statName = strings.TrimSpace(th.ChildText("span:first-child"))
					}
				} else {
					statName = strings.TrimSpace(th.Text)
				}
				// Only append if it's not the "Level" column and not empty
				if statName != "Level" && statName != "" {
					statNames = append(statNames, statName)
				}
			})
		})

		// --- Get the stat values in the body rows
		var levelStatsTableBody string = "body > main > article > section.mb-10 > div.grid.md\\:grid-cols-2.gap-5 > div:nth-child(1) > table > tbody"
		e.ForEach(levelStatsTableBody, func(_ int, el *colly.HTMLElement) {
			el.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
				var levelStat CardLevelStats
				levelStat.Stats = make(map[string]string)
				var currentLevel string

				// Extract the level
				tr.ForEach("th", func(_ int, th *colly.HTMLElement) {
					// Clean up the level text by taking only the first part before any newlines
					// This is to avoid getting the "Level 15\n(Mirrored)\n(Mirr.)" text
					levelText := strings.TrimSpace(th.Text)
					currentLevel = strings.Split(levelText, "\n")[0]
				})

				// Extract the stat values
				tr.ForEach("td", func(i int, td *colly.HTMLElement) {
					if i < len(statNames) {
						value := strings.TrimSpace(td.Text)
						// Only add to stats if the value is not empty
						if value != "" {
							levelStat.Stats[statNames[i]] = value
						}
					}
				})

				// Set the level after collecting all stats
				levelStat.Level = currentLevel

				// Append the levelStat to the cardInfo
				cardInfo.LevelStats = append(cardInfo.LevelStats, levelStat)
			})
		})
	})
}

func checkPropertiesAndRoles(c *colly.Collector, cardInfo *CardInfo) {
	// Checks more general properties and roles of the card

	// --- Check if the card is an air troop
	var airTroopSelector string = "body > main > article > section:nth-child(5) > div:nth-child(2) > div:nth-child(3) > a:nth-child(1)"
	c.OnHTML(airTroopSelector, func(e *colly.HTMLElement) {
		if e.Text == "Air troop" {
			cardInfo.AirTroop = true
		}
	})

	// --- Check if the card is a win condition
	var winConditionSelector string = "body > main > article > section:nth-child(5) > div:nth-child(2) > div.flex.flex-wrap.gap-2.mb-3 > a:nth-child(1)"
	c.OnHTML(winConditionSelector, func(e *colly.HTMLElement) {
		if e.Text == "Win condition" {
			cardInfo.WinCondition = true
		}
	})

	// --- Check if the card is a spell
	var spellSelector string = "body > main > article > section:nth-child(5) > div:nth-child(2) > div:nth-child(3) > a:nth-child(1)"
	c.OnHTML(spellSelector, func(e *colly.HTMLElement) {
		if e.Text == "Spell" {
			cardInfo.Spell = true
		}
	})

	// --- Check if the card is a spawner by checking level stat names
	for _, property := range cardInfo.LevelStats {
		keys := make([]string, 0, len(property.Stats))
		for k := range property.Stats {
			keys = append(keys, k)
		}
		for _, key := range keys {
			if strings.Contains(key, "/card/detail/") {
				cardInfo.IsSpawner = true
				fmt.Println("Spawner:", cardInfo.IsSpawner)
				return
			}
		}
	}
}

func GetCardInfo(cardName string, c *colly.Collector) CardInfo {
	var cardInfo CardInfo

	// VARIABLES
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
	getGeneralStats(c, &cardInfo)
	getLevelStats(c, &cardInfo)
	checkPropertiesAndRoles(c, &cardInfo)

	// --- Visits the card URL and starts the scrape
	c.Visit(cardUrl)

	// --- Pretty-print the full cardInfo data as JSON
	// cardInfoJSON, err := json.MarshalIndent(cardInfo, "", "  ")
	// if err != nil {
	// 	fmt.Println("Error formatting card info as JSON:", err)
	// } else {
	// 	fmt.Printf("Card Info for %s:\n%s\n", cardName, string(cardInfoJSON))
	// }

	return cardInfo
}
