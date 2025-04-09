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

	var cardListDirectoryFilepath string = "data/cardStatFiles/"
	var cardListFilepath string = cardListDirectoryFilepath + cardName + ".json"

	// Creates the file and error check
	// --- Creates the directory if it doesn't exist
	if err := os.MkdirAll(cardListDirectoryFilepath, os.ModePerm); err != nil {
		fmt.Println("Error while creating card info directory:", err)
		return
	}

	// --- Creates the file
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
