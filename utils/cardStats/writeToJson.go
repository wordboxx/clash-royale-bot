package cardStats

// IMPORTS
import (
	"encoding/json"
	"fmt"
	"os"
)

// FUNCTIONS
func MakeCardListJson(cardName string, data []CardInfo) {
	/*
	* This function takes a map and writes it to a JSON file.
	 */

	var cardListFilepath string = "data/" + cardName + ".json"

	// Creates the file and error check
	// TODO: Makes file if not already existing
	file, err := os.Create(cardListFilepath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	// Defer closing the file until the function ends
	defer file.Close()

	// Create a new JSON encoder and sets indentation
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	// Calls the encoder to encode the card names
	// and returns an error if one occurs
	err = encoder.Encode(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
}
