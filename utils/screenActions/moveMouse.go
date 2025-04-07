// Testing, change package name and stuff later
package main

// IMPORTS
import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
	sx, sy := robotgo.GetScreenSize()
	fmt.Println(sx, sy)
}
