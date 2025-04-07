package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/bitmap"
)

func main() {
	// Variables
	var screenShotDir string = "utils/screenActions/screenShots/"
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
	fmt.Println("Bitmap saved as test.png in ", screenShotDir)
}
