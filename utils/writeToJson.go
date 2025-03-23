package utils

// IMPORTS
import (
	"encoding/json"
	"fmt"
	"os"
)

// FUNCTIONS
func MakeCardListJson(data []string, dataFilepath string) {
	/*
	* This function takes a string list and writes it to a JSON file.
	 */

	// Creates the file and error check
	file, err := os.Create(dataFilepath + dataFilepath)
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
