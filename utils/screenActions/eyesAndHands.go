package main

import (
	"fmt"
	"os"

	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/bitmap"
)

var screenShotDir string = "utils/screenActions/screenShots/"

func directoryInit() {
	// Screenshot directory handling
	if _, err := os.Stat(screenShotDir); os.IsNotExist(err) {
		err := os.MkdirAll(screenShotDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
		fmt.Println("Directory created:", screenShotDir)
	}
}

func findObjectOnScreen() {
	// TODO: Gonna have to use OpenCV for this
	sx, sy := robotgo.GetScreenSize()

	// Capture the screen
	bit := robotgo.CaptureScreen(0, 0, sx, sy)
	if bit == nil {
		fmt.Println("Failed to capture screen")
		return
	}
	defer robotgo.FreeBitmap(bit)

	// Save the raw bitmap as PNG
	bitmap.Save(bit, screenShotDir+"test.png")

	// Example: Find an image on the screen
	imagePath := "/home/kingdavid/Documents/clash-royale-bot/utils/screenActions/screenShots/thingToFind.png"
	bitToFind := bitmap.Open(imagePath)
	if bitToFind == nil {
		fmt.Println("Failed to load image:", imagePath)
		return
	}
	defer robotgo.FreeBitmap(bitToFind)

	x, y := bitmap.Find(bitToFind)
	if x != -1 && y != -1 {
		robotgo.Move(x, y)
	} else {
		fmt.Println("Image not found on the screen.")
	}
}

func main() {
	directoryInit()

	for i := 0; i < 10; i++ {
		findObjectOnScreen()
	}
	fmt.Println("done")
}
